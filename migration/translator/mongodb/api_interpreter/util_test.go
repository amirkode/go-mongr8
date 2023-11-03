package api_interpreter

import (
	"reflect"
	"testing"

	"internal/test"

	"github.com/amirkode/go-mongr8/collection/field"

	"go.mongodb.org/mongo-driver/bson"
)

func bsonMAreEqual(obj1, obj2 bson.M) bool {
	if len(obj1) != len(obj2) {
		return false
	}

	// all keys in obj1 must exist in obj2
	for key, item1 := range obj1 {
		item2, ok := obj2[key]
		if !ok {
			return false
		}

		if reflect.TypeOf(item1) != reflect.TypeOf(item2) {
			return false
		}

		if reflect.TypeOf(item1) == reflect.TypeOf(bson.M{}) {
			ok = bsonMAreEqual(item1.(bson.M), item2.(bson.M))
		} else {
			ok = item1 == item2
		}

		if !ok {
			return false
		}
	} 

	return true
}

func TestConvertObjectPayload(t *testing.T) {
	// case 1: with 1 level of object has reached inside a map
	case1Field := field.StringField("field1")
	case1Depth := 0
	case1Payload := convertObjectPayload(case1Field, "$alias", field.TypeInt32, &case1Depth)
	case1ExpectedPayload := bson.M{
		"field1": bson.M{
			convertFunction(field.TypeString, field.TypeInt32): "$$alias.field1",
		},
	}

	test.AssertTrue(t, bsonMAreEqual(case1Payload, case1ExpectedPayload), "Case 1: Unexpected Payload")

	// case 2: with 2 level of object has reached inside a map
	case2Field := field.ObjectField("field1", field.StringField("field2"))
	case2Depth := 0
	case2Payload := convertObjectPayload(case2Field, "$alias", field.TypeInt32, &case2Depth)
	case2ExpectedPayload := bson.M{
		"field1": bson.M{
			"field2": bson.M{
				convertFunction(field.TypeString, field.TypeInt32): "$$alias.field1.field2",
			},
		},
	}

	test.AssertTrue(t, bsonMAreEqual(case2Payload, case2ExpectedPayload), "Case 2: Unexpected Payload")

	// TODO: add more cases
}

func TestConvertMapPayload(t *testing.T) {
	// case 1: plain string array
	case1Field := field.StringField("")
	case1Depth := 0
	case1Payload := convertMapPayload(case1Field, "scores", field.TypeInt32, &case1Depth)
	case1ExpectedPayload := bson.M{
		"$map": bson.M{
			"input": "$scores",
			"as": "alias_1",
			"in": bson.M{
				convertFunction(field.TypeInt32, field.TypeString): "$$alias_1",
			},
		},
	}
	
	test.AssertTrue(t, bsonMAreEqual(case1Payload, case1ExpectedPayload), "Case 1: Unexpected Payload")

	// case 2: array of object with string key
	case2Field := field.ObjectField("", field.StringField("score"))
	case2Depth := 0
	case2Payload := convertMapPayload(case2Field, "scores", field.TypeInt32, &case2Depth)
	case2ExpectedPayload := bson.M{
		"$map": bson.M{
			"input": "$scores",
			"as": "alias_1",
			"in": bson.M{
				"score": bson.M{
					convertFunction(field.TypeInt32, field.TypeString): "$$alias_1.score",
				},
			},
		},
	}
	
	test.AssertTrue(t, bsonMAreEqual(case2Payload, case2ExpectedPayload), "Case 2: Unexpected Payload")

	// case 3: array of plain string array
	case3Field := field.ArrayField("", field.StringField("score"))
	case3Depth := 0
	case3Payload := convertMapPayload(case3Field, "scores", field.TypeInt32, &case3Depth)
	case3ExpectedPayload := bson.M{
		"$map": bson.M{
			"input": "$scores",
			"as": "alias_1",
			"in": bson.M{
				"$map": bson.M{
					"input": "$$alias_1",
					"as": "alias_2",
					"in": bson.M{
						convertFunction(field.TypeInt32, field.TypeString): "$$alias_2",
					},
				},
			},
		},
	}
	
	test.AssertTrue(t, bsonMAreEqual(case3Payload, case3ExpectedPayload), "Case 3: Unexpected Payload")

	// TODO: add more cases
}

func TestConvertSetPayload(t *testing.T) {
	// case 1: plain string field conversion
	case1Field := field.StringField("field1")
	case1Depth := 0
	case1Payload := convertSetPayload(case1Field, "", field.TypeInt32, &case1Depth)
	case1ExpectedPayload := bson.M{
		"field1": bson.M{
			convertFunction(field.TypeInt32, field.TypeString): "$field1",
		},
	}
	
	test.AssertTrue(t, bsonMAreEqual(case1Payload, case1ExpectedPayload), "Case 1: Unexpected Payload")

	// case 2: plain string field in nested object conversion
	case2Field := field.ObjectField("field1", field.ObjectField("field2", field.StringField("field3")))
	case2Depth := 0
	case2Payload := convertSetPayload(case2Field, "", field.TypeInt32, &case2Depth)
	case2ExpectedPayload := bson.M{
		"field1.field2.field3": bson.M{
			convertFunction(field.TypeInt32, field.TypeString): "$field1.field2.field3",
		},
	}
	
	test.AssertTrue(t, bsonMAreEqual(case2Payload, case2ExpectedPayload), "Case 2: Unexpected Payload")

	// case 3: nested array field in nested object conversion
	case3Field := field.ObjectField("field1", 
		field.ObjectField("field2", 
			field.ArrayField("field3",
				field.ArrayField("not required",
					field.StringField("not required"),
				),
			),
		),
	)
	case3Depth := 0
	case3Payload := convertSetPayload(case3Field, "", field.TypeInt32, &case3Depth)
	case3ExpectedPayload := bson.M{
		"field1.field2.field3": bson.M{
			"$map": bson.M{
				"input": "$field1.field2.field3",
				"as": "alias_1",
				"in": bson.M{
					"$map": bson.M{
						"input": "$$alias_1",
						"as": "alias_2",
						"in": bson.M{
							convertFunction(field.TypeInt32, field.TypeString): "$$alias_2",
						},
					},
				},
			},
		},
	}
	
	test.AssertTrue(t, bsonMAreEqual(case3Payload, case3ExpectedPayload), "Case 3: Unexpected Payload")

	// case 4: nested array of nested object in nested object conversion
	case4Field := field.ObjectField("field1", 
		field.ObjectField("field2", 
			field.ArrayField("field3",
				field.ArrayField("not required",
					field.ObjectField("not required",
						field.ObjectField("field4",
							field.StringField("field5"),
						),
					),
				),
			),
		),
	)
	case4Depth := 0
	case4Payload := convertSetPayload(case4Field, "", field.TypeInt32, &case4Depth)
	case4ExpectedPayload := bson.M{
		"field1.field2.field3": bson.M{
			"$map": bson.M{
				"input": "$field1.field2.field3",
				"as": "alias_1",
				"in": bson.M{
					"$map": bson.M{
						"input": "$$alias_1",
						"as": "alias_2",
						"in": bson.M{
							"field4": bson.M{
								"field5": bson.M{
									convertFunction(field.TypeInt32, field.TypeString): "$$alias_2.field4.field5",
								},
							},
						},
					},
				},
			},
		},
	}
	
	test.AssertTrue(t, bsonMAreEqual(case4Payload, case4ExpectedPayload), "Case 4: Unexpected Payload")
}