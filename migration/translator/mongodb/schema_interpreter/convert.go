/*
Copyright (c) 2023 the go-mongr8 Authors and Contributors
[@see Authors file]

Licensed under the MIT License
(https://opensource.org/licenses/MIT)
*/
package schema_interpreter

import (
	"fmt"
	"reflect"
	"time"

	"github.com/amirkode/go-mongr8/migration/translator/dictionary"
)

// convert any value to literal, this function can be called any where
func AnyToLiteral(value interface{}) string {
	if reflect.TypeOf(value).Kind() == reflect.Map &&
		reflect.TypeOf(value).Key().Kind() == reflect.String {
		res := "map[string]interface{}{\n"
		for key, v := range value.(map[string]interface{}) {
			res += fmt.Sprintf(`"%s": %s,`, key, AnyToLiteral(v)) + "\n"
		}
		res += "}"

		return res
	} else if reflect.TypeOf(value).Kind() == reflect.Slice {
		res := "[]interface{}{\n"
		for _, v := range value.([]interface{}) {
			res += fmt.Sprintf("%s\n,", AnyToLiteral(v))
		}
		res += "}"

		return res
	}

	return anyToLiteralString(value)
}

// converts any value to its string literal type
// i.e: str := "hello", this function will return `string("hello")`
func anyToLiteralString(value interface{}) string {
	v := value
	if reflect.TypeOf(value) == reflect.TypeOf(dictionary.ValueType{}) {
		v = value.(dictionary.ValueType).Value
	}

	// handle primitives
	switch reflect.TypeOf(v).Kind() {
	case reflect.Int:
		return fmt.Sprintf("int(%v)", v)
	case reflect.Int8:
		return fmt.Sprintf("int8(%v)", v)
	case reflect.Int16:
		return fmt.Sprintf("int16(%v)", v)
	case reflect.Int32:
		return fmt.Sprintf("int32(%v)", v)
	case reflect.Int64:
		return fmt.Sprintf("int64(%v)", v)
	case reflect.Float32:
		return fmt.Sprintf("float32(%v)", v)
	case reflect.Float64:
		return fmt.Sprintf("float64(%v)", v)
	case reflect.String:
		return fmt.Sprintf(`string("%v")`, v)
	case reflect.Bool:
		return fmt.Sprintf("bool(%v)", v)
	}

	// handle non primitives
	switch reflect.TypeOf(v) {
	case reflect.TypeOf(time.Time{}):
		return timeToLiteralString(v.(time.Time))
	}

	// if none of type is recognized, just return as a string ValueType
	return fmt.Sprintf(`string("%v")`, v)
}

func timeToLiteralString(t time.Time) string {
	year, month, day := t.Date()
	hour, minute, second := t.Clock()
	zoneName, offset := t.Zone()

	return fmt.Sprintf("time.Date(%d, time.%s, %d, %d, %d, %d, %d, time.%s)",
		year, month, day, hour, minute, second, offset, zoneName)
}

// convert a map to literal bson.M map definition in string
func toLiteralStringBsonMap(value interface{}) string {
	if reflect.TypeOf(value).Kind() == reflect.Map &&
		reflect.TypeOf(value).Key().Kind() == reflect.String {
		res := "bson.M{\n"
		for key, v := range value.(map[string]interface{}) {
			res += fmt.Sprintf(`"%s": %s,`, key, toLiteralStringBsonMap(v)) + "\n"
		}
		res += "}"

		return res
	} else if reflect.TypeOf(value).Kind() == reflect.Slice {
		res := "bson.A{\n"
		for _, v := range value.([]interface{}) {
			res += fmt.Sprintf("%s\n,", toLiteralStringBsonMap(v))
		}
		res += "}"

		return res
	}

	return anyToLiteralString(value)
}

// ConvertValueTypeToRealType return value reversal of ValueType(d) structure
// as defined in migration/translator/dictionary
// TODO: find better approach in the future as the ValueType was intended
// because of to many type conversions needed
func ConvertValueTypeToRealType(value interface{}) interface{} {
	if reflect.TypeOf(value).Kind() == reflect.Map &&
		reflect.TypeOf(value).Key().Kind() == reflect.String &&
		reflect.TypeOf(value).Elem().Kind() == reflect.Interface {
		// if value type is a map[string]interface{}
		res := value.(map[string]interface{})
		for key, v := range res {
			res[key] = ConvertValueTypeToRealType(v)
		}

		return res
	} else if reflect.TypeOf(value).Kind() == reflect.Slice &&
		reflect.TypeOf(value).Elem().Kind() == reflect.Interface {
		// if value type is a slice ([]interface{})
		res := value.([]interface{})
		for index, v := range res {
			res[index] = ConvertValueTypeToRealType(v)
		}

		return res
	}

	if reflect.TypeOf(value) == reflect.TypeOf(dictionary.ValueType{}) {
		// convert individual value
		return value.(dictionary.ValueType).Value
	}

	// if none of type is recognized, just return as a string ValueType
	return fmt.Sprintf("%v", value)
}
