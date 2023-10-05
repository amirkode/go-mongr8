package translator

import (
	"fmt"

	"github.com/amirkode/go-mongr8/collection"
	"github.com/amirkode/go-mongr8/migration/translator/dictionary"
	// "github.com/amirkode/go-mongr8/migration/translator/collection_loader"
)

type (
	ProcessorIf interface {
		lnitDictionary()
	}

	Processor struct {
		ProcessorIf,
		Dictionaries []dictionary.Dictionary
		Init bool
	}
)

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

// global processor
var processor Processor

func Proceed(collections []collection.Collection) {
	processor = Processor{
		Init: true,
	}
	processor.initDictionary(collections)
}

// this would be consumed any migrator operation
func GetTranslateDictionaries() (*[]dictionary.Dictionary, error) {
	if !processor.Init {
		return nil, fmt.Errorf("Translation processor has not been initialized.")
	}

	return &processor.Dictionaries, nil
}