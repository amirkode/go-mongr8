package dictionary

import (
	"fmt"
	"internal/util"

	"github.com/amirkode/go-mongr8/collection"
	"github.com/amirkode/go-mongr8/collection/field"
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
	v.validationFuncs = []func() error {
		func() error {
			return validateCollections(v.Collections)
		},
	}

	for _, coll := range v.Collections {
		v.validationFuncs = append(v.validationFuncs, func() error {
			return validateID(coll.Collection().Spec().Name, coll.Fields())
		})
		v.validationFuncs = append(v.validationFuncs, func() error {
			return validateFieldDuplication(coll.Collection().Spec().Name, coll.Fields())
		})
		v.validationFuncs = append(v.validationFuncs, func() error {
			return validateIndexDuplication(coll.Collection().Spec().Name, coll.Indexes())
		})
	}
}

func validateCollections(collections []collection.Collection) error {
	// validate duplicate name
	dup := map[string]bool{}
	for _, coll := range collections {
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

func validateIndexDuplication(collectionName string, indexes []collection.Index) error {
	for _, idx := range indexes {
		dup := map[string]bool{}
		_, ok := dup[idx.Spec().GetKey()]
		if ok {
			return fmt.Errorf("%s: duplicate index found: %s", collectionName, idx.Spec().GetKey())
		}

		dup[idx.Spec().GetKey()] = true
	}

	return nil
}

// TODO: add individual field validation, i.e: array field must not be empty