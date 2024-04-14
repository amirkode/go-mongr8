/*
Copyright (c) 2023-present the go-mongr8 Authors and Contributors
[@see Authors file]

Licensed under the MIT License
(https://opensource.org/licenses/MIT)
*/
package api_interpreter

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/amirkode/go-mongr8/internal/test"

	"github.com/amirkode/go-mongr8/collection/field"

	"go.mongodb.org/mongo-driver/bson"

	. "github.com/smartystreets/goconvey/convey"
)

func bsonAAreEqual(arr1, arr2 bson.A) bool {
	// expects same order of arr1 and arr2
	if len(arr1) != len(arr2) {
		return false
	}

	for index, item1 := range arr1 {
		item2 := arr2[index]
		if reflect.TypeOf(item1) == reflect.TypeOf(bson.M{}) {
			if !bsonMAreEqual(item1.(bson.M), item2.(bson.M)) {
				return false
			}
		} else if reflect.TypeOf(item1) == reflect.TypeOf(bson.A{}) {
			if !bsonAAreEqual(item1.(bson.A), item2.(bson.A)) {
				return false
			}
		} else if item1 != item2 {
			return false
		}
	}

	return true
}

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
		} else if reflect.TypeOf(item1) == reflect.TypeOf(bson.A{}) {
			ok = bsonAAreEqual(item1.(bson.A), item2.(bson.A))
		} else {
			ok = item1 == item2
		}

		if !ok {
			return false
		}
	}

	return true
}

func TestAppendPath(t *testing.T) {
	// case 1: default
	Convey("Case 1: Default", t, func() {
		Convey("Single Parent appended to a child", func() {
			parentPath := "field1"
			childPath := "child1"
			actualPath := appendPath(parentPath, childPath)
			expectedPath := "field1.child1"
			So(actualPath, ShouldEqual, expectedPath)
		})

		Convey("Multi ancestor appended to a child", func() {
			parentPath := "field1.field2"
			childPath := "child1"
			actualPath := appendPath(parentPath, childPath)
			expectedPath := "field1.field2.child1"
			So(actualPath, ShouldEqual, expectedPath)
		})
	})

	// case 2: empty parent
	Convey("Case 2: Empty parent", t, func() {
		Convey("No parent appended to a child", func() {
			parentPath := ""
			childPath := "child1"
			actualPath := appendPath(parentPath, childPath)
			expectedPath := "child1"
			So(actualPath, ShouldEqual, expectedPath)
		})
	})
}

func TestConvertToUpsertPath(t *testing.T) {
	// case 1: default
	Convey("Case 1: Default", t, func() {
		Convey("Array path with zero index", func() {
			originalPath := "field1.0.0.field2"
			actualPath := convertToUpsertPath(originalPath)
			expectedPath := "field1.$[].$[].field2"
			So(actualPath, ShouldEqual, expectedPath)
		})

		Convey("Array path with non zero index", func() {
			originalPath := "field1.4.24.field2"
			actualPath := convertToUpsertPath(originalPath)
			expectedPath := "field1.$[].$[].field2"
			So(actualPath, ShouldEqual, expectedPath)
		})
	})

	// case 2: Unexpected array index
	Convey("Case 2: Unexpected array index", t, func() {
		defer func() {
			if r := recover(); r != nil {
				Convey("Unexpected panic", func() {
					So(fmt.Sprintf("%v", r), ShouldContainSubstring, "Unexpected array index")
				})
			}
		}()

		originalPath := "field1.-1.10.field2"
		convertToUpsertPath(originalPath)
	})
}

func TestGetParentPath(t *testing.T) {
	// case 1: default
	Convey("Case 1: Default", t, func() {
		Convey("With parent", func() {
			originalPath := "field1.0"
			actualPath := getParentPath(originalPath)
			expectedPath := "field1"
			So(actualPath, ShouldEqual, expectedPath)
		})

		Convey("Without parent", func() {
			originalPath := "field1"
			actualPath := getParentPath(originalPath)
			expectedPath := ""
			So(actualPath, ShouldEqual, expectedPath)
		})
	})
}

