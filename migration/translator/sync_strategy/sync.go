package sync_strategy

import (
	"fmt"
	"sort"

	dt "internal/data_type"
	"internal/util"

	"github.com/amirkode/go-mongr8/collection"
	"github.com/amirkode/go-mongr8/collection/field"
	"github.com/amirkode/go-mongr8/migration/migrator"
	si "github.com/amirkode/go-mongr8/migration/translator/mongodb/schema_interpreter"
)

/*
This contains mechanism to sync schema across all sources
mainly from user-defined collection (mongr8/collection)
compared to migration files to define what actions must be added
during new migration file generation

example:
1. current migration files schema:
collections:
user_collection
{
	"name": "[a string]",
	"extras": {
		"ext1": "[a string]",
		"ext2": "[a string]",
		"ext3": "[a string]"
	},
	"login_history": [
		"2023-10-12 01:00:00",
		"2023-10-12 02:00:00",
		"2023-10-12 03:00:00",
	]
}
indexes:
- {"name": -1}
- {"login_history": -1}

2. latest user-defined schema:
collections:
user_collection
{
	"name": "[a string]",
	"age": [an integer],
	"extras": {
		"ext1": "[a string]",
		"ext2": [an integer],
	},
	"login_history": [
		"2023-10-12 01:00:00",
		"2023-10-12 02:00:00",
		"2023-10-12 03:00:00",
	]
}
indexes:
- {"name": 1}
- {"name": 1, "age": -1}
- {"login_history": 1}
- {"login_history": -1}


Actions that must be added:
- drop index `{"name": -1}` (not exist in latest user defined schema)
- drop field "extras"."ext3"
- add field "age" of integer
- add index {"name": 1}
- add index {"name": 1, "age": -1}
- add index {"login_history": 1}
- convert field "extras"."ext2" to integer

As shown on above example, possible actions for each usecase:
- New Field/Index: Will be added as plus-signed Action (Add)
- Unused Field/Index: Will be added as negative-signed Action (Drop)
- Field Type Conversion:
  - Supported: Any to String, Numeric to Numeric (int to double, double to int, etc)
    perform update query, i.e:
	db.collectionName.updateMany({}, [
       { $set: { "fieldName": { $toInt: "$fieldName" } } }
    ])
	*note:
	for other numerics to Int32 could produce error if
	the previous integer value exceed the limit of int32
  - Unsupported: Date to Numeric, Timestamp to Numeric, etc
  - Undefined: String to Any
- Index ordering reversal: drop previous index and followed adding new index

User options:
- Force Conversion: This will force conversion for undefined field type conversion
*/

// This returns initial `SignedCollection` instance from a `Collection`
func getSignedCollection(item collection.Collection) SignedCollection {
	signedMetadata := SignedMetadata{
		Metadata: item.Collection(),
		Sign:     SignPlus,
	}
	signedFields := make([]SignedField, len(item.Fields()))
	for index, f := range item.Fields() {
		signedFields[index] = SignedField{
			Field: f,
			Sign:  SignPlus,
		}
	}
	signedIndexes := make([]SignedIndex, len(item.Indexes()))
	for index, i := range item.Indexes() {
		signedIndexes[index] = SignedIndex{
			Index: i,
			Sign:  SignPlus,
		}
	}

	return SignedCollection{
		Metadata: signedMetadata,
		Fields:   signedFields,
		Indexes:  signedIndexes,
		Sign:     SignPlus,
	}
}

// This returns the syned collection with proper signs after mergin process
// `incoming` is the new schema must be matched, originally from
// the latest user defined collection in `[project dir]/mongr8/collection`
// `origin` is the current schema previosly generated by the migration and fetched from
// the latest migration files in `[project dir]/mongr8/migration`
func SyncCollections(incoming []collection.Collection, origin []collection.Collection) []SignedCollection {
	// convert incomingSchema to signed collections
	source1 := make([]SignedCollection, len(incoming))
	for index, item := range incoming {
		source1[index] = getSignedCollection(item)
	}

	// convert originalSchema to signed collections
	source2 := make([]SignedCollection, len(origin))
	for index, item := range origin {
		source2[index] = getSignedCollection(item)
	}

	// find union
	return Union(source1, source2)
}

