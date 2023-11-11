package api_interpreter

import (
	"context"
	//"errors"
	"fmt"
	//"log"
	"reflect"
	//"strings"
	"testing"
	//"time"

	"internal/convert"
	dt "internal/data_type"
	"internal/test"
	"internal/util"

	"github.com/amirkode/go-mongr8/collection"
	"github.com/amirkode/go-mongr8/collection/field"
	"github.com/amirkode/go-mongr8/collection/index"
	"github.com/amirkode/go-mongr8/collection/metadata"
	"github.com/amirkode/go-mongr8/migration/migrator"
	"github.com/amirkode/go-mongr8/migration/translator/dictionary"
	si "github.com/amirkode/go-mongr8/migration/translator/mongodb/schema_interpreter"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	mim "github.com/ONSdigital/dp-mongodb-in-memory"
)

const MockDb = "mock-db"
const MockCollection = "mock_collection"

// NOTE: All these test are only
func getMockDatabase() (*mim.Server, *mongo.Database, *context.Context) {
	// server, err := memongo.StartWithOptions(&memongo.Options{
	// 	MongoVersion: "4.2.0",
	// 	MongodBin: ,
	// 	ShouldUseReplica: true,
	// })
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// mongoURI := server.URIWithRandomDB()
	// splitedDatabaseName := strings.Split(mongoURI, "/")
	// databaseName := splitedDatabaseName[len(splitedDatabaseName)-1]

	// uri := fmt.Sprintf("%s%s", mongoURI, "?retryWrites=false")
	// client, err := mongo.NewClient(options.Client().ApplyURI(uri))
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// ctx, _ := context.WithTimeout(context.Background(), 10 * time.Second)
	// err = client.Connect(ctx)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// db := client.Database(databaseName)

	// return server, db
	testCtx := context.Background()
	server, err := mim.StartWithOptions(testCtx, "5.0.2", mim.WithPort(27017))
	if err != nil {
		panic(err)
	}

	client, err := mongo.Connect(testCtx, options.Client().ApplyURI(server.URI()))
	if err != nil {
		panic(err)
	}

	//Use client as needed
	err = client.Ping(testCtx, nil)
	if err != nil {
		panic(err)
	}

	db := client.Database(MockDb)

	return server, db, &testCtx
}

func collectionExists(ctx context.Context, db *mongo.Database, name string) bool {
	names, err := db.ListCollectionNames(ctx, bson.D{{}})
	if err != nil {
		return false
	}

	for _, n := range names {
		if n == name {
			return true
		}
	}

	return false
}

func fieldsAreValid(ctx context.Context, db *mongo.Database, collectionName string, mustExist, mustNotExist []collection.Field) bool {
	var res []bson.M
	coll := db.Collection(collectionName)
	cursor, err := coll.Find(ctx, bson.D{})
	if err != nil {
		fmt.Println(err)
		return false
	}

	err = cursor.All(ctx, &res)
	if err != nil {
		fmt.Println(res)
		return false
	}

	if len(res) == 0 {
		return false
	}

	var validateFields func(inc interface{}, origin collection.Field) bool
	validateFields = func(inc interface{}, origin collection.Field) bool {
		// TODO: handle checking for special types, i.e: Geo JSON
		if util.NotInList(origin.Spec().Type, []field.FieldType{
			field.TypeString,
			field.TypeInt32,
			field.TypeInt64,
			field.TypeDouble,
			field.TypeBoolean,
			field.TypeArray,
			field.TypeObject,
			field.TypeTimestamp,
		}) {
			return false
		}

		if reflect.TypeOf(inc) == reflect.TypeOf(bson.A{}) {
			if origin.Spec().Type != field.TypeArray {
				return false
			}

			if len(inc.(bson.A)) == 0 || origin.Spec().ArrayFields == nil || len(*origin.Spec().ArrayFields) != 1 {
				return false
			}

			return validateFields(inc.(bson.A)[0], collection.FieldsFromSpecs(origin.Spec().ArrayFields)[0])
		} else if reflect.TypeOf(inc) == reflect.TypeOf(bson.M{}) {
			if origin.Spec().Type != field.TypeObject {
				return false
			}

			orgObj := *origin.Spec().Object
			incChildren := inc.(bson.M)
			orgChildren := map[string]collection.Field{}

			if origin.Spec().Object == nil || len(orgObj) != len(incChildren) {
				return false
			}

			for i := 0; i < len(orgObj); i++ {
				orgChildren[orgObj[i].Name] = collection.FieldsFromSpecs(&orgObj)[0]
			}

			// cross check inc over origin
			for key, value := range incChildren {
				org, ok := orgChildren[key]
				if !ok {
					return false
				}

				ok = validateFields(value, org)
				if !ok {
					return false
				}
			}
		}

		translatedOrg := dictionary.GetTranslatedField(origin)
		orgObj := translatedOrg.GetObject()
		item := orgObj[origin.Spec().Name]
		if reflect.TypeOf(item) != reflect.TypeOf(dictionary.ValueType{}) {
			return false
		}

		if reflect.TypeOf(convert.ConvertBsonPrimitiveToDefaultType(inc)) != reflect.TypeOf(item.(dictionary.ValueType).Value) {
			return false
		}

		return true
	}

	resMap := res[0]

	// fields must exist on res
	for _, value := range mustExist {
		inc, ok := resMap[value.Spec().Name]
		if !ok {
			return false
		}

		ok = validateFields(inc, value)
		if !ok {
			return false
		}
	}

	// fields must not exist on res
	for _, value := range mustNotExist {
		inc, ok := resMap[value.Spec().Name]
		if !ok {
			continue
		}

		ok = validateFields(inc, value)
		if ok {
			return false
		}
	}

	return true
}

