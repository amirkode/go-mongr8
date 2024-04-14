/*
Copyright (c) 2023-present the go-mongr8 Authors and Contributors
[@see Authors file]

Licensed under the MIT License
(https://opensource.org/licenses/MIT)
*/
package sync_strategy

import (
	"fmt"

	"github.com/amirkode/go-mongr8/collection"
	"github.com/amirkode/go-mongr8/collection/field"
	"github.com/amirkode/go-mongr8/collection/metadata"
)

func fieldsAreEqual(a, b collection.Field) bool {
	// check both types
	if a.Spec().Type != b.Spec().Type {
		fmt.Printf("%s, %s -> Field type are different %s and %s\n", a.Spec().Name, b.Spec().Name, a.Spec().Type.ToString(), b.Spec().Type.ToString())
		return false
	}
	// check array comparation
	if a.Spec().Type == field.TypeArray {
		if a.Spec().ArrayFields == nil || b.Spec().ArrayFields == nil {
			fmt.Println("Array Fields cannot be nil")
			return false
		}

		if len(*a.Spec().ArrayFields) != len(*b.Spec().ArrayFields) {
			fmt.Println("ArrayFields length are different")
			return false
		}

		aArr := map[string]collection.Field{}
		bArr := map[string]collection.Field{}
		for _, aA := range *a.Spec().ArrayFields {
			curr := aA // assign to new address
			aArr[aA.Name] = field.FromFieldSpec(&curr)
		}
		for _, bA := range *b.Spec().ArrayFields {
			curr := bA // assign to new addres
			bArr[bA.Name] = field.FromFieldSpec(&curr)
		}

		// check aArr on bArr
		for _, aA := range *a.Spec().ArrayFields {
			bA, ok := bArr[aA.Name]
			if !ok {
				fmt.Println("Field aA not in bArr:", aA.Name)
				return false
			}

			ok = fieldsAreEqual(field.FromFieldSpec(&aA), bA)
			if !ok {
				fmt.Println("Array Field Depth check field")
				return false
			}
		}

		// check barr on aArr
		for _, bA := range *b.Spec().ArrayFields {
			_, ok := aArr[bA.Name]
			if !ok {
				fmt.Println("Field bA not in aArr:", bA.Name)
				return false
			}
		}
	}

	// check object comparation
	if a.Spec().Type == field.TypeObject {
		if a.Spec().Object == nil || b.Spec().Object == nil {
			fmt.Println("Object Fields cannot be empty")
			return false
		}

		if len(*a.Spec().Object) != len(*b.Spec().Object) {
			fmt.Println("Object Fields length are different")
			return false
		}

		aObj := map[string]collection.Field{}
		bObj := map[string]collection.Field{}
		for _, aA := range *a.Spec().Object {
			fmt.Println("aA:", aA.Name)
			fmt.Println("aA address:", &aA)
			curr := aA /// assign to new address
			aObj[aA.Name] = field.FromFieldSpec(&curr)
		}
		for _, bA := range *b.Spec().Object {
			fmt.Println("bA:", bA.Name)
			fmt.Println("bA address:", &bA)
			curr := bA // assign to new address
			bObj[bA.Name] = field.FromFieldSpec(&curr)
		}

		// check aArr on bArr
		for _, aA := range *a.Spec().Object {
			bA, ok := bObj[aA.Name]
			if !ok {
				fmt.Println("Field aA not in bObj:", aA.Name)
				return false
			}

			ok = fieldsAreEqual(field.FromFieldSpec(&aA), bA)
			if !ok {
				fmt.Println("Object Field Depth check field: ", aA.Name, bA.Spec().Name)
				return false
			}
		}

		// check barr on aArr
		for _, bA := range *b.Spec().Object {
			_, ok := aObj[bA.Name]
			if !ok {
				fmt.Println("Field bA not in aObj:", bA.Name)
				return false
			}
		}
	}

	return true
}

func collectionsAreEqual(a, b collection.Collection) bool {
	// check name
	if a.Collection().Spec().Name != b.Collection().Spec().Name {
		fmt.Println("collection name are different")
		return false
	}

	// check type
	if a.Collection().Spec().Type != b.Collection().Spec().Type {
		fmt.Println("collection type are different")
		return false
	}

	// check options
	if (a.Collection().Spec().Options == nil && b.Collection().Spec().Options != nil) ||
		(a.Collection().Spec().Options != nil && b.Collection().Spec().Options == nil) {
		fmt.Println("Collection option are different 1")
		return false
	}

	if a.Collection().Spec().Options != nil {
		options := metadata.GetAllOptionKeys()
		aOpts := *a.Collection().Spec().Options
		bOpts := *b.Collection().Spec().Options
		for _, opt := range options {
			aOpt, okA := aOpts[opt]
			bOpt, okB := bOpts[opt]
			if okA != okB {
				fmt.Println("Collection option are different 2")
				return false
			}

			if okA && aOpt != bOpt {
				fmt.Println("Collection option are different 3")
				return false
			}
		}
	}

	// check fields
	if len(a.Fields()) != len(b.Fields()) {
		fmt.Println("Fields length are different", len(a.Fields()), "and", len(b.Fields()))
		return false
	}

	aFields := map[string]collection.Field{}
	bFields := map[string]collection.Field{}
	for _, aField := range a.Fields() {
		aFields[aField.Spec().Name] = aField
	}
	for _, bField := range b.Fields() {
		bFields[bField.Spec().Name] = bField
	}

	// check a fields on b fields
	for _, aField := range a.Fields() {
		bField, ok := bFields[aField.Spec().Name]
		if !ok {
			fmt.Printf("Field not found 1: %s\n", aField.Spec().Name)
			return false
		}

		// check deeper level
		ok = fieldsAreEqual(aField, bField)
		if !ok {
			return false
		}
	}

	// check b fields on a fields
	for _, bField := range b.Fields() {
		_, ok := aFields[bField.Spec().Name]
		if !ok {
			fmt.Println("Field not found 2:", bField.Spec().Name)
			return false
		}
	}

	// check indexes
	if len(a.Indexes()) != len(b.Indexes()) {
		fmt.Println("Collection indexes length are different")
		return false
	}

	aIndexes := map[string]bool{}
	bIndexes := map[string]bool{}
	for _, aIndex := range a.Indexes() {
		aIndexes[aIndex.Spec().GetKey()] = true
	}
	for _, bIndex := range b.Indexes() {
		bIndexes[bIndex.Spec().GetKey()] = true
	}

	// check a indexes on b indexes
	for _, aIndex := range a.Indexes() {
		_, ok := bIndexes[aIndex.Spec().GetKey()]
		if !ok {
			fmt.Println("index key are different 1:", aIndex.Spec().GetKey())
			return false
		}
	}

	// check b indexes on a indexes
	for _, bIndex := range b.Indexes() {
		_, ok := aIndexes[bIndex.Spec().GetKey()]
		if !ok {
			fmt.Println("index key are different 2:", bIndex.Spec().GetKey())
			return false
		}
	}

	return true
}
