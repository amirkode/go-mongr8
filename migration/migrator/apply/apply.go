package apply

import (
	"context"
	"fmt"
	"time"
	
	ai "github.com/amirkode/go-mongr8/migration/translator/mongodb/api_interpreter"

	"go.mongodb.org/mongo-driver/mongo"
)

func Run(db *mongo.Database, apis []ai.SubActionApi) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	session, err := db.Client().StartSession()
	if err != nil {
		return err
	}
	defer session.EndSession(ctx)

	// bind everything in a transaction
	return mongo.WithSession(ctx, session, func(sc mongo.SessionContext) error {
		if err := sc.StartTransaction(); err != nil {
			return err
		}

		maxMigrationID := ""
		filteredApis := filterSubActionApi(apis)
		for _, api := range filteredApis {
			// execute the action api synchronously
			err := api.Execute(sc, db)
			if err != nil {
				// rollback
				if rErr := sc.AbortTransaction(ctx); rErr != nil {
					return fmt.Errorf("error rolling back transaction: %v (caused by original error: %v)", rErr, err)
				}
				return err
			}

			if maxMigrationID < api.MigrationID {
				maxMigrationID = api.MigrationID
			}
		}

		// TODO: something to update migration history

		if err := sc.CommitTransaction(ctx); err != nil {
			return err
		}

		return nil	
	})
}