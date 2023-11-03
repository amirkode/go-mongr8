package sync_strategy

import (
	"fmt"
	"testing"

	"internal/util"
	"internal/test"

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
					t.Errorf(msg)
					panic(msg)		
				}
			}
		} else {
			msg := fmt.Sprintf("Case 1: Invalid action %s", action.ActionKey)
			t.Errorf(msg)
			panic(msg)
		}
	}

	// TODO: Add more cases
	// - FIELD CONVERSION
}

func TestGetCollectionFromMigrations(t *testing.T) {
	// Case 1: migrations of different actions: create collection, field, and index
	case1Migrations := []migrator.Migration{
		{
			ID: "1",
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
			ID: "2",
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
	// the number of collections must be 2 (users and customers)
	test.AssertEqual(t, len(case1Collections), 2, "Case 1: The collections length must be 2 (users and customers)")
	// check collection structure
	// TODO: replace this test with a single global function to check structure with DFS
	for _, collection := range case1Collections {
		if !util.InListEq(collection.Collection().Spec().Name, []string{"users", "customers"}) {
			msg := "Case 1: Collection name must be users or customers"
			t.Errorf(msg)
			panic(msg)	
		}

		switch collection.Collection().Spec().Name {
		case "users":
			// check fields
			test.AssertEqual(t, len(collection.Fields()), 3, "Case 1: fields length for collection users must be 3")
			for _, f := range collection.Fields() {
				switch f.Spec().Name {
				case "name":
					test.AssertEqual(t, f.Spec().Type, field.TypeString,
						fmt.Sprintf("Case 1: field type for %s on collection users must be a string", 
							f.Spec().Name,
						),
					)
				case "age":
					test.AssertEqual(t, f.Spec().Type, field.TypeInt32,
						fmt.Sprintf("Case 1: field type for %s on collection users must be an integer", 
							f.Spec().Name,
						),
					)
				case "additional_info":
					test.AssertEqual(t, f.Spec().Type, field.TypeObject,
						fmt.Sprintf("Case 1: field type for %s on collection users must be an integer", 
							f.Spec().Name,
						),
					)
					test.AssertTrue(t, f.Spec().Object != nil && len(*f.Spec().Object) == 2,
					fmt.Sprintf("Case 1: object length for %s on collection users must be 2", 
							f.Spec().Name,
						),
					)
					// assuming after this pass, the children of additional_info are correct
				default:
					msg := fmt.Sprintf("Case 1: Unexpected field name %s for collection users", 
						f.Spec().Name,
					)
					t.Errorf(msg)
					panic(msg)
				}
			}
			// check indexes
			test.AssertEqual(t, len(collection.Indexes()), 1, "Case 1: indexes length for collection users must be 1")
			for _, idx := range collection.Indexes() {
				test.AssertEqual(t, idx.Spec().Type, index.TypeCompound, "Case 1: index type for collection users must be Compound Index")
				test.AssertEqual(t, len(idx.Spec().Fields), 2, "Case 1: index fields length for collection users must be 2")
				for _, idxField := range idx.Spec().Fields {
					switch idxField.Key {
					case "name":
						test.AssertEqual(t, idxField.Value, -1, fmt.Sprintf("Case 1: Compound index field %s value for collection users must be -1", 
							idxField.Key,
						))
					case "age":
						test.AssertEqual(t, idxField.Value, 1, fmt.Sprintf("Case 1: Compound index field %s value for collection users must be -1", 
							idxField.Key,
						))
					default:
						msg := "Case 1: Compound index field for collection users must be either name or age"
						t.Errorf(msg)
						panic(msg)
					}
				}
			}
		case "customers":
			// check fields
			test.AssertEqual(t, len(collection.Fields()), 3, "Case 1: fields length for collection customers must be 3")
			for _, f := range collection.Fields() {
				switch f.Spec().Name {
				case "name":
					test.AssertEqual(t, f.Spec().Type, field.TypeString,
						fmt.Sprintf("Case 1: field type for %s on collection customers must be a string", 
							f.Spec().Name,
						),
					)
				case "bio":
					test.AssertEqual(t, f.Spec().Type, field.TypeString,
						fmt.Sprintf("Case 1: field type for %s on collection customers must be a string", 
							f.Spec().Name,
						),
					)
				case "address":
					test.AssertEqual(t, f.Spec().Type, field.TypeString,
						fmt.Sprintf("Case 1: field type for %s on collection customers must be a string", 
							f.Spec().Name,
						),
					)
				default:
					msg := fmt.Sprintf("Case 1: Unexpected field name %s for collection customers", 
						f.Spec().Name,
					)
					t.Errorf(msg)
					panic(msg)
				}
			}
			// check indexes
			test.AssertEqual(t, len(collection.Indexes()), 2, "Case 1: indexes length for collection customers must be 2")
			for _, idx := range collection.Indexes() {
				test.AssertEqual(t, idx.Spec().Type, index.TypeSingleField, "Case 1: index type for collection customers must be Single Field Index")
				test.AssertEqual(t, len(idx.Spec().Fields), 1, "Case 1: index fields length for collection customers must be 1")
				for _, idxField := range idx.Spec().Fields {
					switch idxField.Key {
					case "bio":
						test.AssertEqual(t, idxField.Value, 1, fmt.Sprintf("Case 1: Single Field index field %s value for collection customers must be 1", 
							idxField.Key,
						))
					case "address":
						test.AssertEqual(t, idxField.Value, -1, fmt.Sprintf("Case 1: Single Field index field %s value for collection customers must be -1", 
							idxField.Key,
						))
					default:
						msg := "Case 1: Single Field index field for collection users must be either bio or address"
						t.Errorf(msg)
						panic(msg)
					}
				}
			}
		}
	}

	// TODO: add more cases
}