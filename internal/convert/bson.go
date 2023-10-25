package convert

import (
	"reflect"

	"go.mongodb.org/mongo-driver/bson"
)

func MapToBsonM(mp map[string]interface{}) bson.M {
	res := bson.M{}
	for key, value := range mp {
		// use reflect.TypeOf to check if the value is a map[string]interface{}
		if reflect.TypeOf(value).Kind() == reflect.Map &&
			reflect.TypeOf(value).Key().Kind() == reflect.String &&
			reflect.TypeOf(value).Elem().Kind() == reflect.Interface {
			// if it's a map, recursively convert it to bson.M
			res[key] = MapToBsonM(value.(map[string]interface{}))
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
