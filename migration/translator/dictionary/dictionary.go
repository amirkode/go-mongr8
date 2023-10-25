package dictionary

import (
	"github.com/amirkode/go-mongr8/collection"
	"github.com/amirkode/go-mongr8/collection/metadata"
)

type (
	SchemaValidationIf interface {
		// serialize current schema to string
		toJsonString()
		// collection schema represented as a map doc
		// the schema validation generated from collection.Fields()
		getCollectionDoc() map[string]interface{}
	}

	SchemaValidation struct {
		SchemaValidationIf
	}

	// translated field to bson.M doc
	// @ee field_impl.go for implementation
	TranslatedFieldIf interface {
		// still not find the proper usecase
		// for GetArray()
		GetArray() []interface{}

		// get object of current field
		GetObject() map[string]interface{}

		// get value of current field
		getItem() interface{}
	}

	// translated field to bson.M doc
	// @ee index_impl.go for implementation
	TranslatedIndexIf interface {
		// get object of indexes
		getObject() map[string]interface{}
		// get rules
		getRules() *map[string]interface{}
	}

	DictIf interface {
		initValidateFuncs()
		setSchemaValidation()
		// Translated Properties
		GetPrimaryKey()
		GetOptions() *map[metadata.CollectionOption]interface{}
		GetDocument() map[string]interface{}
	}

	Dictionary struct {
		DictIf
		// raw collection data
		Collection collection.Collection

		validateFuncs []func() error

		// if schema validation exists, should be added to validateFuncs on setSchemaValidation
		// this is also used to describe current schma validation
		schemeValidation *SchemaValidation
	}
)

func (dict Dictionary) Translate() {
	dict.initValidateFuncs()
}

func (dict Dictionary) validate() error {
	for _, v := range dict.validateFuncs {
		if err := v(); err != nil {
			return err
		}
	}

	return nil
}

func (dict Dictionary) setValidateFuncs(funcs ...func() error) {
	for _, f := range funcs {
		dict.validateFuncs = append(dict.validateFuncs, f)
	}
}

func (dict Dictionary) GetPrimaryKey() {
	// pkField := dict.Collection.LookupField("_id")
	// if pkField != nil {
	// 	//
	// }
}

func (dict Dictionary) GetCollectionDoc() map[string]interface{} {
	res := map[string]interface{}{}

	return res
}
