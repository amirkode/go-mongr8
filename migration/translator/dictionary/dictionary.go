package dictionary

import (
	"github.com/amirkode/go-mongr8/collection"
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

	ValidationIf interface {
		Validate()
	}

	Validation struct {
		ValidationIf
		// raw collection data
		Collections []collection.Collection

		validationFuncs []func() error

		// // if schema validation exists, should be added to validateFuncs on setSchemaValidation
		// // this is also used to describe current schma validation
		// schemeValidation *SchemaValidation
	}
)
