package translator

import (
	"context"
	"fmt"

	dt "internal/data_type"

	"github.com/amirkode/go-mongr8/collection"
	"github.com/amirkode/go-mongr8/migration/migrator"
	"github.com/amirkode/go-mongr8/migration/translator/dictionary"
	"github.com/amirkode/go-mongr8/migration/translator/sync_strategy"
	si "github.com/amirkode/go-mongr8/migration/translator/mongodb/schema_interpreter"
)

type (

	ProcessorIf interface {
		ValidateCollection(collections []collection.Collection) error
		Generate(collections []collection.Collection, migrations []migrator.Migration) dt.Pair[[]si.Action, []si.Action]
		Consolidate(collections []collection.Collection, dbCollections []collection.Collection, migrations []migrator.Migration)
	}

	Processor struct {
		ProcessorIf,
		Ctx context.Context
		Dictionaries []dictionary.Dictionary
		Init bool
	}
)

func (p Processor) ValidateCollection(collections []collection.Collection) error {
	// implement something
	return nil
}

func (p Processor) Generate(collections []collection.Collection, migrations []migrator.Migration) dt.Pair[[]si.Action, []si.Action] {
	collectionsFromMigrations := sync_strategy.GetCollectionFromMigrations(migrations)
	// TODO: validate collection
	return sync_strategy.GetActions(collections, collectionsFromMigrations)
}

func (p Processor)  Consolidate(collections []collection.Collection, dbCollections []collection.Collection, migrations []migrator.Migration) {
	// TODO: something to find consolidation resulting intersection of user-defined collections, db collections, migration-generated sub action schemas
}

func (p Processor) initDictionary(collections []collection.Collection) {
	p.Dictionaries = make([]dictionary.Dictionary, len(collections))
	for index, collection := range collections {
		dict := dictionary.Dictionary{
			Collection: collection,
		}
		// now translate to MongoDB manners
		dict.Translate()
		// set dictionary to current index array
		p.Dictionaries[index] = dict
	}
}

// this would be consumed any migrator operation
func (p Processor) GetTranslateDictionaries() (*[]dictionary.Dictionary, error) {
	if !p.Init {
		return nil, fmt.Errorf("Translation processor has not been initialized.")
	}

	return &p.Dictionaries, nil
}

func NewProcessor(ctx context.Context) ProcessorIf {
	return Processor{
		Ctx: ctx,
		Init: true,
	}	
}
