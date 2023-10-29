package translator

import (
	"context"

	dt "internal/data_type"

	"github.com/amirkode/go-mongr8/collection"
	"github.com/amirkode/go-mongr8/migration/migrator"
	"github.com/amirkode/go-mongr8/migration/translator/dictionary"
	si "github.com/amirkode/go-mongr8/migration/translator/mongodb/schema_interpreter"
	"github.com/amirkode/go-mongr8/migration/translator/sync_strategy"
)

type (
	ProcessorIf interface {
		validateCollection(collections []collection.Collection, panic bool) error
		Generate(collections []collection.Collection, migrations []migrator.Migration) dt.Pair[[]si.Action, []si.Action]
		Consolidate(collections []collection.Collection, dbCollections []collection.Collection, migrations []migrator.Migration)
	}

	Processor struct {
		ProcessorIf,
		Ctx *context.Context
		Init bool
	}
)

func (p Processor) validateCollection(collections []collection.Collection, panicOnError bool) error {
	validation := dictionary.Validation{
		Collections: collections,
	}
	err := validation.Validate()
	if err != nil {
		if panicOnError {
			panic(err.Error())
		}

		return err
	}

	return nil
}

func (p Processor) Generate(collections []collection.Collection, migrations []migrator.Migration) dt.Pair[[]si.Action, []si.Action] {
	// validate incoming collections
	p.validateCollection(collections, true)
	collectionsFromMigrations := sync_strategy.GetCollectionFromMigrations(migrations)
	// validate exisiting collections
	err := p.validateCollection(collections, false)
	if err != nil {
		// TODO: should it automatically consolidate?
		panic(err)
	}

	return sync_strategy.GetActions(collections, collectionsFromMigrations)
}

func (p Processor) Consolidate(collections []collection.Collection, dbCollections []collection.Collection, migrations []migrator.Migration) {
	// TODO: something to find consolidation resulting intersection of user-defined collections, db collections, migration-generated sub action schemas
}

func NewProcessor(ctx *context.Context) ProcessorIf {
	return Processor{
		Ctx:  ctx,
		Init: true,
	}
}
