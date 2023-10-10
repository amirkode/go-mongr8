package api_interpreter

import (
	"fmt"

	"github.com/amirkode/go-mongr8/collection/field"
	"github.com/amirkode/go-mongr8/collection/index"
	"github.com/amirkode/go-mongr8/collection/metadata"
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
		Collection metadata.MetadataSpec
		Fields     *[]field.FieldSpec
		Indexes    *[]index.IndexSpec
	}

	// SubAction handles atomic action as part of main Action
	SubAction struct {
		// GetStatement returns code statement in the official MongoDB API
		GetStatement func() string
		// this returns string literal definition of ActionSchema
		GetSchemaLiteral func() string
		// could be the name of collection, index, field, or etc
		Type SubActionType
		// schema changed in this action
		ActionSchema SubActionSchema
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

	ActionIf interface {
		GetStatements() string
		// this returns the schema whether added or removed by current action
	}

	Action struct {
		ActionIf
		Type ActionType
		// sub actions should be execute respectively
		SubActions []SubAction
	}
)

func (a Action) GetRawStatements() string {
	res := ""
	for index, subAction := range a.SubActions {
		res += fmt.Sprintf("Step %d:\n%s\n", index+1, subAction.GetStatement())
	}

	return res
}
