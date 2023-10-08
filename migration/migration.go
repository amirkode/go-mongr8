package migration

import (
	"time"

	"github.com/amirkode/go-mongr8/collection"
	"github.com/amirkode/go-mongr8/migration/migrator/generate"
	"github.com/amirkode/go-mongr8/migration/translator"

	"go.mongodb.org/mongo-driver/mongo"
)

const (
	Mongr8Version = "v0.0.1-alpha"
)

type (
	Cmd interface {
		ApplyMigration() error
		ConsolidateMigration(collections []collection.Collection) error
		GenerateMigration(collections []collection.Collection) error
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

func (Migration) ConsolidateMigration(collections []collection.Collection) error {
	return nil
}

func (Migration) GenerateMigration(collections []collection.Collection) error {
	processor := translator.NewProcessor()
	processor.Proceed(collections)
	// get translated dictionary
	translatedDictionaries, err := processor.GetTranslateDictionaries()
	if err != nil {
		return err
	}

	return generate.Run(*translatedDictionaries)
}