func TestCheckPathExistPayloads(t *testing.T) {
	comparePaths := func(path1, path2 *[]bson.D) bool {
		if path1 == nil || path2 == nil {
			return path1 == nil && path2 == nil
		}

		if len(*path1) != len(*path2) {
			return false
		}

		for idx, path := range *path1 {
			if path[0].Key != (*path2)[idx][0].Key ||
				path[1].Value != (*path2)[idx][1].Value ||
				path[2].Value != (*path2)[idx][2].Value {
				return false
			}
		}

		return true
	}
	// case 1: Default
	Convey("Case 1: Default", t, func() {
		Convey("Nesting array with object", func() {
			payload := bson.D{
				{Key: "arr", Value: bson.A{
					bson.D{
						{Key: "inner_arr", Value: bson.A{
							bson.D{
								{Key: "new_field", Value: ""},
							},
						}},
					},
				}},
			}
			actualPaths := []bson.D{}
			checkPathExistPayloads(payload, "", &actualPaths)
			expectedPaths := []bson.D{
				{
					{Key: "arr", Value: bson.M{"$exists": true}},
					{Key: "is_array", Value: false},
					{Key: "wants_array", Value: true},
				},
				{
					{Key: "arr.0", Value: bson.M{"$exists": true}},
					{Key: "is_array", Value: true},
					{Key: "wants_array", Value: false},
				},
				{
					{Key: "arr.0.inner_arr", Value: bson.M{"$exists": true}},
					{Key: "is_array", Value: false},
					{Key: "wants_array", Value: true},
				},
				{
					{Key: "arr.0.inner_arr.0", Value: bson.M{"$exists": true}},
					{Key: "is_array", Value: true},
					{Key: "wants_array", Value: false},
				},
			}
			So(comparePaths(&actualPaths, &expectedPaths), ShouldBeTrue)
		})

		Convey("Nesting array with another array", func() {
			payload := bson.D{
				{Key: "arr", Value: bson.A{
					bson.A{
						bson.D{
							{Key: "new_field", Value: ""},
						},
					},
				}},
			}
			actualPaths := []bson.D{}
			checkPathExistPayloads(payload, "", &actualPaths)
			expectedPaths := []bson.D{
				{
					{Key: "arr", Value: bson.M{"$exists": true}},
					{Key: "is_array", Value: false},
					{Key: "wants_array", Value: true},
				},
				{
					{Key: "arr.0", Value: bson.M{"$exists": true}},
					{Key: "is_array", Value: true},
					{Key: "wants_array", Value: true},
				},
				{
					{Key: "arr.0.0", Value: bson.M{"$exists": true}},
					{Key: "is_array", Value: true},
					{Key: "wants_array", Value: false},
				},
			}
			So(comparePaths(&actualPaths, &expectedPaths), ShouldBeTrue)
		})
	})
}

