package translator

import (
	"fmt"

	"github.com/amirkode/go-mongr8/collection"
	"github.com/amirkode/go-mongr8/migration/translator/dictionary"
	// "github.com/amirkode/go-mongr8/migration/translator/collection_loader"
	"github.com/amirkode/go-mongr8/migration/translator/mongodb/api_interpreter"
)

type (

	ProcessorIf interface {
		initDictionary(collections []collection.Collection)
		GetTranslateDictionaries() (*[]dictionary.Dictionary, error)
		Generate(collections []collection.Collection, migrationSubActionSchemas []api_interpreter.SubActionSchema)
		Consolidate(collections []collection.Collection, dbCollections []collection.Collection, migrationSubActionSchemas []api_interpreter.SubActionSchema)
	}

	Processor struct {
		ProcessorIf,
		Dictionaries []dictionary.Dictionary
		Init bool
	}
)

func (p Processor) Generate(collections []collection.Collection, migrationSubActionSchemas []api_interpreter.SubActionSchema) {
	//p.initDictionary(collections)
	// TODO: something to find intersection of user-defined collections, migration-generated sub action schemas
	
}

func (p Processor)  Consolidate(collections []collection.Collection, dbCollections []collection.Collection, migrationSubActionSchemas []api_interpreter.SubActionSchema) {
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

func NewProcessor() ProcessorIf {
	return Processor{
		Init: true,
	}	
}
