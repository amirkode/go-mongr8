/*
Copyright (c) 2023 the go-mongr8 Authors and Contributors
[@see Authors file]

Licensed under the MIT License
(https://opensource.org/licenses/MIT)
*/
package apply

import (
	"context"
	"fmt"

	"github.com/amirkode/go-mongr8/migration/migrator"
	"github.com/amirkode/go-mongr8/migration/option"
	ai "github.com/amirkode/go-mongr8/migration/translator/mongodb/api_interpreter"

	"go.mongodb.org/mongo-driver/mongo"
)

func execSubActions(ctx context.Context, db *mongo.Database, apis []ai.SubActionApi) error {
	filteredApis, err := filterSubActionApi(apis, ctx, db)
	if err != nil {
		return err
	}

	if len(*filteredApis) > 0 {
		migrations := []migrator.Migration{}
		added := map[string]bool{}
		for _, api := range *filteredApis {
			// execute the action api synchronously
			err := api.Execute(ctx, db)
			if err != nil {
				return err
			}

			_, exists := added[api.Migration.ID]
			if !exists {
				migrations = append(migrations, api.Migration)
				added[api.Migration.ID] = true
			}
		}

		// update migration history
		err = updateMigrationHistory(migrations, ctx, db)
		if err != nil {
			return err
		}

		fmt.Printf("All Migration files has been migrated with IDs %s..%s", 
			(*filteredApis)[0].Migration.ID, 
			(*filteredApis)[len(*filteredApis) - 1].Migration.ID,
		)
	} else {
		fmt.Printf("Nothing to migrate.")
	}

	return nil
}

func Run(ctx *context.Context, db *mongo.Database, apis []ai.SubActionApi) error {
	useTransaction := option.GetMigrationOptionFromContext(ctx).UseTransaction
	if !useTransaction {
		// executes everything with individually
		err := execSubActions(*ctx, db, apis)
		if err != nil {
			return err
		}

		return nil
	}

	// ctx, cancel := context.WithTimeout(context.Background(), 30 * time.Second)
	// defer cancel()

	session, err := db.Client().StartSession()
	if err != nil {
		return err
	}
	defer session.EndSession(*ctx)

	// bind everything in a transaction
	return mongo.WithSession(*ctx, session, func(sc mongo.SessionContext) error {
		if err := sc.StartTransaction(); err != nil {
			return err
		}

		// execute sub actions
		err := execSubActions(sc, db, apis)
		if err != nil {
			// rollback
			if rErr := sc.AbortTransaction(*ctx); rErr != nil {
				return fmt.Errorf("error rolling back transaction: %v (caused by original error: %v)", rErr, err)
			}
			return err
		}

		return nil
	})
}
