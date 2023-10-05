package dictionary

import (
	"github.com/amirkode/go-mongr8/collection"
	"go.mongodb.org/mongo-driver/bson"
)

type (
	SchemaValidationIf interface {
		// serialize current schema to string
		toJsonString()
		// collection schema represented as a bson.M doc
		// the schema validation generated from collection.Fields()
		getCollectionDoc() bson.M
	}

	SchemaValidation struct {
		SchemaValidationIf
	}

	// translated field to bson.M doc
	// @ee field_impl.go for implementation
	TranslatedField struct {

	}

	// translated field to bson.M doc
	// @ee index_impl.go for implementation
	TranslatedIndex struct {
		
	}

	DictIf interface {
		initValidateFuncs()
		setSchemaValidation()
		// Translated Properties
		Fields() []TranslatedField
		Indexes() []TranslatedField
	}


	Dictionary struct {
		DictIf
		// raw collection data
		Collection       collection.Collection

		validateFuncs    []func() error

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

func (dict Dictionary) getPrimaryKey() {
	// pkField := dict.Collection.LookupField("_id")
	// if pkField != nil {
	// 	//
	// }
}

func (dict Dictionary) GetCollectionDoc() bson.M {
	res := bson.M{}	

	return res
}