func indexesAreValid(ctx context.Context, db *mongo.Database, collectionName string, mustExist, mustNotExist []collection.Index) bool {
	cursor, err := db.Collection(collectionName).Indexes().List(ctx)
	if err != nil {
		return false
	}

	var res []bson.M
	if err = cursor.All(ctx, &res); err != nil {
		return false
	}

	indexMap := map[string]bool{}
	for _, curr := range res {
		indexMap[curr["name"].(string)] = true
	}

	// indexes must exist
	for _, curr := range mustExist {
		key := curr.Spec().GetName()
		_, ok := indexMap[key]
		if !ok {
			return false
		}
	}

	// indexes must not exist
	for _, curr := range mustNotExist {
		key := curr.Spec().GetName()
		_, ok := indexMap[key]
		if ok {
			return false
		}
	}

	// var names []string
	// for _, index := range res {
	// 	names = append(names, index["name"].(string)) 
	// }

	return true
}

func setupCollection(ctx context.Context, db *mongo.Database) error {
	opt := options.CreateCollectionOptions{}
	err := db.CreateCollection(ctx, MockCollection, &opt)
	if err != nil {
		return err
	}

	subAction := si.SubAction{
		ActionSchema: si.SubActionSchema{
			Collection: metadata.InitMetadata("users"),
			Fields: []collection.Field{
				field.StringField("name"),
				field.Int32Field("age"),
			},
			Indexes: []collection.Index{
				index.CompoundIndex(
					index.Field("name", 1),
					index.Field("age", 1),
				),
			},
		},
	}

	// init few fields and indexes
	err = createField(ctx, db, MockCollection, subAction.GetFieldsBsonD(), false)
		if err != nil {
			return err
		}

	return createIndexes(ctx, db, MockCollection, subAction.GetIndexesBsonD())
}

// Test exeuctor functions for all available SubActionApis

func TestSubActionApiCreateCollection(t *testing.T) {
	server, db, ctx := getMockDatabase()
	defer server.Stop(*ctx)

	// Case 1: default
	case1SubActionApi := SubActionApiCreateCollection(dt.NewPair(
		migrator.Migration{},
		*si.SubActionCreateCollection(si.SubActionSchema{
			Collection: metadata.InitMetadata("users"),
			Fields: []collection.Field{
				field.StringField("name"),
				field.Int32Field("age"),
			},
			Indexes: []collection.Index{
				index.CompoundIndex(
					index.Field("name", 1),
					index.Field("age", 1),
				),
			},
		}),
	))
	case1Err := case1SubActionApi.Execute(*ctx, db)

	test.AssertTrue(t, case1Err == nil, "Case 1: Unexpected error")
	// check created collection
	test.AssertTrue(t, collectionExists(*ctx, db, 
		case1SubActionApi.SubAction.ActionSchema.Collection.Spec().Name,
	), "Case 1: Collection does not exist")
	test.AssertTrue(t, fieldsAreValid(*ctx, db, 
		case1SubActionApi.SubAction.ActionSchema.Collection.Spec().Name, 
		case1SubActionApi.SubAction.ActionSchema.Fields,
		[]collection.Field{},
	), "Case 1: Unexpected Fields")
	test.AssertTrue(t, indexesAreValid(*ctx, db, 
		case1SubActionApi.SubAction.ActionSchema.Collection.Spec().Name, 
		case1SubActionApi.SubAction.ActionSchema.Indexes,
		[]collection.Index{},
	), "Case 1: Unexpected Indexes")

	// TODO: add more cases
}

