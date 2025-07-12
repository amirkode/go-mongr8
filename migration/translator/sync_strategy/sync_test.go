/*
Copyright (c) 2023-present the go-mongr8 Authors and Contributors
[@see Authors file]

Licensed under the MIT License
(https://opensource.org/licenses/MIT)
*/
package sync_strategy

import (
	"fmt"
	"testing"

	"github.com/amirkode/go-mongr8/internal/test"
	"github.com/amirkode/go-mongr8/internal/util"

	"github.com/amirkode/go-mongr8/collection"
	"github.com/amirkode/go-mongr8/collection/field"
	"github.com/amirkode/go-mongr8/collection/index"
	"github.com/amirkode/go-mongr8/collection/metadata"
	"github.com/amirkode/go-mongr8/migration/migrator"
	si "github.com/amirkode/go-mongr8/migration/translator/mongodb/schema_interpreter"
)

func TestGetSignedCollection(t *testing.T) {
	// case 1: default case
	case1Collection := collection.NewCollection(
		metadata.InitMetadata("users"),
		[]collection.Field{
			field.StringField("name"),
			field.Int32Field("age"),
			field.ObjectField("additional_info",
				field.StringField("address"),
				field.ArrayField("score_history",
					field.Int32Field("score"),
				),
			),
		},
		[]collection.Index{
			index.CompoundIndex(
				index.Field("name", -1),
				index.Field("age", 1),
			),
		},
	)
	case1SignedCollection := getSignedCollection(case1Collection)
	// check items
	test.AssertEqual(t, case1Collection.Collection().Spec().Name, case1SignedCollection.Metadata.Spec().Name,
		"Case 1: Collection name must be equal",
	)
	test.AssertEqual(t, len(case1Collection.Fields()), len(case1SignedCollection.Fields), "Case 1: Fields length must be equal")
	test.AssertEqual(t, len(case1Collection.Indexes()), len(case1SignedCollection.Indexes), "Case 1: Indexes length must be equal")
	// check signs
	test.AssertEqual(t, case1SignedCollection.Sign, SignPlus, "Case 1: Default entity sign for collection must be SignPlus")
	for _, f := range case1SignedCollection.Fields {
		test.AssertEqual(t, f.Sign, SignPlus, "Case 1: Default enetity sign for field must be SignPlus")
	}
	for _, idx := range case1SignedCollection.Indexes {
		test.AssertEqual(t, idx.Sign, SignPlus, "Case 1: Default enetity sign for index must be SignPlus")
	}
}

func TestSyncCollections(t *testing.T) {
	// case 1: new collection
	case1Incoming := []collection.Collection{
		collection.NewCollection(
			metadata.InitMetadata("users"),
			[]collection.Field{
				field.StringField("name"),
				field.Int32Field("age"),
				field.ObjectField("additional_info",
					field.StringField("address"),
					field.ArrayField("score_history",
						field.Int32Field("score"),
					),
				),
			},
			[]collection.Index{
				index.CompoundIndex(
					index.Field("name", -1),
					index.Field("age", 1),
				),
			},
		),
	}
	case1Origin := []collection.Collection{}
	case1Synced := SyncCollections(case1Incoming, case1Origin)

	test.AssertEqual(t, len(case1Incoming), len(case1Synced),
		"Case 1: Synced collection length must be equal to incoming on new collection",
	)
	// check all the signs are SignPlus
	for _, s := range case1Synced {
		test.AssertEqual(t, s.Sign, SignPlus, "Case 1: All synced collection must SingPlus on new collection")
	}

	// Case 1: delete all collection
	case2Incoming := []collection.Collection{}
	case2Origin := []collection.Collection{
		collection.NewCollection(
			metadata.InitMetadata("users"),
			[]collection.Field{
				field.StringField("name"),
				field.Int32Field("age"),
				field.ObjectField("additional_info",
					field.StringField("address"),
					field.ArrayField("score_history",
						field.Int32Field("score"),
					),
				),
			},
			[]collection.Index{
				index.CompoundIndex(
					index.Field("name", -1),
					index.Field("age", 1),
				),
			},
		),
	}
	case2Synced := SyncCollections(case2Incoming, case2Origin)

	test.AssertEqual(t, len(case2Origin), len(case2Synced),
		"Case 1: Synced collection length must be equal to origin on new collection",
	)
	// check all the signs are SignMinus
	for _, s := range case2Synced {
		test.AssertEqual(t, s.Sign, SignMinus, "Case 1: All synced collection must SingMinus on new collection")
	}

	// TODO: complete other cases
}

