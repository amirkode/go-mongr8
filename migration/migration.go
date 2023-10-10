package migration

import (
	"time"

	"github.com/amirkode/go-mongr8/collection"
	"github.com/amirkode/go-mongr8/migration/migrator/generate"
	"github.com/amirkode/go-mongr8/migration/migrator/loader"
	"github.com/amirkode/go-mongr8/migration/translator"
	"github.com/amirkode/go-mongr8/migration/translator/mongodb/api_interpreter"

	"go.mongodb.org/mongo-driver/mongo"
)

const (
	Mongr8Version = "v0.0.1-alpha"
)

type (
	Cmd interface {
		ApplyMigration() error
		ConsolidateMigration(collections []collection.Collection, migrationSubActionSchemas []api_interpreter.SubActionSchema) error
		GenerateMigration(collections []collection.Collection, migrationSubActionSchemas []api_interpreter.SubActionSchema) error
	}

	Migration struct {
		Cmd
		db   *mongo.Database
		date string
	}
)

func NewMigration(db *mongo.Database) Cmd {
	return Migration{
		db:   db,
		date: time.Now().Format("2006-01-02"),
	}
}

func (Migration) ApplyMigration() error {
	return nil
}

func (Migration) ConsolidateMigration(collections []collection.Collection, migrationSubActionSchemas []api_interpreter.SubActionSchema) error {
	dbSchemas := loader.GetSchemaFromDB()
	processor := translator.NewProcessor()
	processor.Consolidate(collections, dbSchemas, migrationSubActionSchemas)
	return nil
}

func (Migration) GenerateMigration(collections []collection.Collection, migrationSubActionSchemas []api_interpreter.SubActionSchema) error {
	processor := translator.NewProcessor()
	processor.Generate(collections, migrationSubActionSchemas)
	// get translated dictionary
	translatedDictionaries, err := processor.GetTranslateDictionaries()
	if err != nil {
		return err
	}

	return generate.Run(*translatedDictionaries)
}
