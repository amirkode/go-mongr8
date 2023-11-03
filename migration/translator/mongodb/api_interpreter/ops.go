package api_interpreter

// translate sub action to operation with mongodb client interface

import (
	"context"
	"fmt"
	"reflect"

	dt "internal/data_type"

	"github.com/amirkode/go-mongr8/collection"
	"github.com/amirkode/go-mongr8/collection/field"
	"github.com/amirkode/go-mongr8/collection/metadata"
	si "github.com/amirkode/go-mongr8/migration/translator/mongodb/schema_interpreter"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func createField(ctx context.Context, db *mongo.Database, collName string, payload bson.D, update bool) error {
	collection := db.Collection(collName)
	// if it's not an update, then create/insert entire documents
	if !update {
		if _, err := collection.InsertOne(ctx, payload); err != nil {
			return err
		}

		return nil
	}

	// set field expects 1 path
	var setPayload func(curr interface{}, path string) bson.M
	setPayload = func(curr interface{}, path string) bson.M {
		if reflect.TypeOf(curr) == reflect.TypeOf(bson.D{}) {
			d := curr.(bson.D)
			return setPayload(d[0].Value, fmt.Sprintf("%s.%s", path, d[0].Key))
		} else if reflect.TypeOf(curr) == reflect.TypeOf(bson.A{}) {
			a := curr.(bson.A)
			return setPayload(a[0], fmt.Sprintf("%s.$[]", path))
		}

		return bson.M{
			path: curr,
		}
	}

	updatePayload := bson.M{
		"$set": setPayload(payload, ""),
	}

	_, err := collection.UpdateMany(ctx, bson.M{}, updatePayload)

	return err
}


// TODO: complete this + test
func convertField(ctx context.Context, db *mongo.Database, collName string, to collection.Field, from field.FieldType) error {
	// depth as suffix of map alias to maintain the uniqueness of the alias
	depth := 0 
	updatePayload := bson.M{
		"$set": convertSetPayload(to, "", from, &depth),
	}

	collection := db.Collection(collName)
	_, err := collection.UpdateMany(ctx, bson.M{}, updatePayload)

	return err
}

func createIndexes(ctx context.Context, db *mongo.Database, collName string, indexes []dt.Pair[bson.D, bson.D]) error {
	collection := db.Collection(collName)
	for _, idx := range indexes {
		keys := idx.First
		rules := idx.Second
		opt := options.Index()
		// init options
		for _, rule := range rules {
			switch rule.Key {
			case "partialFilterExpression":
				opt = opt.SetPartialFilterExpression(rule.Value)
			case "sparse":
				opt = opt.SetSparse(true)
			case "collation":
				collation := options.Collation{}
				// based on https://www.mongodb.com/docs/manual/reference/collation/
				for _, c := range rule.Value.(bson.D) {
					switch c.Key {
					case "locale":
						collation.Locale = c.Value.(string)
					case "caseLevel":
						collation.CaseLevel = c.Value.(bool)
					case "caseFirst":
						collation.CaseFirst = c.Value.(string)
					case "strength":
						collation.Strength = c.Value.(int)
					case "numericOrdering":
						collation.NumericOrdering = c.Value.(bool)
					case "alternate":
						collation.Alternate = c.Value.(string)
					case "maxVariable":
						collation.MaxVariable = c.Value.(string)
					case "backwards":
						collation.Backwards = c.Value.(bool)
					}
				}

				opt = opt.SetCollation(&collation)
			case "unique":
				opt = opt.SetUnique(true)
			}
		}

		indexModel := mongo.IndexModel{
			Keys:    keys,
			Options: opt,
		}

		_, err := collection.Indexes().CreateOne(ctx, indexModel)
		if err != nil {
			return err
		}
	}

	return nil
}

func SubActionApiCreateCollection(subAction dt.Pair[string, si.SubAction]) SubActionApi {
	collectionName := subAction.Second.ActionSchema.Collection.Spec().Name
	exec := func(ctx context.Context, db *mongo.Database) error {
		opt := options.CreateCollectionOptions{}
		schemaOption := subAction.Second.ActionSchema.Collection.Spec().Options
		if schemaOption != nil {
			_, capped := (*schemaOption)[metadata.CollectionOptionCapped]
			if capped {
				cappedSize, _ := (*schemaOption)[metadata.CollectionOptionCappedSize]

				opt.SetCapped(true)
				opt.SetSizeInBytes(cappedSize.(int64))
			}

			ttl, useTTL := (*schemaOption)[metadata.CollectionOptionExpiredAfterSeconds]
			if useTTL {
				opt.SetExpireAfterSeconds(ttl.(int64))
			}
		}

		err := db.CreateCollection(ctx, collectionName, &opt)
		if err != nil {
			return err
		}

		err = createField(ctx, db, collectionName, subAction.Second.GetFieldsBsonD(), false)
		if err != nil {
			return err
		}

		return createIndexes(ctx, db, collectionName, subAction.Second.GetIndexesBsonD())
	}

	return SubActionApi{
		MigrationID: subAction.First,
		SubAction:   subAction.Second,
		Execute:     exec,
	}
}

func SubActionApiCreateIndex(subAction dt.Pair[string, si.SubAction]) SubActionApi {
	collectionName := subAction.Second.ActionSchema.Collection.Spec().Name
	exec := func(ctx context.Context, db *mongo.Database) error {
		return createIndexes(ctx, db, collectionName, subAction.Second.GetIndexesBsonD())
	}

	return SubActionApi{
		MigrationID: subAction.First,
		SubAction:   subAction.Second,
		Execute:     exec,
	}
}

func SubActionApiCreateField(subAction dt.Pair[string, si.SubAction]) SubActionApi {
	collectionName := subAction.Second.ActionSchema.Collection.Spec().Name
	exec := func(ctx context.Context, db *mongo.Database) error {
		return createField(ctx, db, collectionName, subAction.Second.GetFieldsBsonD(), true)
	}

	return SubActionApi{
		MigrationID: subAction.First,
		SubAction:   subAction.Second,
		Execute:     exec,
	}
}

func SubActionApiConvertField(subAction dt.Pair[string, si.SubAction]) SubActionApi {
	collectionName := subAction.Second.ActionSchema.Collection.Spec().Name
	exec := func(ctx context.Context, db *mongo.Database) error {
		return convertField(ctx, db, collectionName, subAction.Second.ActionSchema.Fields[0], *subAction.Second.ActionSchema.FieldConvertFrom)
	}

	return SubActionApi{
		MigrationID: subAction.First,
		SubAction:   subAction.Second,
		Execute:     exec,
	}
}
