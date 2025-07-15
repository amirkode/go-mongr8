/*
Copyright (c) 2023-present the go-mongr8 Authors and Contributors
[@see Authors file]

Licensed under the MIT License
(https://opensource.org/licenses/MIT)
*/
package dictionary

import (
	"fmt"
	"strings"

	"github.com/amirkode/go-mongr8/collection"
	"github.com/amirkode/go-mongr8/collection/field"
	"github.com/amirkode/go-mongr8/collection/index"
	"github.com/amirkode/go-mongr8/internal/util"
	"github.com/amirkode/go-mongr8/migration/common"
)

func (v *Validation) Validate() error {
	v.initValidationFunctions()

	for _, validate := range v.validationFuncs {
		if err := validate(); err != nil {
			return err
		}
	}

	return nil
}

func (v *Validation) initValidationFunctions() {
	v.validationFuncs = []func() error{
		func() error {
			return validateCollections(v.Collections)
		},
	}

	for _, coll := range v.Collections {
		v.validationFuncs = append(v.validationFuncs, func() error {
			return validateID(coll.Collection().Spec().Name, coll.Fields())
		})
		v.validationFuncs = append(v.validationFuncs, func() error {
			return validateFields(coll.Collection().Spec().Name, coll.Fields())
		})
		v.validationFuncs = append(v.validationFuncs, func() error {
			return validateIndexes(coll.Collection().Spec().Name, coll.Fields(), coll.Indexes())
		})
	}
}

func validateCollections(collections []collection.Collection) error {
	// validate duplicate name
	dup := map[string]bool{}
	for _, coll := range collections {
		// collection name cannot be the same as default migration history collection name
		if coll.Collection().Spec().Name == common.MigrationHistoryCollection {
			return fmt.Errorf("Collection name cannot be %s", common.MigrationHistoryCollection)
		}

		_, ok := dup[coll.Collection().Spec().Name]
		if ok {
			return fmt.Errorf("Duplicate collection found with name: %s", coll.Collection().Spec().Name)
		}

		dup[coll.Collection().Spec().Name] = true
	}

	return nil
}

func validateID(collectionName string, fields []collection.Field) error {
	// allowed _id field types:
	// - default Object ID
	// - integers
	// - double/float
	// - string
	for _, _field := range fields {
		if _field.Spec().Name == "_id" {
			if !util.InListEq(_field.Spec().Type, []field.FieldType{
				field.TypeInt32,
				field.TypeInt64,
				field.TypeDouble,
				field.TypeString,
			}) {
				return fmt.Errorf("%s: ID field type invalid. Allowed types: integers, double, and string", collectionName)
			}
		}
	}

	return nil
}

func validateFieldDuplication(collectionName string, fields []collection.Field) error {
	var findDuplicate func(parent string, fs []collection.Field) error
	findDuplicate = func(parent string, fs []collection.Field) error {
		dup := map[string]bool{}
		for _, f := range fs {
			_, ok := dup[f.Spec().Name]
			if ok {
				return fmt.Errorf("%s: duplicate field found: %s%s", collectionName, parent, f.Spec().Name)
			}

			dup[f.Spec().Name] = true

			// check deeper
			if f.Spec().ArrayFields != nil {
				err := findDuplicate(fmt.Sprintf("%s.", f.Spec().Name), collection.FieldsFromSpecs(f.Spec().ArrayFields))
				if err != nil {
					return err
				}
			}

			if f.Spec().Object != nil {
				err := findDuplicate(fmt.Sprintf("%s.", f.Spec().Name), collection.FieldsFromSpecs(f.Spec().Object))
				if err != nil {
					return err
				}
			}
		}

		return nil
	}

	return findDuplicate("", fields)
}

