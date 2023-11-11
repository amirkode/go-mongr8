/*
Copyright (c) 2023 the go-mongr8 Authors and Contributors
[@see Authors file]

Licensed under the MIT License
(https://opensource.org/licenses/MIT)
*/
package dictionary

import (
	"time"

	"github.com/amirkode/go-mongr8/collection"
	"github.com/amirkode/go-mongr8/collection/field"
)

type (
	TranslatedField struct {
		TranslatedFieldIf
		field collection.Field
	}

	translatedString struct {
		TranslatedField
	}

	translatedInt32 struct {
		TranslatedField
	}

	translatedInt64 struct {
		TranslatedField
	}

	translatedDouble struct {
		TranslatedField
	}

	translatedBoolean struct {
		TranslatedField
	}

	translatedArray struct {
		TranslatedField
	}

	translatedObject struct {
		TranslatedField
	}

	translatedTimestamp struct {
		TranslatedField
	}

	// base geo json
	translatedGeoJSON struct {
		TranslatedField
	}

	translatedGeoJSONPoint struct {
		translatedGeoJSON
	}

	translatedGeoJSONLineString struct {
		translatedGeoJSON
	}

	translatedGeoJSONPolygonSingleRing struct {
		translatedGeoJSON
	}

	translatedGeoJSONPolygonMultipleRing struct {
		translatedGeoJSON
	}

	translatedGeoJSONMultiPoint struct {
		translatedGeoJSON
	}

	translatedGeoJSONMultiLineString struct {
		translatedGeoJSON
	}

	translatedGeoJSONMultiPolygon struct {
		translatedGeoJSON
	}

	translatedGeoJSONGeometryCollection struct {
		translatedGeoJSON
	}

	translatedLegacyCoordinate struct {
		TranslatedField
	}
)

// translation for string field type
func newTranslatedString(field collection.Field) translatedString {
	return translatedString{
		TranslatedField: TranslatedField{
			field: field,
		},
	}
}

func (t translatedString) GetObject() map[string]interface{} {
	key := t.field.Spec().Name
	return map[string]interface{}{
		key: String(""),
	}
}

func (t translatedString) GetArray() []interface{} {
	return Array()
}

// translation for int32 field type
func newTranslatedInt32(field collection.Field) translatedInt32 {
	return translatedInt32{
		TranslatedField{
			field: field,
		},
	}
}

func (t translatedInt32) GetObject() map[string]interface{} {
	key := t.field.Spec().Name
	return map[string]interface{}{
		key: Int32(0),
	}
}

func (t translatedInt32) GetArray() []interface{} {
	return Array()
}

// translation for int64 field type
func newTranslatedInt64(field collection.Field) translatedInt64 {
	return translatedInt64{
		TranslatedField{
			field: field,
		},
	}
}

func (t translatedInt64) GetObject() map[string]interface{} {
	key := t.field.Spec().Name
	return map[string]interface{}{
		key: Int64(0),
	}
}

func (t translatedInt64) GetArray() []interface{} {
	return Array()
}

// translation for double field type
func newTranslatedDouble(field collection.Field) translatedDouble {
	return translatedDouble{
		TranslatedField{
			field: field,
		},
	}
}

func (t translatedDouble) GetObject() map[string]interface{} {
	key := t.field.Spec().Name
	return map[string]interface{}{
		key: Float64(0),
	}
}

func (t translatedDouble) GetArray() []interface{} {
	return Array()
}

// translation for boolean field type
func newTranslatedBoolean(field collection.Field) translatedBoolean {
	return translatedBoolean{
		TranslatedField{
			field: field,
		},
	}
}

func (t translatedBoolean) GetObject() map[string]interface{} {
	key := t.field.Spec().Name
	return map[string]interface{}{
		key: Boolean(false),
	}
}

func (t translatedBoolean) GetArray() []interface{} {
	return Array()
}

// translation for array field type
func newTranslatedArray(field collection.Field) translatedArray {
	return translatedArray{
		TranslatedField{
			field: field,
		},
	}
}

func (t translatedArray) GetObject() map[string]interface{} {
	key := t.field.Spec().Name
	arrayFields := t.field.Spec().ArrayFields
	res := []interface{}{}
	for _, _field := range *arrayFields {
		currObj := GetTranslatedField(field.FromFieldSpec(&_field)).GetObject()
		item := currObj[_field.Name]
		res = append(res, item)
	}

	return map[string]interface{}{
		key: res,
	}
}

