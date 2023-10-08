package action

import (
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
)

// TODO: NEED TO FIND BETTER APPROACH
// after some considerations, there no need to define another layer above MongoDB API
// to make the migration files more concise
// this supposed to generate raw string of MongoDB API statements

type (
	// SubAction handles atomic action as part of main Action
	SubAction struct {
		Execute func()
		GetRawStatement func() string
		// could be the name of collection, index, field, or etc
		Name string
		// payload could any object format
		// passed as argument in the official MongoDB API
		// could be more than one arguments
		Payload []bson.M
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
		Execute()
		GetRawStatements() string
	}

	Action struct {
		ActionIf
		Type ActionType
		// sub actions should be execute respectively
		SubActions []SubAction
	}
)

func (a Action) Execute() {
	// execute all sub actions
	for _, subAction := range a.SubActions {
		subAction.Execute()
	}
}

func (a Action) GetRawStatements() string {
	res := ""
	for index, subAction := range a.SubActions {
		res += fmt.Sprintf("Step %d:\n%s\n", index + 1, subAction.GetRawStatement)
	}
	
	return res
}
