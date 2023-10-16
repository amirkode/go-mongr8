package sync_strategy

import (
	"github.com/amirkode/go-mongr8/collection"
	"github.com/amirkode/go-mongr8/migration/translator/mongodb/api_interpreter"
)

/*
this contains mechanism to sync schema across all sources
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

as shown on above example, possible actions for each usecase:
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

user options:
- Force Conversion: This will force conversion for undefined field type conversion

*/

func getSignedCollection(item collection.Collection) SignedCollection {
	signedMetadata := SignedMetadata{
		Metadata: item.Collection(),
		Sign: SignPlus,
	}
	signedFields := make([]SignedField, len(item.Fields()))
	for index, f := range item.Fields() {
		signedFields[index] = SignedField{
			Field: f,
			Sign: SignPlus,
		}
	}
	signedIndexes := make([]SignedIndex, len(item.Indexes()))
	for index, i := range item.Indexes() {
		signedIndexes[index] = SignedIndex{
			Index: i,
			Sign: SignPlus,
		}
	}

	return SignedCollection{
		Metadata: signedMetadata,
		Fields: signedFields,
		Indexes: signedIndexes,	
		Sign:       SignPlus,
	}
}

func getCollectionFromSubActions(subActions []api_interpreter.SubActionSchema) []collection.Collection {
	// implement something
	return []collection.Collection{}
}

func syncCollections(latestSchema []collection.Collection, existingSchema []collection.Collection) []SignedCollection {
	// convert latestSchema to signed collections
	source1 := make([]SignedCollection, len(latestSchema))
	for index, item := range latestSchema {
		source1[index] = getSignedCollection(item)
	}

	// convert existingSchema to signed collections
	source2 := make([]SignedCollection, len(existingSchema))
	for index, item := range existingSchema {
		source2[index] = getSignedCollection(item)
	}

	// find union
	return Union(source1, source2)
}
