package api_interpreter

import (
	"context"
	"fmt"

	"github.com/amirkode/go-mongr8/collection/metadata"
	"go.mongodb.org/mongo-driver/mongo"
)

var ctx context.Context
var db *mongo.Database

// init subaction in MongoDB manner
// we assume all pre-validations have passed

func NewCreateCollectionSubAction(collectionName string, opts map[metadata.CollectionOption]interface{}) SubAction {
	stmt := fmt.Sprintf("%s := options.CreateCollectionOptions{}\n", VarNameCreateOptions)
	for key, value := range opts {
		if key == metadata.CollectionOptionCapped {
			stmt += fmt.Sprintf("%s.SetCapped(%s)\n", VarNameCreateOptions, anyToLiteralString(value))
		} else if key == metadata.CollectionOptionSize {
			stmt += fmt.Sprintf("%s.SetSizeInBytes(%s)\n", VarNameCreateOptions, anyToLiteralString(value))
		} else if key == metadata.CollectionOptionExpiredAfterSeconds {
			stmt += fmt.Sprintf("%s.SetExpireAfterSeconds(%s)\n", VarNameCreateOptions, anyToLiteralString(value))
		}
	}

	stmt += fmt.Sprintf(`if %s := %s.CreateCollection(%s, "%s", &%s); %s != nil { return %s }\n`,
		VarNameError,
		VarNameDatabase,
		VarNameContext,
		collectionName,
		VarNameCreateOptions,
		VarNameError,
		VarNameError,
	)

	return SubAction{
		GetStatement: func() string {
			return stmt
		},
		Type: SubActionCreateCollection,
	}
}

func NewInsertSingleSubAction(initCollection bool, collectionName string, payload map[string]interface{}) SubAction {
	stmt := ""
	// if it requires to initialize collection instance
	if initCollection {
		stmt += fmt.Sprintf(`%s := %s.Collection("%s")`, VarNameCollection, VarNameDatabase, collectionName) + "\n"
	}

	stmt += fmt.Sprintf(`count, _ := %s.CountDocuments(%s, bson.M{"_id": %s})`, VarNameCollection, VarNameContext, anyToLiteralString(payload["_id"])) + "\n"
	stmt += fmt.Sprintf(`if count > 0 { return fmt.Errorf("Collection already exists, cannot insert initial data") }`) + "\n"
	stmt += fmt.Sprintf(`if _, %s := %s.InsertOne(%s, %s); %s != nil { return %s }`,
		VarNameError,
		VarNameCollection,
		VarNameContext,
		toLiteralStringBsonMap(payload),
		VarNameError,
		VarNameError,
	)

	return SubAction{
		GetStatement: func() string {
			return stmt
		},
		Type: SubActionInsertOne,
	}
}