func (t translatedArray) GetArray() []interface{} {
	return Array()
}

// translation for object field type
func newTranslatedObject(field collection.Field) translatedObject {
	return translatedObject{
		TranslatedField{
			field: field,
		},
	}
}

func (t translatedObject) GetObject() map[string]interface{} {
	res := map[string]interface{}{}
	for _, _field := range *t.field.Spec().Object {
		currObj := GetTranslatedField(field.FromFieldSpec(&_field)).GetObject()
		res[_field.Name] = currObj[_field.Name]
	}

	key := t.field.Spec().Name

	return map[string]interface{}{
		key: res,
	}
}

func (t translatedObject) GetArray() []interface{} {
	return Array()
}

// translation for timestamp field type
func newTranslatedTimestamp(field collection.Field) translatedTimestamp {
	return translatedTimestamp{
		TranslatedField{
			field: field,
		},
	}
}

func (t translatedTimestamp) GetObject() map[string]interface{} {
	key := t.field.Spec().Name
	return map[string]interface{}{
		key: Time(time.Now()),
	}
}

func (t translatedTimestamp) GetArray() []interface{} {
	return Array()
}

// Geo JSON Section
func (t translatedGeoJSON) getCoordinateObject(key, _type string, child interface{}) map[string]interface{} {
	return map[string]interface{}{
		key: map[string]interface{}{
			"type":        _type,
			"coordinates": child,
		},
	}
}

// translation for geo json point field type
func newTranslatedGeoJSONPoint(field collection.Field) translatedGeoJSONPoint {
	return translatedGeoJSONPoint{
		translatedGeoJSON{
			TranslatedField{
				field: field,
			},
		},
	}
}

func (t translatedGeoJSONPoint) GetObject() map[string]interface{} {
	key := t.field.Spec().Name
	return t.getCoordinateObject(key, "Point", Array(Float64(0)))
}

func (t translatedGeoJSONPoint) GetArray() []interface{} {
	return []interface{}{}
}

// translation for geo json line string field type
func newTranslatedGeoJSONLineString(field collection.Field) translatedGeoJSONLineString {
	return translatedGeoJSONLineString{
		translatedGeoJSON{
			TranslatedField{
				field: field,
			},
		},
	}
}

func (t translatedGeoJSONLineString) GetObject() map[string]interface{} {
	key := t.field.Spec().Name
	return t.getCoordinateObject(key, "LineString", Array(Array(Float64(0))))
}

func (t translatedGeoJSONLineString) GetArray() []interface{} {
	return Array()
}

// translation for geo json polygon single ring
func newTranslatedGeoJSONPolygonSingleRing(field collection.Field) translatedGeoJSONPolygonSingleRing {
	return translatedGeoJSONPolygonSingleRing{
		translatedGeoJSON{
			TranslatedField{
				field: field,
			},
		},
	}
}

func (t translatedGeoJSONPolygonSingleRing) GetObject() map[string]interface{} {
	key := t.field.Spec().Name
	return t.getCoordinateObject(key, "Polygon", Array(Array(Array(Float64(0)))))
}

func (t translatedGeoJSONPolygonSingleRing) GetArray() []interface{} {
	return Array()
}

// translation for geo json polygon multiple ring
func newTranslatedGeoJSONPolygonMultipleRing(field collection.Field) translatedGeoJSONPolygonMultipleRing {
	return translatedGeoJSONPolygonMultipleRing{
		translatedGeoJSON{
			TranslatedField{
				field: field,
			},
		},
	}
}

func (t translatedGeoJSONPolygonMultipleRing) GetObject() map[string]interface{} {
	key := t.field.Spec().Name
	return t.getCoordinateObject(key, "Polygon", Array(Array(Array(Float64(0))), Array(Array(Float64(0)))))
}

func (t translatedGeoJSONPolygonMultipleRing) GetArray() []interface{} {
	return Array()
}

// translation for geo json multi point
func newTranslatedGeoJSONMultiPoint(field collection.Field) translatedGeoJSONMultiPoint {
	return translatedGeoJSONMultiPoint{
		translatedGeoJSON{
			TranslatedField{
				field: field,
			},
		},
	}
}

func (t translatedGeoJSONMultiPoint) GetObject() map[string]interface{} {
	key := t.field.Spec().Name
	return t.getCoordinateObject(key, "MultiPoint", Array(Array(Float64(0), Float64(0))))
}