func validateIndividualField(collectionName, path string, _field collection.Field, insideArray bool) error {
	if len(_field.Spec().Name) > 128 {
		return fmt.Errorf("%s: Cannot have a field name more than 128 characters len on field: %s.%s", collectionName, path, _field.Spec().Name)
	}

	// any field type except a must not be empty
	if !insideArray {
		if _field.Spec().Name == "" {
			return fmt.Errorf("%s: Field name must not be empty for path: %s, type: %s", collectionName, path, _field.Spec().Type.ToString())
		}
	}

	if path != "" {
		path += "."
	}

	switch _field.Spec().Type {
	case field.TypeArray:
		if _field.Spec().ArrayFields == nil {
			return fmt.Errorf("%s: ArrayFields must not be empty for path: %s, type: %s", collectionName, path, _field.Spec().Type.ToString())
		}

		if len(*_field.Spec().ArrayFields) != 1 {
			return fmt.Errorf("%s: ArrayFields must have exactly 1 child for path: %s, type: %s", collectionName, path, _field.Spec().Type.ToString())
		}

		for _, child := range *_field.Spec().ArrayFields {
			err := validateIndividualField(collectionName, path+_field.Spec().Name, collection.FieldFromSpec(&child), true)
			if err != nil {
				return err
			}
		}
	case field.TypeObject:
		if _field.Spec().Object == nil {
			return fmt.Errorf("%s: Object must not be empty for path: %s, type: %s", collectionName, path, _field.Spec().Type.ToString())
		}

		for _, child := range *_field.Spec().Object {
			err := validateIndividualField(collectionName, path+_field.Spec().Name, collection.FieldFromSpec(&child), false)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func validateFields(collectionName string, fields []collection.Field) error {
	err := validateFieldDuplication(collectionName, fields)
	if err != nil {
		return err
	}

	for _, _field := range fields {
		if err = validateIndividualField(collectionName, "", _field, false); err != nil {
			return err
		}
	}

	return err
}

func validateIndexDuplication(collectionName string, indexes []collection.Index) error {
	dup := map[string]bool{}
	for _, idx := range indexes {
		_, ok := dup[idx.Spec().GetKey()]
		if ok {
			return fmt.Errorf("%s: duplicate index found: %s", collectionName, idx.Spec().GetKey())
		}

		dup[idx.Spec().GetKey()] = true
	}

	// TODO: differentiate usecase for each option. i.e: an option might not be necessary for comparison

	return nil
}

func validateIndexWithFields(collectionName string, fields []collection.Field, _index collection.Index) error {
	// Index Fields cannot be empty
	if len(_index.Spec().Fields) == 0 {
		return fmt.Errorf("%s: Index Fields cannot be empty: %s", collectionName, _index.Spec().GetName())
	}

	// check whether index keys are in fields
	var fieldExists func(path []string, _field collection.Field) bool
	fieldExists = func(path []string, _field collection.Field) bool {
		// assuming path is not empty
		curr := path[0]
		if curr != _field.Spec().Name {
			return false
		}

		res := false
		if len(path) > 1 {
			switch _field.Spec().Type {
			case field.TypeArray:
				if _field.Spec().ArrayFields != nil {
					for _, child := range *_field.Spec().ArrayFields {
						res = fieldExists(path[1:], collection.FieldFromSpec(&child))
						if res {
							break
						}
					}
				}
			case field.TypeObject:
				if _field.Spec().Object != nil {
					for _, child := range *_field.Spec().Object {
						res = fieldExists(path[1:], collection.FieldFromSpec(&child))
						if res {
							break
						}
					}
				}
			}
		} else {
			res = true
		}

		return res
	}

	var checkFieldType func(path []string, _field collection.Field, expectedType field.FieldType) bool
	checkFieldType = func(path []string, _field collection.Field, expectedType field.FieldType) bool {
		// assuming path is not empty
		curr := path[0]
		if curr != _field.Spec().Name {
			return false
		}

		if len(path) > 1 {
			switch _field.Spec().Type {
			case field.TypeArray:
				if _field.Spec().ArrayFields != nil {
					for _, child := range *_field.Spec().ArrayFields {
						res := checkFieldType(path[1:], collection.FieldFromSpec(&child), expectedType)
						if res {
							return true
						}
					}
				}
			case field.TypeObject:
				if _field.Spec().Object != nil {
					for _, child := range *_field.Spec().Object {
						res := checkFieldType(path[1:], collection.FieldFromSpec(&child), expectedType)
						if res {
							return true
						}
					}
				}
			}
		}

		return _field.Spec().Type == expectedType
	}

	pathExists := func(path []string) bool {
		ok := false
		for _, _field := range fields {
			ok = fieldExists(path, _field)
			if ok {
				break
			}
		}

		return ok
	}

	pathHasType := func(path []string, expectedType field.FieldType) bool {
		ok := false
		for _, _field := range fields {
			ok = checkFieldType(path, _field, expectedType)
			if ok {
				break
			}
		}

		return ok
	}

	// cross checking index keys x fields
	for _, indexField := range _index.Spec().Fields {
		path := strings.Split(indexField.Key, ".")
		if !pathExists(path) {
			return fmt.Errorf("%s: index key is invalid: %s", collectionName, indexField.Key)
		}
	}

	// validate by index type
	switch _index.Spec().Type {
		// TODO: complete if the usecase is clear
		// related to this commit, migt be moved here:
		// https://github.com/amirkode/go-mongr8/commit/45060b493e03b7631b5c81b2684f760d10305d09
		//
		// These cases should be handled:
		// - 2 text indenxes on the same collection
		//   source: https://www.mongodb.com/docs/manual/core/indexes/index-types/index-text/text-index-restrictions/?utm_source=chatgpt.com#text-index-restrictions-on-self-managed-deployments
		// - compound indexes invalid type with text index
		//   source: https://www.mongodb.com/docs/manual/core/indexes/index-types/index-text/text-index-restrictions/?utm_source=chatgpt.com#compound-text-index
		// - ascending and descending index in the same field
	}

	// validate index options
	if _index.Spec().Rules != nil {
		rules := *_index.Spec().Rules
		if _index.Spec().HasRule(index.OptionPartialFilterExp) {
			filters := rules[index.OptionPartialFilterExp].(map[string]interface{})
			for key, _ := range filters {
				path := strings.Split(key, ".")
				if !pathExists(path) {
					return fmt.Errorf("%s: Partial filter key is invalid: %s", collectionName, key)
				}
			}
		}

		if _index.Spec().HasRule(index.OptionTTL) {
			// There must be a least a timestamp field
			ok := false
			for _, indexField := range _index.Spec().Fields {
				path := strings.Split(indexField.Key, ".")
				ok = pathHasType(path, field.TypeTimestamp)
				if ok {
					break
				}
			}

			if !ok {
				return fmt.Errorf("%s: Timestamp field must exist in TTL index: %s", collectionName, _index.Spec().GetName())
			}
		}
	}

	return nil
}

func validateIndexes(collectionName string, fields []collection.Field, indexes []collection.Index) error {
	err := validateIndexDuplication(collectionName, indexes)
	if err != nil {
		return err
	}

	for _, _index := range indexes {
		if err = validateIndexWithFields(collectionName, fields, _index); err != nil {
			return err
		}
	}

	return nil
}

// TODO: validation for "(InvalidOptions) 'expireAfterSeconds' is only supported on time-series collections"