func TestCreateFieldSetPayload(t *testing.T) {
	// case 1: simple addition
	// {
	//   "new_string_field"	: ""
	// }
	case1Field := bson.D{{Key: "new_string_field", Value: ""}}
	case1Payload := createFieldSetPayload(case1Field, "")
	case1ExpectedPayload := bson.M{
		"new_string_field": "",
	}

	test.AssertTrue(t, bsonMAreEqual(case1Payload, case1ExpectedPayload), "Case 1: Unexpected Payload")

	// case 2: addition inside array
	// {
	// 	"arr": [
	// 		{
	// 			"new_field": true
	// 		}
	// 	]
	// }
	case2Field := bson.D{
		{Key: "arr", Value: bson.A{
			bson.D{
				{Key: "new_field", Value: true},
			},
		}},
	}
	case2Payload := createFieldSetPayload(case2Field, "")
	case2ExpectedPayload := bson.M{
		"arr.$[].new_field": true,
	}

	test.AssertTrue(t, bsonMAreEqual(case2Payload, case2ExpectedPayload), "Case 2: Unexpected Payload")

	// case 3: addition inside nested array
	// {
	// 	"arr": [
	// 		{
	// 			"inner_arr": [
	// 				{
	// 					"new_field": ""
	// 				}
	// 			]
	// 		}
	// 	]
	// }
	case3Field := bson.D{
		{Key: "arr", Value: bson.A{
			bson.D{
				{Key: "inner_arr", Value: bson.A{
					bson.D{
						{Key: "new_field", Value: ""},
					},
				}},
			},
		}},
	}
	case3Payload := createFieldSetPayload(case3Field, "")
	case3ExpectedPayload := bson.M{
		"arr.$[].inner_arr.$[].new_field": "",
	}

	test.AssertTrue(t, bsonMAreEqual(case3Payload, case3ExpectedPayload), "Case 3: Unexpected Payload")

	// case 4: addition of array field inside nested array
	// {
	// 	"arr": [
	// 		{
	// 			"inner_arr": [
	// 				{
	// 					"new_field": [0]
	// 				}
	// 			]
	// 		}
	// 	]
	// }
	case4Field := bson.D{
		{Key: "arr", Value: bson.A{
			bson.D{
				{Key: "inner_arr", Value: bson.A{
					bson.D{
						{Key: "new_field", Value: bson.A{0}},
					},
				}},
			},
		}},
	}
	case4Payload := createFieldSetPayload(case4Field, "")
	case4ExpectedPayload := bson.M{
		"arr.$[].inner_arr.$[].new_field": bson.A{0},
	}

	test.AssertTrue(t, bsonMAreEqual(case4Payload, case4ExpectedPayload), "Case 4: Unexpected Payload")

	// case 4: array of array
	// {
	// 	 "arr": [
	//	    [
	// 		     [0]
	// 		]
	// 	 ]
	// }
	case5Field := bson.D{
		{Key: "arr", Value: bson.A{
			bson.A{
				bson.A{0},
			},
		}},
	}
	case5Payload := createFieldSetPayload(case5Field, "")
	case5ExpectedPayload := bson.M{
		"arr.$[].$[]": bson.A{0},
	}

	test.AssertTrue(t, bsonMAreEqual(case5Payload, case5ExpectedPayload), "Case 5: Unexpected Payload")

	// TODO: add more cases
}

func TestDropFieldUnsetPayload(t *testing.T) {
	// case 1: simple addition
	// {
	//   "new_string_field"	: ""
	// }
	case1Field := bson.D{{Key: "new_string_field", Value: ""}}
	case1Payload := dropFieldUnsetPayload(case1Field, "")
	case1ExpectedPayload := bson.M{
		"new_string_field": "",
	}

	test.AssertTrue(t, bsonMAreEqual(case1Payload, case1ExpectedPayload), "Case 1: Unexpected Payload")

	// case 2: addition inside array
	// {
	// 	"arr": [
	// 		{
	// 			"new_field": true
	// 		}
	// 	]
	// }
	case2Field := bson.D{
		{Key: "arr", Value: bson.A{
			bson.D{
				{Key: "new_field", Value: true},
			},
		}},
	}
	case2Payload := dropFieldUnsetPayload(case2Field, "")
	case2ExpectedPayload := bson.M{
		"arr.$[].new_field": "",
	}

	test.AssertTrue(t, bsonMAreEqual(case2Payload, case2ExpectedPayload), "Case 2: Unexpected Payload")

	// case 3: addition inside nested array
	// {
	// 	"arr": [
	// 		{
	// 			"inner_arr": [
	// 				{
	// 					"new_field": ""
	// 				}
	// 			]
	// 		}
	// 	]
	// }
	case3Field := bson.D{
		{Key: "arr", Value: bson.A{
			bson.D{
				{Key: "inner_arr", Value: bson.A{
					bson.D{
						{Key: "new_field", Value: ""},
					},
				}},
			},
		}},
	}
	case3Payload := dropFieldUnsetPayload(case3Field, "")
	case3ExpectedPayload := bson.M{
		"arr.$[].inner_arr.$[].new_field": "",
	}

	test.AssertTrue(t, bsonMAreEqual(case3Payload, case3ExpectedPayload), "Case 3: Unexpected Payload")

	// case 4: addition of array field inside nested array
	// {
	// 	"arr": [
	// 		{
	// 			"inner_arr": [
	// 				{
	// 					"new_field": [0]
	// 				}
	// 			]
	// 		}
	// 	]
	// }
	case4Field := bson.D{
		{Key: "arr", Value: bson.A{
			bson.D{
				{Key: "inner_arr", Value: bson.A{
					bson.D{
						{Key: "new_field", Value: bson.A{0}},
					},
				}},
			},
		}},
	}
	case4Payload := dropFieldUnsetPayload(case4Field, "")
	case4ExpectedPayload := bson.M{
		"arr.$[].inner_arr.$[].new_field": "",
	}

	test.AssertTrue(t, bsonMAreEqual(case4Payload, case4ExpectedPayload), "Case 4: Unexpected Payload")

	// TODO: add more cases
}

