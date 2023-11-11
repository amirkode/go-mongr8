/*
Copyright (c) 2023 the go-mongr8 Authors and Contributors
[@see Authors file]

Licensed under the MIT License
(https://opensource.org/licenses/MIT)
*/
package api_interpreter

import (
	"context"

	dt "github.com/amirkode/go-mongr8/internal/data_type"
	"github.com/amirkode/go-mongr8/collection"
	"github.com/amirkode/go-mongr8/migration/migrator"
	si "github.com/amirkode/go-mongr8/migration/translator/mongodb/schema_interpreter"

	"go.mongodb.org/mongo-driver/mongo"
)

type (
	SubActionApi struct {
		Migration migrator.Migration
		// TODO: decide whether SubAction is always attached to SubActionApi (?), since not direct usage required
		SubAction si.SubAction
		Execute   func(ctx context.Context, db *mongo.Database) error
	}
)

// This returns the list of SubActionApi(s)
// `subActions` retrieved from migration files
// `dbSchemas` is current schema from database formatted in Collection manner
func GetSubActionApis(subActions []dt.Pair[migrator.Migration, si.SubAction], dbSchemas []collection.Collection) []SubActionApi {
	// TODO: for now, we are assuming that all subactions are valid.
	// in fact, schema might changed dynamically on database in real case
	// we will have to implement the validation eventually
	res := []SubActionApi{}
	for _, subAction := range subActions {
		switch subAction.Second.Type {
		case si.SubActionTypeCreateCollection:
			res = append(res, SubActionApiCreateCollection(subAction))
		case si.SubActionTypeCreateIndex:
			res = append(res, SubActionApiCreateIndex(subAction))
		case si.SubActionTypeCreateField:
			res = append(res, SubActionApiCreateField(subAction))
		case si.SubActionTypeConvertField:
			res = append(res, SubActionApiConvertField(subAction))
		case si.SubActionTypeDropIndex:
			res = append(res, SubActionApiDropCollection(subAction))
		case si.SubActionTypeDropField:
			res = append(res, SubActionApiDropField(subAction))
		}
	}

	return res
}
