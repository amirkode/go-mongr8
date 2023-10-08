package translator

import (
	"fmt"

	"github.com/amirkode/go-mongr8/collection"
	"github.com/amirkode/go-mongr8/migration/translator/dictionary"
	// "github.com/amirkode/go-mongr8/migration/translator/collection_loader"
)

type (

	ProcessorIf interface {
		initDictionary(collections []collection.Collection)
		GetTranslateDictionaries() (*[]dictionary.Dictionary, error)
		Proceed(collections []collection.Collection)
	}

	Processor struct {
		ProcessorIf,
		Dictionaries []dictionary.Dictionary
		Init bool
	}
)

func (p Processor) Proceed(collections []collection.Collection) {
	p.initDictionary(collections)
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
