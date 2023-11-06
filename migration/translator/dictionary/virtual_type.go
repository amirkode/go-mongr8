package dictionary

import (
	"fmt"
	"reflect"
	"time"
)

type (
	DataType  string
	ValueType struct {
		Type  DataType
		Value interface{}
	}
)

// list of primitive and non-primitive data types
const (
	DataTypeInt     = "int"
	DataTypeInt8    = "int8"
	DataTypeInt16   = "int16"
	DataTypeInt32   = "int32"
	DataTypeInt64   = "int64"
	DataTypeFloat32 = "float32"
	DataTypeFloat64 = "float64"
	DataTypeString  = "string"
	DataTypeBoolean = "boolean"
	DataTypeTime    = "time"
	// TODO: add more types
	// but, so far, not all data types are needed to initialize the mongodb field
)

func Int(value int) ValueType {
	return ValueType{
		Type:  DataTypeInt,
		Value: value,
	}
}

func Int8(value int8) ValueType {
	return ValueType{
		Type:  DataTypeInt8,
		Value: value,
	}
}

func Int16(value int16) ValueType {
	return ValueType{
		Type:  DataTypeInt16,
		Value: value,
	}
}

func Int32(value int32) ValueType {
	return ValueType{
		Type:  DataTypeInt32,
		Value: value,
	}
}

func Int64(value int64) ValueType {
	return ValueType{
		Type:  DataTypeInt64,
		Value: value,
	}
}

func Float32(value float32) ValueType {
	return ValueType{
		Type:  DataTypeFloat32,
		Value: value,
	}
}

func Float64(value float64) ValueType {
	return ValueType{
		Type:  DataTypeFloat64,
		Value: value,
	}
}

func String(value string) ValueType {
	return ValueType{
		Type:  DataTypeString,
		Value: value,
	}
}

func Boolean(value bool) ValueType {
	return ValueType{
		Type:  DataTypeBoolean,
		Value: value,
	}
}

func Time(value time.Time) ValueType {
	return ValueType{
		Type:  DataTypeTime,
		Value: value,
	}
}

// helpers
func Array(children ...interface{}) []interface{} {
	return []interface{}{children}
}

// keep all elements/payloads in the same logical level type
// to ensure the integrity of data if other layer might have breaking changes
func ConvertAnyToValueType(value interface{}) interface{} {
	if reflect.TypeOf(value).Kind() == reflect.Map &&
		reflect.TypeOf(value).Key().Kind() == reflect.String &&
		reflect.TypeOf(value).Elem().Kind() == reflect.Interface {
		// if value type is a map[string]interface{}
		res := value.(map[string]interface{})
		for key, v := range res {
			res[key] = ConvertAnyToValueType(v)
		}

		return res
	} else if reflect.TypeOf(value).Kind() == reflect.Slice &&
		reflect.TypeOf(value).Elem().Kind() == reflect.Interface {
		// if value type is an array ([]interface{})
		res := value.([]interface{})
		for index, v := range res {
			res[index] = ConvertAnyToValueType(v)
		}

		return res
	}

	// convert individual value

	// handle primitives
	switch reflect.TypeOf(value).Kind() {
	case reflect.Int:
		return Int(value.(int))
	case reflect.Int8:
		return Int8(value.(int8))
	case reflect.Int16:
		return Int16(value.(int16))
	case reflect.Int32:
		return Int32(value.(int32))
	case reflect.Int64:
		return Int64(value.(int64))
	case reflect.Float32:
		return Float32(value.(float32))
	case reflect.Float64:
		return Float64(value.(float64))
	case reflect.String:
		return String(value.(string))
	case reflect.Bool:
		return Boolean(value.(bool))
	}

	// handle non primitives
	switch reflect.TypeOf(value) {
	case reflect.TypeOf(time.Time{}):
		return Time(value.(time.Time))
	}

	// if none of type is recognized, just return as a string ValueType
	return String(fmt.Sprintf("%v", value))
}
