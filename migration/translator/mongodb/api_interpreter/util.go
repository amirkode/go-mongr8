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
	"strconv"
	"strings"

	"github.com/amirkode/go-mongr8/collection"
	"github.com/amirkode/go-mongr8/collection/field"
	"github.com/amirkode/go-mongr8/internal/test"

	"go.mongodb.org/mongo-driver/bson"
)

func appendPath(parent, child string) string {
	if parent == "" {
		return child
	}

	return fmt.Sprintf("%s.%s", parent, child)
}

// This function converts any path to valid upsert path
// for instance, conversion from path checking:
// field.0.0 -> field.$[].$[]
// it means expecting upsertion for all indexes in the array
func convertToUpsertPath(path string) string {
	splitted := strings.Split(path, ".")
	// check and convert any index path to all indexes specifier
	for idx, curr := range splitted {
		if value, err := strconv.Atoi(curr); err == nil {
			if value < 0 {
				panic(fmt.Sprintf("Unexpected array index: %d", value))
			}
			// if current path is an index or number
			splitted[idx] = "$[]"
		}
	}

	return strings.Join(splitted, ".")
}

func getParentPath(path string) string {
	splitted := strings.Split(path, ".")
	return strings.Join(splitted[:len(splitted)-1], ".")
}

// This functions returns all possible paths based on the payload
func checkPathExistPayloads(curr interface{}, path string, res *[]bson.D) {
	if reflect.TypeOf(curr) == reflect.TypeOf(bson.D{}) {
		d := curr.(bson.D)
		if reflect.TypeOf(d[0].Value) == reflect.TypeOf(bson.A{}) ||
			reflect.TypeOf(d[0].Value) == reflect.TypeOf(bson.D{}) {
			path = appendPath(path, d[0].Key)
			*res = append(*res, bson.D{
				{Key: path, Value: bson.M{"$exists": true}},
				{Key: "is_array", Value: false},
				{Key: "wants_array", Value: reflect.TypeOf(d[0].Value) == reflect.TypeOf(bson.A{})},
			})
			checkPathExistPayloads(d[0].Value, path, res)
		}
	} else if reflect.TypeOf(curr) == reflect.TypeOf(bson.A{}) {
		a := curr.(bson.A)
		if reflect.TypeOf(a[0]) == reflect.TypeOf(bson.A{}) ||
			reflect.TypeOf(a[0]) == reflect.TypeOf(bson.D{}) {
			path = appendPath(path, "0")
			*res = append(*res, bson.D{
				{Key: path, Value: bson.M{"$exists": true}},
				{Key: "is_array", Value: true},
				{Key: "wants_array", Value: reflect.TypeOf(a[0]) == reflect.TypeOf(bson.A{})},
			})
			checkPathExistPayloads(a[0], path, res)
		}
	}
	// the deepest key (not an array nor an object) won't be checked
}

// This functions recursively construct payload for field creation
// any nested field expected to be only one way, for example:
// correct:
//
//	{
//		 "field": {
//	    "sub_field1": {
//	       "sub_field2": "a value goes here"
//	     }
//	  }
//	}
//
// wrong:
//
//	{
//		 "field": {
//	    "sub_field1": {
//	       "sub_field2": "a value goes here"
//	     }
//	    "sub_field1": {
//	       "sub_field2_1": "a value goes here",
//	       "sub_field2_2": 0
//	     }
//	  }
//	}
//
// this because it's guaranteed that any payload passed here
// must be a one way path as assured in migration generation
// @see migration/translation/sync_strategy/ for more clarity
//
// Parameters:
// `curr` represents the current payload yet to explore
// `path` represents the path has been explored so far
func createFieldSetPayload(curr interface{}, path string) bson.M {
	if reflect.TypeOf(curr) == reflect.TypeOf(bson.D{}) {
		d := curr.(bson.D)
		test.Assert(len(d) == 1, "createFieldSetPayload", "Object is not one way path")
		return createFieldSetPayload(d[0].Value, appendPath(path, d[0].Key))
	} else if reflect.TypeOf(curr) == reflect.TypeOf(bson.A{}) {
		a := curr.(bson.A)
		test.Assert(len(a) == 1, "createFieldSetPayload", "Array is not one way path")
		// if there's no deeper search, then just set the current path value to the array
		if reflect.TypeOf(a[0]) == reflect.TypeOf(bson.A{}) ||
			reflect.TypeOf(a[0]) == reflect.TypeOf(bson.D{}) {
			return createFieldSetPayload(a[0], appendPath(path, "$[]"))
		}
	}

	return bson.M{
		path: curr,
	}
}