// this returns Up and Down actions used for generating migration 
// `incoming` is the new schema must be matched, originally from
// the latest user defined collection in `[project dir]/mongr8/collection`
// `origin` is the current schema previosly generated by the migration and fetched from
// the latest migration files in `[project dir]/mongr8/migration`
func GetActions(incoming []collection.Collection, origin []collection.Collection) dt.Pair[[]si.Action, []si.Action] {
	upActionMap := map[string]si.Action{}
	downActionMap := map[string]si.Action{}
	signedCollections := SyncCollections(incoming, origin)

	// fill upActionMap and downActionMap
	fillActionMap := func(signedCollection SignedCollection, subActionType si.SubActionType) {
		// init up action
		upAction, ok := upActionMap[signedCollection.Key()]
		if !ok {
			upAction = si.Action{
				ActionKey: signedCollection.Key(),
			}
		}
		// init down action
		downAction, ok := downActionMap[signedCollection.Key()]
		if !ok {
			downAction = si.Action{
				ActionKey: signedCollection.Key(),
			}
		}

		var upSubAction *si.SubAction
		var downSubAction *si.SubAction

		schema := si.SubActionSchema{
			Collection: signedCollection.Metadata,
			Fields:     signedCollection.GetFields(),
			Indexes:    signedCollection.GetIndexes(),
		}

		switch subActionType {
		case si.SubActionTypeCreateCollection:
			upSubAction = si.SubActionCreateCollection(schema)
			downSubAction = si.SubActionDropCollection(schema)
		case si.SubActionTypeCreateIndex:
			upSubAction = si.SubActionCreateIndex(schema)
			downSubAction = si.SubActionDropIndex(schema)
		case si.SubActionTypeCreateField:
			upSubAction = si.SubActionCreateField(schema)
			downSubAction = si.SubActionDropField(schema)
		case si.SubActionTypeConvertField:
			// TODO: test this
			convertFromType := signedCollection.Fields[0].ConvertFrom().FieldDeepestType()
			convertToType := signedCollection.Fields[0].FieldDeepestType()
			// set up conversion
			schema.FieldConvertFrom = &convertFromType
			upSubAction = si.SubActionConvertField(schema)
			// set down conversion
			downSignedCollection := signedCollection.RefreshFieldAddresses()
			downSignedCollection.Fields[0].SetFieldDeepestType(convertFromType)
			downSchema := schema // assign new address
			downSchema.Fields = downSignedCollection.GetFields()
			downSchema.FieldConvertFrom = &convertToType
			downSubAction = si.SubActionConvertField(downSchema)
		case si.SubActionTypeDropCollection:
			upSubAction = si.SubActionDropCollection(schema)
			downSubAction = si.SubActionCreateCollection(schema)
		case si.SubActionTypeDropIndex:
			upSubAction = si.SubActionDropIndex(schema)
			downSubAction = si.SubActionCreateIndex(schema)
		case si.SubActionTypeDropField:
			upSubAction = si.SubActionDropField(schema)
			downSubAction = si.SubActionCreateField(schema)
		}

		// push sub actions
		if upSubAction != nil {
			upAction.SubActions = append(upAction.SubActions, *upSubAction)
			upActionMap[signedCollection.Key()] = upAction
		}
		if downSubAction != nil {
			downAction.SubActions = append(downAction.SubActions, *downSubAction)
			downActionMap[signedCollection.Key()] = downAction
		}
	}

	for _, signedCollection := range signedCollections {
		if signedCollection.IsIntersection {
			// actions for fields
			for _, signedField := range signedCollection.Fields {
				if signedField.Sign == SignPlus {
					// create
					fillActionMap(signedCollection, si.SubActionTypeCreateField)
				} else if signedField.Sign == SignMinus {
					// drop
					fillActionMap(signedCollection, si.SubActionTypeDropField)
				} else {
					// convert
					fillActionMap(signedCollection, si.SubActionTypeConvertField)
				}
			}
			// actions for indexes
			for _, signedIndex := range signedCollection.Indexes {
				if signedIndex.Sign == SignPlus {
					// create
					fillActionMap(signedCollection, si.SubActionTypeCreateIndex)
				} else {
					// drop
					fillActionMap(signedCollection, si.SubActionTypeDropIndex)
				}
			}
		} else if signedCollection.Sign == SignPlus {
			// create new collection
			fillActionMap(signedCollection, si.SubActionTypeCreateCollection)
		} else {
			// drop collection
			fillActionMap(signedCollection, si.SubActionTypeCreateCollection)
		}
	}

	upActions := []si.Action{}
	downActions := []si.Action{}

	sortSubActions := func(subActions []si.SubAction) []si.SubAction {
		// sort subActions
		sort.SliceStable(subActions, func(i, j int) bool {
			if subActions[i].IsUp() && !subActions[j].IsUp() {
				return false
			} else if !subActions[i].IsUp() && subActions[j].IsUp() {
				return true
			}

			// both i and j is Up
			// the order is collection -> field -> index
			return subActions[i].Type == si.SubActionTypeCreateCollection ||
				(subActions[i].Type == si.SubActionTypeCreateField && util.InListEq(subActions[j].Type, []si.SubActionType{
					si.SubActionTypeCreateField,
					si.SubActionTypeCreateIndex,
					si.SubActionTypeConvertField,
				})) ||
				(subActions[i].Type == si.SubActionTypeConvertField && util.InListEq(subActions[j].Type, []si.SubActionType{
					si.SubActionTypeCreateIndex,
					si.SubActionTypeConvertField,
				}))
		})

		return subActions
	}

	// sort both up actions and down actions
	for _, action := range upActionMap {
		action.SubActions = sortSubActions(action.SubActions)
		upActions = append(upActions, action)
	}

	for _, action := range downActionMap {
		action.SubActions = sortSubActions(action.SubActions)
		downActions = append(downActions, action)
	}

	return dt.NewPair(upActions, downActions)
}

