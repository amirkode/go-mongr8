/*
Copyright (c) 2023-present the go-mongr8 Authors and Contributors
[@see Authors file]

Licensed under the MIT License
(https://opensource.org/licenses/MIT)
*/
package schema_interpreter

import (
	"github.com/amirkode/go-mongr8/collection"
	"github.com/amirkode/go-mongr8/collection/field"
	dt "github.com/amirkode/go-mongr8/internal/data_type"

	"go.mongodb.org/mongo-driver/bson"
)

// TODO: NEED TO FIND BETTER APPROACH
// after some considerations, there no need to define another layer above MongoDB API
// to make the migration files more concise
// this supposed to generate raw string of MongoDB API statements

type (
	// SubActionSchema holds what are being changed/adjusted in the current SubAction
	// this data will be used to compare the generated migration files schema
	// with the higher level user-defined schema in /mongr8/collection
	SubActionSchema struct {
		// Collection must always be defined
		// if Fields and Indexes are not defined, then it must be
		// related to collection option, i.e: Create Collection
		Collection collection.Metadata
		Fields     []collection.Field
		Indexes    []collection.Index
		// field type for conversion
		// we're expecting only a single field conversion
		// each sub action
		FieldConvertFrom *field.FieldType
	}

	SubActionIf interface {
		// flag whether current sub action is Up or Down
		IsUp() bool
		// this will return pairs of index and rule
		GetIndexesBsonD() []dt.Pair[string, dt.Pair[bson.D, bson.D]]
		GetIndexesBsonM() []dt.Pair[string, dt.Pair[bson.M, bson.M]]
		// this will return whole field schema in the collection
		GetFieldsBsonD() bson.D
		GetFieldsBsonM() bson.M
	}

	// SubAction handles atomic action as part of main Action
	SubAction struct {
		SubActionIf
		// could be the name of collection, index, field, or etc
		Type SubActionType
		// schema changed in this action
		ActionSchema SubActionSchema
		// validate function
		validate func()
	}

	// Action is an entity for storing actionable item to run on MongoDB
	// here are supported operations in this package:
	// - Create new collection
	// - Create indexes
	// - Create fields to existing collection
	// - Drop collection
	// - Drop index
	// - Drop field
	// future support:
	// - Add advanced collection rules
	// - Change field type (for now it's a bit challenging)
	//
	// since this is a usecase level entity, there might be more than one sub actions
	// e.g: in Create new collection, we need to create collection definition,
	// then insert some insert values to main the stucture integrity
	//
	// also since, this is the higher level represention of MongoDB statement
	// there must be implementation of each action

	// ActionIf interface {
	// 	GetStatements() string
	// 	// this returns the schema whether added or removed by current action
	// }

	// Action holds operations that will be performed in a collection
	Action struct {
		// ActionIf
		ActionKey string
		// sub actions should be execute respectively
		SubActions []SubAction
	}
)

// func (a Action) GetRawStatements() string {
// 	res := ""
// 	for index, subAction := range a.SubActions {
// 		res += fmt.Sprintf("Step %d:\n%s\n", index+1, subAction.GetStatement())
// 	}

// 	return res
// }
