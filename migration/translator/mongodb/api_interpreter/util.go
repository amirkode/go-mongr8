package api_interpreter

import (
	"fmt"
	"reflect"

	"github.com/amirkode/go-mongr8/collection"
	"github.com/amirkode/go-mongr8/collection/field"

	"go.mongodb.org/mongo-driver/bson"
)

func appendPath(parent, child string) string {
	if parent == "" {
		return child
	}

	return fmt.Sprintf("%s.%s", parent, child)
}

func createFieldSetPayload(curr interface{}, path string) bson.M {
	if reflect.TypeOf(curr) == reflect.TypeOf(bson.D{}) {
		d := curr.(bson.D)
		return createFieldSetPayload(d[0].Value, appendPath(path, d[0].Key))
	} else if reflect.TypeOf(curr) == reflect.TypeOf(bson.A{}) {
		a := curr.(bson.A)
		// if there's no deeper search, then just set the current path to the array
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
		return "$toBoolean"
	case field.TypeTimestamp:
		return "$toDate"
	case field.TypeInt32:
		return "$toInt"
	case field.TypeInt64:
		return "$toLong"
	// TODO: complete for the future usecases
	}

	panic("Conversion is not supported")
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
// `curr` represents the current instance of Field
// `path` represents the current path of fields so far
// `from` represents the type of conversion from
// `depth` represents the the depth of map operations has reached
func convertFieldSetPayload(curr collection.Field, path string, from field.FieldType, depth *int) bson.M {
	// TODO: FIXME: by default, this will ignore other properties other than the converted field
	// it's because of the $map behaviour
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