func TestConvertFunction(t *testing.T) {
	// case 1: default
	Convey("Case 1: Default", t, func() {
		functionMap := map[field.FieldType]string{
			field.TypeString:    "$toString",
			field.TypeBoolean:   "$toBool",
			field.TypeTimestamp: "$toDate",
			field.TypeInt32:     "$toInt",
			field.TypeInt64:     "$toLong",
			field.TypeDouble:    "$toDouble",
		}

		Convey("String Conversion", func() {
			So(functionMap[field.TypeString], ShouldEqual, convertFunction(field.TypeString, field.TypeString))
		})
		Convey("Boolean Conversion", func() {
			So(functionMap[field.TypeBoolean], ShouldEqual, convertFunction(field.TypeBoolean, field.TypeString))
		})
		Convey("Timestamp Conversion", func() {
			So(functionMap[field.TypeTimestamp], ShouldEqual, convertFunction(field.TypeTimestamp, field.TypeString))
		})
		Convey("Int32 Conversion", func() {
			So(functionMap[field.TypeInt32], ShouldEqual, convertFunction(field.TypeInt32, field.TypeString))
		})
		Convey("Int64 Conversion", func() {
			So(functionMap[field.TypeInt64], ShouldEqual, convertFunction(field.TypeInt64, field.TypeString))
		})
		Convey("Double Conversion", func() {
			So(functionMap[field.TypeDouble], ShouldEqual, convertFunction(field.TypeDouble, field.TypeString))
		})
	})

	// case 2: unsupported conversion
	Convey("Case 2: Unsupported conversion", t, func() {
		defer func() {
			if r := recover(); r != nil {
				Convey("Unexpected panic", func() {
					So(fmt.Sprintf("%v", r), ShouldContainSubstring, "not supported")
				})
			}
		}()

		convertFunction(field.TypeGeoJSONPoint, field.TypeString)
	})
}