func TestGetActions(t *testing.T) {
	// Case 1: Create collection, create field, and create indexe
	case1Incoming := []collection.Collection{
		collection.NewCollection(
			metadata.InitMetadata("users"),
			[]collection.Field{
				field.StringField("name"),
				field.Int32Field("age"),
				field.ObjectField("additional_info",
					field.StringField("address"),
					field.ArrayField("score_history",
						field.Int32Field("score"),
					),
				),
			},
			[]collection.Index{
				index.CompoundIndex(
					index.Field("name", -1),
					index.Field("age", 1),
				),
			},
		),
		collection.NewCollection(
			metadata.InitMetadata("customers"),
			[]collection.Field{
				field.StringField("name"),
				field.StringField("bio"),
			},
			[]collection.Index{
				index.SingleFieldIndex(index.Field("bio", 1)),
			},
		),
	}
	case1Origin := []collection.Collection{
		collection.NewCollection(
			metadata.InitMetadata("customers"),
			[]collection.Field{
				field.StringField("name"),
			},
			[]collection.Index{},
		),
	}
	case1Actions := GetActions(case1Incoming, case1Origin)

	// check the actions length, it must be 3 SubActions:
	// - create new collection `Users`
	// - create new field `bio` on collection `customers`
	// - create new single field index `{"bio": 1}` on collection `customers`
	// with 2 grouped actions by collection name
	test.AssertEqual(t, len(case1Actions.First), 2, "Case 1: Up Actions length must be 2")
	test.AssertEqual(t, len(case1Actions.Second), 2, "Case 1: Down Actions length must be 2")
	// check Up Actions
	for _, action := range case1Actions.First {
		if action.ActionKey == "users" {
			test.AssertEqual(t, len(action.SubActions), 1, "Case 1: Subactions length for users collection must be 1")
			test.AssertEqual(t, action.SubActions[0].Type, si.SubActionTypeCreateCollection,
				"Case 1: Sub Action Type for users collection must SubActionTypeCreateCollection",
			)
		} else if action.ActionKey == "customers" {
			test.AssertEqual(t, len(action.SubActions), 2, "Case 1: Subactions length for customers collection must be 2")
			for _, subAction := range action.SubActions {
				if !util.InListEq(subAction.Type, []si.SubActionType{
					si.SubActionTypeCreateField,
					si.SubActionTypeCreateIndex,
				}) {
					msg := "Case 1: Sub Action Type for customers collection must SubActionTypeCreateField or SubActionTypeCreateIndex"
					t.Errorf("%s", msg)
					panic(msg)
				}
			}
		} else {
			msg := fmt.Sprintf("Case 1: Invalid action %s", action.ActionKey)
			t.Errorf("%s", msg)
			panic(msg)
		}
	}

	// Case 2: Field Conversion
	case2Incoming := []collection.Collection{
		collection.NewCollection(
			metadata.InitMetadata("customers"),
			[]collection.Field{
				field.ObjectField("other_info",
					field.ArrayField("tx_history",
						field.ObjectField("",
							field.StringField("amount"),
						),
					),
				),
			},
			[]collection.Index{},
		),
	}
	case2Origin := []collection.Collection{
		collection.NewCollection(
			metadata.InitMetadata("customers"),
			[]collection.Field{
				field.ObjectField("other_info",
					field.ArrayField("tx_history",
						field.ObjectField("",
							field.Int32Field("amount"),
						),
					),
				),
			},
			[]collection.Index{},
		),
	}
	case2Actions := GetActions(case2Incoming, case2Origin)

	test.AssertEqual(t, len(case2Actions.First), 1, "Case 2: Up Actions length must be 1")
	test.AssertEqual(t, len(case2Actions.Second), 1, "Case 2: Down Actions length must be 1")
	// check Up Actions
	test.AssertEqual(t, len(case2Actions.First[0].SubActions), 1, "Case 2: Unexpected Sub Actions length")
	test.AssertEqual(t, *case2Actions.First[0].SubActions[0].ActionSchema.FieldConvertFrom, field.TypeInt32, "Case 2: Unexpected Convert From Type")
	case2ExpectedPayloadAsCollection := case2Incoming[0]
	case2ActualPayloadAsCollection := collection.NewCollection(
		case2Actions.First[0].SubActions[0].ActionSchema.Collection,
		case2Actions.First[0].SubActions[0].ActionSchema.Fields,
		[]collection.Index{},
	)

	test.AssertTrue(t, collectionsAreEqual(case2ExpectedPayloadAsCollection, case2ActualPayloadAsCollection), "Case 2: Unexpected Action Payload")

	// TODO: Add more cases
}

