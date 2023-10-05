package migration

import (
	"time"

	"github.com/amirkode/go-mongr8/collection"
	"github.com/amirkode/go-mongr8/migration/migrator/generate"
	"github.com/amirkode/go-mongr8/migration/translator"
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
		date string
	}
)

func NewMigration() Cmd {
	return Migration{
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
	translator.Proceed(collections)
	// get translated dictionary
	translatedDictionaries, err := translator.GetTranslateDictionaries()
	if err != nil {
		return err
	}

	return generate.Run(*translatedDictionaries)
}