func TestConvertFieldObjectPayload(t *testing.T) {
	// case 1: with 1 level of object has reached inside a map
	case1Field := field.StringField("field1")
	case1Depth := 0
	case1Payload := convertFieldObjectPayload(case1Field, "$alias", field.TypeInt32, &case1Depth)
	case1ExpectedPayload := bson.M{
		"field1": bson.M{
			convertFunction(field.TypeString, field.TypeInt32): "$$alias.field1",
		},
	}

	test.AssertTrue(t, bsonMAreEqual(case1Payload, case1ExpectedPayload), "Case 1: Unexpected Payload")

	// case 2: with 2 level of object has reached inside a map
	case2Field := field.ObjectField("field1", field.StringField("field2"))
	case2Depth := 0
	case2Payload := convertFieldObjectPayload(case2Field, "$alias", field.TypeInt32, &case2Depth)
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

func TestConvertFieldMapPayload(t *testing.T) {
	// case 1: plain string array
	case1Field := field.StringField("")
	case1Depth := 0
	case1Payload := convertFieldMapPayload(case1Field, "scores", field.TypeInt32, &case1Depth)
	case1ExpectedPayload := bson.M{
		"$map": bson.M{
			"input": "$scores",
			"as":    "alias_1",
			"in": bson.M{
				convertFunction(field.TypeString, field.TypeInt32): "$$alias_1",
			},
		},
	}

	test.AssertTrue(t, bsonMAreEqual(case1Payload, case1ExpectedPayload), "Case 1: Unexpected Payload")

	// case 2: array of object with string key
	case2Field := field.ObjectField("", field.StringField("score"))
	case2Depth := 0
	case2Payload := convertFieldMapPayload(case2Field, "scores", field.TypeInt32, &case2Depth)
	case2ExpectedPayload := bson.M{
		"$map": bson.M{
			"input": "$scores",
			"as":    "alias_1",
			"in": bson.M{
				"$mergeObjects": bson.A{
					"$$alias_1",
					bson.M{
						"score": bson.M{
							convertFunction(field.TypeString, field.TypeInt32): "$$alias_1.score",
						},
					},
				},
			},
		},
	}

	test.AssertTrue(t, bsonMAreEqual(case2Payload, case2ExpectedPayload), "Case 2: Unexpected Payload")

	// case 3: array of plain string array
	case3Field := field.ArrayField("", field.StringField("score"))
	case3Depth := 0
	case3Payload := convertFieldMapPayload(case3Field, "scores", field.TypeInt32, &case3Depth)
	case3ExpectedPayload := bson.M{
		"$map": bson.M{
			"input": "$scores",
			"as":    "alias_1",
			"in": bson.M{
				"$map": bson.M{
					"input": "$$alias_1",
					"as":    "alias_2",
					"in": bson.M{
						convertFunction(field.TypeString, field.TypeInt32): "$$alias_2",
					},
				},
			},
		},
	}

	test.AssertTrue(t, bsonMAreEqual(case3Payload, case3ExpectedPayload), "Case 3: Unexpected Payload")

	// TODO: add more cases
}

func TestConvertFieldSetPayload(t *testing.T) {
	// case 1: plain string field conversion
	case1Field := field.StringField("field1")
	case1Depth := 0
	case1Payload := convertFieldSetPayload(case1Field, "", field.TypeInt32, &case1Depth)
	case1ExpectedPayload := bson.M{
		"field1": bson.M{
			convertFunction(field.TypeString, field.TypeInt32): "$field1",
		},
	}

	test.AssertTrue(t, bsonMAreEqual(case1Payload, case1ExpectedPayload), "Case 1: Unexpected Payload")

	// case 2: plain string field in nested object conversion
	case2Field := field.ObjectField("field1", field.ObjectField("field2", field.StringField("field3")))
	case2Depth := 0
	case2Payload := convertFieldSetPayload(case2Field, "", field.TypeInt32, &case2Depth)
	case2ExpectedPayload := bson.M{
		"field1.field2.field3": bson.M{
			convertFunction(field.TypeString, field.TypeInt32): "$field1.field2.field3",
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
	case3Payload := convertFieldSetPayload(case3Field, "", field.TypeInt32, &case3Depth)
	case3ExpectedPayload := bson.M{
		"field1.field2.field3": bson.M{
			"$map": bson.M{
				"input": "$field1.field2.field3",
				"as":    "alias_1",
				"in": bson.M{
					"$map": bson.M{
						"input": "$$alias_1",
						"as":    "alias_2",
						"in": bson.M{
							convertFunction(field.TypeString, field.TypeInt32): "$$alias_2",
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
	case4Payload := convertFieldSetPayload(case4Field, "", field.TypeInt32, &case4Depth)
	case4ExpectedPayload := bson.M{
		"field1.field2.field3": bson.M{
			"$map": bson.M{
				"input": "$field1.field2.field3",
				"as":    "alias_1",
				"in": bson.M{
					"$map": bson.M{
						"input": "$$alias_1",
						"as":    "alias_2",
						"in": bson.M{
							"$mergeObjects": bson.A{
								"$$alias_2",
								bson.M{
									"field4": bson.M{
										"field5": bson.M{
											convertFunction(field.TypeString, field.TypeInt32): "$$alias_2.field4.field5",
										},
									},
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