func (t translatedGeoJSONMultiPoint) GetArray() []interface{} {
	return Array()
}

// translation for geo json multi line string
func newTranslatedGeoJSONMultiLineString(field collection.Field) translatedGeoJSONMultiLineString {
	return translatedGeoJSONMultiLineString{
		translatedGeoJSON{
			TranslatedField{
				field: field,
			},
		},
	}
}

func (t translatedGeoJSONMultiLineString) GetObject() map[string]interface{} {
	key := t.field.Spec().Name
	return t.getCoordinateObject(key, "MultiLineString", Array(Array(Array(Float64(0), Float64(0)))))
}

func (t translatedGeoJSONMultiLineString) GetArray() []interface{} {
	return Array()
}

// translation for geo json multi polygon
func newTranslatedGeoJSONMultiPolygon(field collection.Field) translatedGeoJSONMultiPolygon {
	return translatedGeoJSONMultiPolygon{
		translatedGeoJSON{
			TranslatedField{
				field: field,
			},
		},
	}
}

func (t translatedGeoJSONMultiPolygon) GetObject() map[string]interface{} {
	key := t.field.Spec().Name
	return t.getCoordinateObject(key, "MultiPolygon", Array(
		Array(Array(Array(Float64(0), Float64(0)))),
		Array(Array(Array(Float64(0), Float64(0)))),
	))
}

func (t translatedGeoJSONMultiPolygon) GetArray() []interface{} {
	return Array()
}

// translation for geo json geometry collection
func newTranslatedGeoJSONGeometryCollection(field collection.Field) translatedGeoJSONGeometryCollection {
	return translatedGeoJSONGeometryCollection{
		translatedGeoJSON{
			TranslatedField{
				field: field,
			},
		},
	}
}

func (t translatedGeoJSONGeometryCollection) GetObject() map[string]interface{} {
	key := t.field.Spec().Name
	arrayFields := t.field.Spec().ArrayFields
	collections := Array()
	for _, _field := range *arrayFields {
		item := GetTranslatedField(field.FromFieldSpec(&_field)).GetObject()
		collections = append(collections, item)
	}

	return map[string]interface{}{
		key: map[string]interface{}{
			"type":       "GeometryCollection",
			"geometries": collections,
		},
	}
}

// translation for legacy coordinate field
func newTranslatedLegacyCoordinate(field collection.Field) translatedLegacyCoordinate {
	return translatedLegacyCoordinate{
		TranslatedField{
			field: field,
		},
	}
}

func (t translatedLegacyCoordinate) GetObject() map[string]interface{} {
	return map[string]interface{}{}
}

func (t translatedLegacyCoordinate) GetArray() []interface{} {
	return Array()
}

// map field to correct translated field
func GetTranslatedField(_field collection.Field) TranslatedFieldIf {
	switch _field.Spec().Type {
	case field.TypeString:
		return newTranslatedString(_field)
	case field.TypeInt32:
		return newTranslatedInt32(_field)
	case field.TypeInt64:
		return newTranslatedInt64(_field)
	case field.TypeDouble:
		return newTranslatedDouble(_field)
	case field.TypeBoolean:
		return newTranslatedBoolean(_field)
	case field.TypeArray:
		return newTranslatedArray(_field)
	case field.TypeObject:
		return newTranslatedObject(_field)
	case field.TypeTimestamp:
		return newTranslatedTimestamp(_field)
	case field.TypeGeoJSONPoint:
		return newTranslatedGeoJSONPoint(_field)
	case field.TypeGeoJSONLineString:
		return newTranslatedGeoJSONLineString(_field)
	case field.TypeGeoJSONPolygonSingleRing:
		return newTranslatedGeoJSONPolygonSingleRing(_field)
	case field.TypeGeoJSONPolygonMultipleRing:
		return newTranslatedGeoJSONPolygonMultipleRing(_field)
	case field.TypeGeoJSONMultiPoint:
		return newTranslatedGeoJSONMultiPoint(_field)
	case field.TypeGeoJSONMultiLineString:
		return newTranslatedGeoJSONMultiLineString(_field)
	case field.TypeGeoJSONMultiPolygon:
		return newTranslatedGeoJSONMultiPolygon(_field)
	case field.TypeGeoJSONGeometryCollection:
		return newTranslatedGeoJSONGeometryCollection(_field)
	}

	return TranslatedField{}
}
