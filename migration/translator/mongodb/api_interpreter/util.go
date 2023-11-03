package api_interpreter

import (
	"fmt"

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

func convertFunction(to field.FieldType, from field.FieldType) string {
	return "implement this"
}

// This returns the bosn.M object payload of a map
func convertObjectPayload(curr collection.Field, path string, from field.FieldType, depth *int) bson.M {
	var child bson.M
	switch curr.Spec().Type {
	case field.TypeArray:
		child = convertMapPayload(field.FromFieldSpec(&(*curr.Spec().ArrayFields)[0]), appendPath(path, curr.Spec().Name), from, depth)
	case field.TypeObject:
		child = convertObjectPayload(field.FromFieldSpec(&(*curr.Spec().Object)[0]), appendPath(path, curr.Spec().Name), from, depth)
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
func convertMapPayload(curr collection.Field, path string, from field.FieldType, depth *int) bson.M {
	var child bson.M
	*depth += 1
	currAlias := fmt.Sprintf("alias_%d", *depth)
	switch curr.Spec().Type {
	case field.TypeArray:
		child = convertMapPayload(field.FromFieldSpec(&(*curr.Spec().ArrayFields)[0]), fmt.Sprintf("$%s", currAlias), from, depth)
	case field.TypeObject:
		// this must be a child of map operation
		child = convertObjectPayload(field.FromFieldSpec(&(*curr.Spec().Object)[0]), fmt.Sprintf("$%s", currAlias), from, depth)
	default:
		child = bson.M{
			convertFunction(curr.Spec().Type, from): fmt.Sprintf("$$%s", currAlias),
		}
	}

	return bson.M{
		"$map": bson.M{
			"input": fmt.Sprintf("$%s", path),
			"as": currAlias,
			"in": child,
		},
	}
}

// This returns the payload of conversions
// `curr` represents the current instance of Field
// `path` represents the current path of fields so far
// `from` represents the type of conversion from
// `depth` represents the the depth of map operations has reached
func convertSetPayload(curr collection.Field, path string, from field.FieldType, depth *int) bson.M {
	currPath := appendPath(path, curr.Spec().Name)
	switch curr.Spec().Type {
	case field.TypeArray:
		return bson.M{
			currPath: convertMapPayload(field.FromFieldSpec(&(*curr.Spec().ArrayFields)[0]), appendPath(path, curr.Spec().Name), from, depth),
		}
	case field.TypeObject:
		return convertSetPayload(field.FromFieldSpec(&(*curr.Spec().Object)[0]), currPath, from, depth)
	}
	
	return bson.M{
		currPath: bson.M{
			convertFunction(curr.Spec().Type, from): fmt.Sprintf("$%s", currPath),
		},
	}
}
