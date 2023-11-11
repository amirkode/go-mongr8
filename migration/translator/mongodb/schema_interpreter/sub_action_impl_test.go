/*
Copyright (c) 2023 the go-mongr8 Authors and Contributors
[@see Authors file]

Licensed under the MIT License
(https://opensource.org/licenses/MIT)
*/
package schema_interpreter

import (
	// "fmt"
	"testing"

	// "internal/test"

	"github.com/amirkode/go-mongr8/collection"
	"github.com/amirkode/go-mongr8/collection/field"
	"github.com/amirkode/go-mongr8/collection/index"
	"github.com/amirkode/go-mongr8/collection/metadata"
)

func TestGetLiteralInstance(t *testing.T) {
	// case 1: general case
	case1SubAction := SubAction{
		Type: SubActionTypeCreateCollection,
		ActionSchema: SubActionSchema{
			Collection: metadata.InitMetadata("users"),
			Fields: []collection.Field{
				field.StringField("name"),
				field.Int32Field("age"),
				field.ObjectField("additional_info",
					field.StringField("address"),
					field.ArrayField("score_history",
						field.Int32Field("score"),
					),
				),
			},
			Indexes: []collection.Index{
				index.CompoundIndex(
					index.Field("name", -1),
					index.Field("age", 1),
				),
			},
		},
	}

	// if no panic, then it should work properly
	case1SubAction.GetLiteralInstance("", false)
	
	// TODO: add value checking
}

// TODO: write some other tests