// This constructs list of collection from a list of migrations
// `migrations` is a list of migration declarations in `[project dir]/mongr8/migration`
// for now, implement everything in this function
func GetCollectionFromMigrations(migrations []migrator.Migration) []collection.Collection {
	// assuming migrations are sorted based on ID
	// populate collection to a map of name and instance
	collections := map[string]collection.Collection{}

	var mergeFields func(a, b []collection.Field, c bool) []collection.Field
	mergeFields = func(origin, incoming []collection.Field, isRemove bool) []collection.Field {
		res := []collection.Field{}
		mergeDeeper := func(org, inc collection.Field) {
			if org.Spec().Type != inc.Spec().Type {
				if isRemove {
					panic("Field type for both origin and incoming must be same for removal")
				}

				// this should be field conversion, add the incoming
				res = append(res, inc)
			}

			if org.Spec().Type == field.TypeArray {
				merged := mergeFields(
					collection.FieldsFromSpecs(org.Spec().ArrayFields),
					collection.FieldsFromSpecs(inc.Spec().ArrayFields),
					isRemove,
				)
				// merged array fields not empty
				if len(merged) > 0 {
					arrayFields := collection.SpecsFromFields(merged)
					org.Spec().ArrayFields = &arrayFields
					res = append(res, org)
				}
			} else if org.Spec().Type == field.TypeObject {
				merged := mergeFields(
					collection.FieldsFromSpecs(org.Spec().Object),
					collection.FieldsFromSpecs(inc.Spec().Object),
					isRemove,
				)
				// merged object fields not empty
				if len(merged) > 0 {
					objectFields := collection.SpecsFromFields(merged)
					org.Spec().Object = &objectFields
					res = append(res, org)
				}
			}
		}

		// populate incoming and origin fields to map
		incFields := map[string]collection.Field{}
		for _, inc := range incoming {
			incFields[inc.Spec().Name] = inc
		}

		orgFields := map[string]collection.Field{}
		for _, org := range origin {
			orgFields[org.Spec().Name] = org
		}

		if isRemove {
			// check origin over incoming
			for _, org := range origin {
				inc, found := incFields[org.Spec().Name]
				if found {
					mergeDeeper(org, inc)
				} else {
					// skip, no need to remove
					res = append(res, org)
				}
			}
		} else {
			// add all origins those are not found in incoming
			for _, org := range origin {
				_, found := incFields[org.Spec().Name]
				if !found {
					res = append(res, org)
				}
			}
			
			// check incoming over origin
			for _, inc := range incoming {
				org, found := orgFields[inc.Spec().Name]
				if found {
					mergeDeeper(org, inc)
				} else {
					// if not found in origin, just add it
					res = append(res, inc)
				}
			}
		}

		return res
	}
	mergeToCollections := func(subAction *si.SubAction, migrationID string) {
		collectionName := subAction.ActionSchema.Collection.Spec().Name
		coll, ok := collections[collectionName]
		if subAction.Type == si.SubActionTypeCreateCollection {
			// add new collection
			collections[collectionName] = collection.NewCollection(
				subAction.ActionSchema.Collection,
				subAction.ActionSchema.Fields,
				subAction.ActionSchema.Indexes,
			)
		} else if ok {
			if subAction.Type == si.SubActionTypeDropCollection {
				// delete collection from map
				delete(collections, subAction.ActionSchema.Collection.Spec().Name)
			} else {
				// flag whether to add or remove
				// we consider a convert as an add
				// since, we only care of incoming type
				isRemove := util.InListEq(subAction.Type, []si.SubActionType{
					si.SubActionTypeDropField,
					si.SubActionTypeDropIndex,
				})
				newIndexes := []collection.Index{}
				if isRemove {
					// cross checking on indexes
					for _, currIndex := range coll.Indexes() {
						found := false
						for _, index := range subAction.ActionSchema.Indexes {
							if index.Spec().GetKey() == index.Spec().GetKey() {
								found = true
								break
							}
						}
						// skip to remove current index
						if found {
							continue
						}

						newIndexes = append(newIndexes, currIndex)
					}
				} else {
					// set indexes as initial state
					newIndexes = coll.Indexes()
					// cross checking on indexes
					for _, incomingIndex := range subAction.ActionSchema.Indexes {
						found := false
						for _, currIndex := range coll.Indexes() {
							if currIndex.Spec().GetKey() == incomingIndex.Spec().GetKey() {
								found = true
								break
							}
						}
						// add if not found
						if !found {
							newIndexes = append(newIndexes, incomingIndex)
						}
					}
				}

				// merge collections
				newFields := mergeFields(coll.Fields(), subAction.ActionSchema.Fields, isRemove)

				// finally update as new collection instance
				collections[collectionName] = collection.NewCollection(
					subAction.ActionSchema.Collection,
					newFields,
					newIndexes,
				)
			}
		} else {
			panic(fmt.Sprintf("Inconsitent migration found (%s) with action: %s\n", migrationID, subAction.Type.ToString()))
		}
	}

	for _, migration := range migrations {
		// we only care of UP actions
		for _, action := range migration.Up {
			for _, subAction := range action.SubActions {
				mergeToCollections(&subAction, migration.ID)
			}
		}
	}

	res := []collection.Collection{}
	for _, coll := range collections {
		res = append(res, coll)
	}

	return res
}
