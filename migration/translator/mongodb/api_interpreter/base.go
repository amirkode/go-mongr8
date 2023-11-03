package api_interpreter

import (
	"context"

	dt "internal/data_type"

	"github.com/amirkode/go-mongr8/collection"
	si "github.com/amirkode/go-mongr8/migration/translator/mongodb/schema_interpreter"

	"go.mongodb.org/mongo-driver/mongo"
)

type (
	SubActionApi struct {
		MigrationID string
		// TODO: decide whether SubAction is always attached to SubActionApi (?)
		SubAction si.SubAction
		Execute   func(ctx context.Context, db *mongo.Database) error
	}
)

// This returns the list of SubActionApi(s)
// `subActions` retrieved from migration files
// `dbSchemas` is current schema from database formatted in Collection manner
func GetSubActionApis(subActions []dt.Pair[string, si.SubAction], dbSchemas []collection.Collection) []SubActionApi {
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
			// TODO: complement this
		}
	}

	return res
}
