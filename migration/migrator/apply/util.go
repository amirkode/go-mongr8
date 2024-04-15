/*
Copyright (c) 2023-present the go-mongr8 Authors and Contributors
[@see Authors file]

Licensed under the MIT License
(https://opensource.org/licenses/MIT)
*/
package apply

import (
	"context"
	"time"

	"github.com/amirkode/go-mongr8/internal/constant"

	"github.com/amirkode/go-mongr8/migration/common"
	"github.com/amirkode/go-mongr8/migration/migrator"
	ai "github.com/amirkode/go-mongr8/migration/translator/mongodb/api_interpreter"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type MigrationHistory struct {
	MigrationID string    `bson:"_id"`
	Desc        string    `bson:"desc"`
	MigratedAt  time.Time `bson:"migrated_at"`
}

func getLatestMigrationID(ctx context.Context, db *mongo.Database) (*string, error) {
	res := constant.MinTimevalue().Format("20060102_150405")
	coll := db.Collection(common.MigrationHistoryCollection)
	cursor, err := coll.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	for cursor.Next(ctx) {
		var history MigrationHistory
		err = cursor.Decode(&history)
		if err != nil {
			return nil, err
		}

		if history.MigrationID > res {
			res = history.MigrationID
		}
	}

	return &res, nil
}

func updateMigrationHistory(migrations []migrator.Migration, ctx context.Context, db *mongo.Database) error {
	migratedAt := time.Now()
	payload := []interface{}{}
	for _, m := range migrations {
		payload = append(payload, MigrationHistory{
			MigrationID: m.ID,
			Desc:        m.Desc,
			MigratedAt:  migratedAt,
		})
	}

	coll := db.Collection(common.MigrationHistoryCollection)
	_, err := coll.InsertMany(ctx, payload)

	return err
}

func filterSubActionApi(apis []ai.SubActionApi, ctx context.Context, db *mongo.Database) (*[]ai.SubActionApi, error) {
	res := []ai.SubActionApi{}
	latestMigrationID, err := getLatestMigrationID(ctx, db)
	if err != nil {
		return nil, err
	}

	for _, api := range apis {
		if api.Migration.ID > *latestMigrationID {
			res = append(res, api)
		}
	}

	return &res, nil
}