func dropFieldUnsetPayload(curr interface{}, path string) bson.M {
	if reflect.TypeOf(curr) == reflect.TypeOf(bson.D{}) {
		d := curr.(bson.D)
		return dropFieldUnsetPayload(d[0].Value, appendPath(path, d[0].Key))
	} else if reflect.TypeOf(curr) == reflect.TypeOf(bson.A{}) {
		a := curr.(bson.A)
		// if there's no deeper search, then just set the current path to the array
		if reflect.TypeOf(a[0]) == reflect.TypeOf(bson.A{}) ||
			reflect.TypeOf(a[0]) == reflect.TypeOf(bson.D{}) {
			return dropFieldUnsetPayload(a[0], appendPath(path, "$[]"))
		}
	}

	return bson.M{
		path: "",
	}
}

// This returns original conversion function in MongoDB
// assuming all the conversions are valid
func convertFunction(to field.FieldType, from field.FieldType) string {
	// TODO: utilize `from` param in the future if needed
	switch to {
	case field.TypeString:
		return "$toString"
	case field.TypeBoolean:
		return "$toBool"
	case field.TypeTimestamp:
		return "$toDate"
	case field.TypeInt32:
		return "$toInt"
	case field.TypeInt64:
		return "$toLong"
	case field.TypeDouble:
		return "$toDouble"
		// TODO: complete for the future usecases
	}

	panic(fmt.Sprintf("Conversion from %s to %s is not supported", from, to))
}

// This returns the bosn.M object payload of a map
func convertFieldObjectPayload(curr collection.Field, path string, from field.FieldType, depth *int) bson.M {
	var child bson.M
	switch curr.Spec().Type {
	case field.TypeArray:
		child = convertFieldMapPayload(field.FromFieldSpec(&(*curr.Spec().ArrayFields)[0]), appendPath(path, curr.Spec().Name), from, depth)
	case field.TypeObject:
		child = convertFieldObjectPayload(field.FromFieldSpec(&(*curr.Spec().Object)[0]), appendPath(path, curr.Spec().Name), from, depth)
	default:
		child = bson.M{
			convertFunction(curr.Spec().Type, from): fmt.Sprintf("$%s", appendPath(path, curr.Spec().Name)),
		}
	}

	return bson.M{
		curr.Spec().Name: child,
	}
}

// This returns map operation in the bson.M representation
func convertFieldMapPayload(curr collection.Field, path string, from field.FieldType, depth *int) bson.M {
	var child bson.M
	*depth += 1
	currAlias := fmt.Sprintf("alias_%d", *depth)
	switch curr.Spec().Type {
	case field.TypeArray:
		child = convertFieldMapPayload(field.FromFieldSpec(&(*curr.Spec().ArrayFields)[0]), fmt.Sprintf("$%s", currAlias), from, depth)
	case field.TypeObject:
		// this must be a child of map operation
		child = convertFieldObjectPayload(field.FromFieldSpec(&(*curr.Spec().Object)[0]), fmt.Sprintf("$%s", currAlias), from, depth)
	default:
		child = bson.M{
			convertFunction(curr.Spec().Type, from): fmt.Sprintf("$$%s", currAlias),
		}
	}

	mp := bson.M{
		"input": fmt.Sprintf("$%s", path),
		"as":    currAlias,
		"in":    child,
	}

	// merge objects
	if curr.Spec().Type == field.TypeObject {
		mp["in"] = bson.M{
			"$mergeObjects": bson.A{
				fmt.Sprintf("$$%s", currAlias),
				child,
			},
		}
	}

	return bson.M{
		"$map": mp,
	}
}

// This returns the payload of conversions
// Parameters:
// `curr` represents the current instance of Field
// `path` represents the current path of fields so far
// `from` represents the type of conversion from
// `depth` represents the the depth of map operations has reached
func convertFieldSetPayload(curr collection.Field, path string, from field.FieldType, depth *int) bson.M {
	currPath := appendPath(path, curr.Spec().Name)
	switch curr.Spec().Type {
	case field.TypeArray:
		return bson.M{
			currPath: convertFieldMapPayload(field.FromFieldSpec(&(*curr.Spec().ArrayFields)[0]), appendPath(path, curr.Spec().Name), from, depth),
		}
	case field.TypeObject:
		return convertFieldSetPayload(field.FromFieldSpec(&(*curr.Spec().Object)[0]), currPath, from, depth)
	}

	return bson.M{
		currPath: bson.M{
			convertFunction(curr.Spec().Type, from): fmt.Sprintf("$%s", currPath),
		},
	}
}