func TestSubActionApiCreateIndex(t *testing.T) {
	server, db, ctx := getMockDatabase()
	defer server.Stop(*ctx)

	err := setupCollection(*ctx, db)
	test.AssertTrue(t, err == nil, "Error while creating collection")
	
	// Case 1: create single field index on name
	case1SubActionApi := SubActionApiCreateIndex(dt.NewPair(
		migrator.Migration{},
		*si.SubActionCreateIndex(si.SubActionSchema{
			Collection: metadata.InitMetadata(MockCollection),
			Indexes: []collection.Index{
				index.SingleFieldIndex(
					index.Field("name", 1),
				),
			},
		}),
	))
	case1Err := case1SubActionApi.Execute(*ctx, db)

	test.AssertTrue(t, case1Err == nil, "Case 1: Unexpected error")
	// check created index
	test.AssertTrue(t, indexesAreValid(*ctx, db, 
		case1SubActionApi.SubAction.ActionSchema.Collection.Spec().Name, 
		case1SubActionApi.SubAction.ActionSchema.Indexes,
		[]collection.Index{},
	), "Case 1: Unexpected Indexes")

	// Case 2: create single field index on age
	case2SubActionApi := SubActionApiCreateIndex(dt.NewPair(
		migrator.Migration{},
		*si.SubActionCreateIndex(si.SubActionSchema{
			Collection: metadata.InitMetadata(MockCollection),
			Indexes: []collection.Index{
				index.SingleFieldIndex(
					index.Field("age", 1),
				),
			},
		}),
	))
	case2Err := case2SubActionApi.Execute(*ctx, db)

	test.AssertTrue(t, case2Err == nil, "Case 2: Unexpected error")
	// check created index
	test.AssertTrue(t, indexesAreValid(*ctx, db, 
		case2SubActionApi.SubAction.ActionSchema.Collection.Spec().Name, 
		case2SubActionApi.SubAction.ActionSchema.Indexes,
		[]collection.Index{},
	), "Case 2: Unexpected Indexes")

	// TODO: add more cases
}

func TestSubActionApiCreateField(t *testing.T) {
	server, db, ctx := getMockDatabase()
	defer server.Stop(*ctx)

	err := setupCollection(*ctx, db)
	test.AssertTrue(t, err == nil, "Error while creating collection")
	
	// Case 1: create a new timestamp field
	case1SubActionApi := SubActionApiCreateField(dt.NewPair(
		migrator.Migration{},
		*si.SubActionCreateField(si.SubActionSchema{
			Collection: metadata.InitMetadata(MockCollection),
			Fields: []collection.Field{
				field.TimestampField("created_at"),
			},
		}),
	))
	case1Err := case1SubActionApi.Execute(*ctx, db)

	test.AssertTrue(t, case1Err == nil, "Case 1: Unexpected error")
	// check created index
	test.AssertTrue(t, fieldsAreValid(*ctx, db, 
		case1SubActionApi.SubAction.ActionSchema.Collection.Spec().Name, 
		case1SubActionApi.SubAction.ActionSchema.Fields,
		[]collection.Field{},
	), "Case 1: Unexpected Fields")

	// TODO: add more cases
}

func TestSubActionApiConvertField(t *testing.T) {
	server, db, ctx := getMockDatabase()
	defer server.Stop(*ctx)

	err := setupCollection(*ctx, db)
	test.AssertTrue(t, err == nil, "Error while creating collection")
	
	// Case 1: convert age field to string
	case1SubActionApi := SubActionApiConvertField(dt.NewPair(
		migrator.Migration{},
		*si.SubActionConvertField(si.SubActionSchema{
			Collection: metadata.InitMetadata(MockCollection),
			Fields: []collection.Field{
				field.StringField("age"),
			},
			FieldConvertFrom: field.GetTypePointer(field.TypeInt32),
		}),
	))
	case1Err := case1SubActionApi.Execute(*ctx, db)

	test.AssertTrue(t, case1Err == nil, "Case 1: Unexpected error")
	// check created index
	test.AssertTrue(t, fieldsAreValid(*ctx, db, 
		case1SubActionApi.SubAction.ActionSchema.Collection.Spec().Name, 
		case1SubActionApi.SubAction.ActionSchema.Fields,
		[]collection.Field{},
	), "Case 1: Unexpected Fields")

	// TODO: add more cases
}

