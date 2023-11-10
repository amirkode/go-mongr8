package migration

import (
	"context"
	"time"

	"github.com/amirkode/go-mongr8/collection"
	"github.com/amirkode/go-mongr8/migration/migrator"
	"github.com/amirkode/go-mongr8/migration/migrator/apply"
	"github.com/amirkode/go-mongr8/migration/migrator/generate"
	"github.com/amirkode/go-mongr8/migration/migrator/loader"
	"github.com/amirkode/go-mongr8/migration/translator"

	"go.mongodb.org/mongo-driver/mongo"
)

type (
	Cmd interface {
		ApplyMigration(migrations []migrator.Migration) error
		ConsolidateMigration(collections []collection.Collection, migrations []migrator.Migration) error
		GenerateMigration(collections []collection.Collection, migrations []migrator.Migration) error
	}

	Migration struct {
		Cmd
		ctx  *context.Context
		db   *mongo.Database
		date string
	}
)

func NewMigration(ctx *context.Context, db *mongo.Database) Cmd {
	return &Migration{
		ctx:  ctx,
		db:   db,
		date: time.Now().Format("2006-01-02"),
	}
}

func (m *Migration) ApplyMigration(migrations []migrator.Migration) error {
	dbSchemas := loader.GetSchemaFromDB()
	processor := translator.NewProcessor(m.ctx)
	apis := processor.GetApi(migrations, dbSchemas)

	return apply.Run(m.ctx, m.db, apis)
}

func (m *Migration) ConsolidateMigration(collections []collection.Collection, migrations []migrator.Migration) error {
	dbSchemas := loader.GetSchemaFromDB()
	processor := translator.NewProcessor(m.ctx)
	processor.Consolidate(collections, dbSchemas, migrations)
	return nil
}

func (m *Migration) GenerateMigration(collections []collection.Collection, migrations []migrator.Migration) error {
	processor := translator.NewProcessor(m.ctx)
	actions := processor.Generate(collections, migrations)

	return generate.Run(m.ctx, actions)
}