func TestGetCollectionFromMigrations(t *testing.T) {
	// Case 1: migrations of different actions: create collection, field, and index
	case1Migrations := []migrator.Migration{
		{
			ID:   "1",
			Desc: "a description",
			Up: []si.Action{
				{
					ActionKey: "users",
					SubActions: []si.SubAction{
						{
							Type: si.SubActionTypeCreateCollection,
							ActionSchema: si.SubActionSchema{
								Collection: metadata.InitMetadata("users"),
								Fields: []collection.Field{
									field.StringField("name"),
									field.Int32Field("age"),
								},
							},
						},
						{
							Type: si.SubActionTypeCreateIndex,
							ActionSchema: si.SubActionSchema{
								Collection: metadata.InitMetadata("users"),
								Indexes: []collection.Index{
									index.CompoundIndex(
										index.Field("name", -1),
										index.Field("age", 1),
									),
								},
							},
						},
						{
							Type: si.SubActionTypeCreateField,
							ActionSchema: si.SubActionSchema{
								Collection: metadata.InitMetadata("users"),
								Fields: []collection.Field{
									field.ObjectField("additional_info",
										field.StringField("address"),
										field.ArrayField("score_history",
											field.Int32Field("score"),
										),
									),
								},
							},
						},
					},
				},
				{
					ActionKey: "customers",
					SubActions: []si.SubAction{
						{
							Type: si.SubActionTypeCreateCollection,
							ActionSchema: si.SubActionSchema{
								Collection: metadata.InitMetadata("customers"),
								Fields: []collection.Field{
									field.StringField("name"),
								},
							},
						},
						{
							Type: si.SubActionTypeCreateField,
							ActionSchema: si.SubActionSchema{
								Collection: metadata.InitMetadata("customers"),
								Fields: []collection.Field{
									field.StringField("bio"),
								},
							},
						},
						{
							Type: si.SubActionTypeCreateIndex,
							ActionSchema: si.SubActionSchema{
								Collection: metadata.InitMetadata("customers"),
								Indexes: []collection.Index{
									index.SingleFieldIndex(index.Field("bio", 1)),
								},
							},
						},
					},
				},
			},
		},
		{
			ID:   "2",
			Desc: "a description",
			Up: []si.Action{
				{
					ActionKey: "customers",
					SubActions: []si.SubAction{
						{
							Type: si.SubActionTypeCreateField,
							ActionSchema: si.SubActionSchema{
								Collection: metadata.InitMetadata("customers"),
								Fields: []collection.Field{
									field.StringField("address"),
								},
								Indexes: []collection.Index{
									index.SingleFieldIndex(index.Field("address", -1)),
								},
							},
						},
					},
				},
			},
		},
	}
	case1Collections := GetCollectionFromMigrations(case1Migrations)
	case1ExpectedCollections := map[string]collection.Collection{
		"users": collection.NewCollection(
			metadata.InitMetadata("users"),
			[]collection.Field{
				field.StringField("name"),
				field.Int32Field("age"),
				field.ObjectField("additional_info",
					field.StringField("address"),
					field.ArrayField("score_history",
						field.Int32Field("score"),
					),
				),
			},
			[]collection.Index{
				index.CompoundIndex(
					index.Field("name", -1),
					index.Field("age", 1),
				),
			},
		),
		"customers": collection.NewCollection(
			metadata.InitMetadata("customers"),
			[]collection.Field{
				field.StringField("name"),
				field.StringField("bio"),
				field.StringField("address"),
			},
			[]collection.Index{
				index.SingleFieldIndex(index.Field("bio", 1)),
				index.SingleFieldIndex(index.Field("address", -1)),
			},
		),
	}
	// the number of collections must be 2 (users and customers)
	test.AssertEqual(t, len(case1Collections), 2, "Case 1: The collections length must be 2 (users and customers)")
	// check collection structure
	for _, collection := range case1Collections {
		compCollection, ok := case1ExpectedCollections[collection.Collection().Spec().Name]
		if !ok {
			msg := "Case 1: Collection name must be users or customers"
			t.Errorf("%s", msg)
			panic(msg)
		}

		test.AssertTrue(t, collectionsAreEqual(collection, compCollection), fmt.Sprintf("Case 1: Unexpected Collection %s", collection.Collection().Spec().Name))
	}

	// Case 2: migrations of different actions with drop field action
	case2Migrations := []migrator.Migration{
		{
			ID:   "1",
			Desc: "a description",
			Up: []si.Action{
				{
					ActionKey: "users",
					SubActions: []si.SubAction{
						*si.SubActionCreateCollection(si.SubActionSchema{
							Collection: metadata.InitMetadata("users"),
							Fields: []collection.Field{
								field.StringField("name"),
							},
							Indexes: []collection.Index{
								index.SingleFieldIndex(index.Field("name", int(1))),
							},
						}),
					},
				},
			},
		},
		{
			ID:   "2",
			Desc: "a description",
			Up: []si.Action{
				{
					ActionKey: "users",
					SubActions: []si.SubAction{
						*si.SubActionCreateField(si.SubActionSchema{
							Collection: metadata.InitMetadata("users"),
							Fields: []collection.Field{
								field.Int32Field("age"),
							},
							Indexes: []collection.Index{},
						}),
						*si.SubActionCreateIndex(si.SubActionSchema{
							Collection: metadata.InitMetadata("users"),
							Fields:     []collection.Field{},
							Indexes: []collection.Index{
								index.CompoundIndex(
									index.Field("name", int(-1)),
									index.Field("age", int(1)),
								),
							},
						}),
					},
				},
			},
		},
		{
			ID:   "3",
			Desc: "a description",
			Up: []si.Action{
				{
					ActionKey: "customers",
					SubActions: []si.SubAction{
						*si.SubActionCreateCollection(si.SubActionSchema{
							Collection: metadata.InitMetadata("customers"),
							Fields: []collection.Field{
								field.StringField("name"),
								field.StringField("addres"),
								field.ObjectField("other_info",
									field.ArrayField("tx_history",
										field.ObjectField("",
											field.StringField("id"),
											field.StringField("desc"),
											field.StringField("amount"),
										),
									),
								),
							},
							Indexes: []collection.Index{},
						}),
					},
				},
			},
		},
		{
			ID:   "4",
			Desc: "a description",
			Up: []si.Action{
				{
					ActionKey: "customers",
					SubActions: []si.SubAction{
						*si.SubActionDropField(si.SubActionSchema{
							Collection: metadata.InitMetadata("customers"),
							Fields: []collection.Field{
								field.ObjectField("other_info",
									field.ArrayField("tx_history",
										field.ObjectField("",
											field.StringField("amount"),
										),
									),
								),
							},
							Indexes: []collection.Index{},
						}),
					},
				},
			},
		},
	}
	case2Collections := GetCollectionFromMigrations(case2Migrations)
	case2ExpectedCollections := map[string]collection.Collection{
		"users": collection.NewCollection(
			metadata.InitMetadata("users"),
			[]collection.Field{
				field.StringField("name"),
				field.Int32Field("age"),
			},
			[]collection.Index{
				index.SingleFieldIndex(index.Field("name", 1)),
				index.CompoundIndex(
					index.Field("name", -1),
					index.Field("age", 1),
				),
			},
		),
		"customers": collection.NewCollection(
			metadata.InitMetadata("customers"),
			[]collection.Field{
				field.StringField("name"),
				field.StringField("addres"),
				field.ObjectField("other_info",
					field.ArrayField("tx_history",
						field.ObjectField("",
							field.StringField("id"),
							field.StringField("desc"),
						),
					),
				),
			},
			[]collection.Index{},
		),
	}
	// the number of collections must be 2 (users and customers)
	test.AssertEqual(t, len(case2Collections), len(case2ExpectedCollections), "Case 2: Unexpected collections length")
	// check collection structure
	for _, collection := range case2Collections {
		compCollection, ok := case2ExpectedCollections[collection.Collection().Spec().Name]
		if !ok {
			msg := "Case 2: Collection name must be users or customers"
			t.Errorf("%s", msg)
			panic(msg)
		}

		test.AssertTrue(t, collectionsAreEqual(collection, compCollection), fmt.Sprintf("Case 2: Unexpected Collection %s", collection.Collection().Spec().Name))
	}

	// TODO: add more cases
}