func TestSubActionApiDropCollection(t *testing.T) {
	server, db, ctx := getMockDatabase()
	defer server.Stop(*ctx)

	err := setupCollection(*ctx, db)
	test.AssertTrue(t, err == nil, "Error while creating collection")

	// Case 1: default
	case1SubActionApi := SubActionApiDropCollection(dt.NewPair(
		migrator.Migration{},
		*si.SubActionDropCollection(si.SubActionSchema{
			Collection: metadata.InitMetadata(MockCollection),
		}),
	))
	case1Err := case1SubActionApi.Execute(*ctx, db)
	
	test.AssertTrue(t, case1Err == nil, "Case 1: Unexpected error")
	test.AssertTrue(t, !collectionExists(*ctx, db, 
		MockCollection,
	), "Case 1: Collection exists")

	// TODO: add more cases
}

func TestSubActionApiDropIndex(t *testing.T) {
	server, db, ctx := getMockDatabase()
	defer server.Stop(*ctx)

	err := setupCollection(*ctx, db)
	test.AssertTrue(t, err == nil, "Error while creating collection")
	
	// Case 1: drop compound fields of name and age
	case1SubActionApi := SubActionApiDropIndex(dt.NewPair(
		migrator.Migration{},
		*si.SubActionDropIndex(si.SubActionSchema{
			Collection: metadata.InitMetadata(MockCollection),
			Indexes: []collection.Index{
				index.CompoundIndex(
					index.Field("name", 1),
					index.Field("age", 1),
				),
			},
		}),
	))
	case1Err := case1SubActionApi.Execute(*ctx, db)

	test.AssertTrue(t, case1Err == nil, "Case 1: Unexpected error")
	// check dropped index
	test.AssertTrue(t, indexesAreValid(*ctx, db, 
		case1SubActionApi.SubAction.ActionSchema.Collection.Spec().Name, 
		[]collection.Index{},
		case1SubActionApi.SubAction.ActionSchema.Indexes,
	), "Case 1: Unexpected Indexes")

	// TODO: add more cases
}

func TestSubActionApiDropField(t *testing.T) {
	server, db, ctx := getMockDatabase()
	defer server.Stop(*ctx)

	err := setupCollection(*ctx, db)
	test.AssertTrue(t, err == nil, "Error while creating collection")
	
	// Case 1: drop name field
	case1SubActionApi := SubActionApiDropField(dt.NewPair(
		migrator.Migration{},
		*si.SubActionCreateField(si.SubActionSchema{
			Collection: metadata.InitMetadata(MockCollection),
			Fields: []collection.Field{
				field.StringField("name"),
			},
		}),
	))
	case1Err := case1SubActionApi.Execute(*ctx, db)

	test.AssertTrue(t, case1Err == nil, "Case 1: Unexpected error")
	// check created index
	test.AssertTrue(t, fieldsAreValid(*ctx, db, 
		case1SubActionApi.SubAction.ActionSchema.Collection.Spec().Name, 
		[]collection.Field{},
		case1SubActionApi.SubAction.ActionSchema.Fields,
	), "Case 1: Unexpected Fields")

	// Case 2: drop age field
	case2SubActionApi := SubActionApiDropField(dt.NewPair(
		migrator.Migration{},
		*si.SubActionCreateField(si.SubActionSchema{
			Collection: metadata.InitMetadata(MockCollection),
			Fields: []collection.Field{
				field.Int32Field("name"),
			},
		}),
	))
	case2Err := case2SubActionApi.Execute(*ctx, db)

	test.AssertTrue(t, case2Err == nil, "Case 2: Unexpected error")
	// check created index
	test.AssertTrue(t, fieldsAreValid(*ctx, db, 
		case2SubActionApi.SubAction.ActionSchema.Collection.Spec().Name, 
		[]collection.Field{},
		case2SubActionApi.SubAction.ActionSchema.Fields,
	), "Case 2: Unexpected Fields")

	// TODO: add more cases
}