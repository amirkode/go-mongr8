/*
Copyright (c) 2023-present the go-mongr8 Authors and Contributors
[@see Authors file]

Licensed under the MIT License
(https://opensource.org/licenses/MIT)
*/
package convert

import (
	"reflect"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func ArrToBsonA(arr []interface{}, isD bool) bson.A {
	res := bson.A{}
	for _, value := range arr {
		if reflect.TypeOf(value).Kind() == reflect.Map &&
			reflect.TypeOf(value).Key().Kind() == reflect.String &&
			reflect.TypeOf(value).Elem().Kind() == reflect.Interface {
			// if it's a map, recursively convert it to bson.M
			if isD {
				res = append(res, MapToBsonD(value.(map[string]interface{})))
			} else {
				res = append(res, MapToBsonM(value.(map[string]interface{})))
			}
		} else if reflect.TypeOf(value).Kind() == reflect.Slice &&
			reflect.TypeOf(value).Elem().Kind() == reflect.Interface {
			// if value type is an array ([]interface{})
			nextArr := value.([]interface{})
			res = append(res, ArrToBsonA(nextArr, isD))
		} else {
			// otherwise, use the value as is
			res = append(res, value)
		}
	}

	return res
}

func MapToBsonM(mp map[string]interface{}) bson.M {
	res := bson.M{}
	for key, value := range mp {
		// use reflect.TypeOf to check if the value is a map[string]interface{}
		if reflect.TypeOf(value).Kind() == reflect.Map &&
			reflect.TypeOf(value).Key().Kind() == reflect.String &&
			reflect.TypeOf(value).Elem().Kind() == reflect.Interface {
			// if it's a map, recursively convert it to bson.M
			res[key] = MapToBsonM(value.(map[string]interface{}))
		} else if reflect.TypeOf(value).Kind() == reflect.Slice &&
			reflect.TypeOf(value).Elem().Kind() == reflect.Interface {
			// if value type is an array ([]interface{})
			arr := value.([]interface{})
			res[key] = ArrToBsonA(arr, false)
		} else {
			// otherwise, use the value as is
			res[key] = value
		}
	}

	return res
}

func MapToBsonD(mp map[string]interface{}) bson.D {
	res := bson.D{}
	for key, value := range mp {
		// use reflect.TypeOf to check if the value is a map[string]interface{}
		if reflect.TypeOf(value).Kind() == reflect.Map &&
			reflect.TypeOf(value).Key().Kind() == reflect.String &&
			reflect.TypeOf(value).Elem().Kind() == reflect.Interface {
			// if it's a map, recursively convert it to bson.M
			res = append(res, bson.E{
				Key:   key,
				Value: MapToBsonD(value.(map[string]interface{})),
			})
		} else if reflect.TypeOf(value).Kind() == reflect.Slice &&
			reflect.TypeOf(value).Elem().Kind() == reflect.Interface {
			// if value type is an array ([]interface{})
			arr := value.([]interface{})
			res = append(res, bson.E{
				Key:   key,
				Value: ArrToBsonA(arr, true),
			})
		} else {
			// otherwise, use the value as is
			res = append(res, bson.E{
				Key:   key,
				Value: value,
			})
		}
	}

	return res
}

func ConvertBsonPrimitiveToDefaultType(value interface{}) interface{} {
	switch reflect.TypeOf(value) {
	case reflect.TypeOf(primitive.DateTime(0)):
		return value.(primitive.DateTime).Time()
	case reflect.TypeOf(primitive.ObjectID{}):
		return value.(primitive.ObjectID).Hex()
	case reflect.TypeOf(primitive.Binary{}):
		org := value.(primitive.Binary)
		data := make([]byte, len(org.Data))
		copy(data, org.Data)
		return data
	case reflect.TypeOf(primitive.Decimal128{}):
		return value.(primitive.Decimal128).String() // convert to string because it can be too long
	}

	return value
}
