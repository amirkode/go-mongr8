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

func (t TranslatedField) getItem() interface{} {
	key := t.field.Spec().Name
	return t.getObject()[key]
}

// translation for string field type
func newTranslatedString(field collection.Field) translatedString {
	return translatedString{
		TranslatedField: TranslatedField{
			field: field,
		},
	}
}

func (t translatedString) getObject() map[string]interface{} {
	key := t.field.Spec().Name
	return map[string]interface{}{
		key: String(""),
	}
}

func (t translatedString) getArray() []interface{} {
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

func (t translatedInt32) getObject() map[string]interface{} {
	key := t.field.Spec().Name
	return map[string]interface{}{
		key: Int32(0),
	}
}

func (t translatedInt32) getArray() []interface{} {
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

func (t translatedInt64) getObject() map[string]interface{} {
	key := t.field.Spec().Name
	return map[string]interface{}{
		key: Int64(0),
	}
}

func (t translatedInt64) getArray() []interface{} {
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

func (t translatedDouble) getObject() map[string]interface{} {
	key := t.field.Spec().Name
	return map[string]interface{}{
		key: Float64(0),
	}
}

func (t translatedDouble) getArray() []interface{} {
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

func (t translatedBoolean) getObject() map[string]interface{} {
	key := t.field.Spec().Name
	return map[string]interface{}{
		key: Boolean(false),
	}
}

func (t translatedBoolean) getArray() []interface{} {
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

func (t translatedArray) getOject() map[string]interface{} {
	key := t.field.Spec().Name
	arrayFields := t.field.Spec().ArrayFields
	res := []interface{}{}
	for _, _field := range *arrayFields {
		item := getTranslatedField(field.FromFieldSpec(&_field)).getItem()
		res = append(res, item)
	}

	return map[string]interface{}{
		key: res,
	}
}

func (t translatedArray) getArray() []interface{} {
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

func (t translatedObject) getObject() map[string]interface{} {
	res := map[string]interface{}{}
	for _, _field := range *t.field.Spec().Object {
		res[_field.Name] = getTranslatedField(field.FromFieldSpec(&_field)).getItem()
	}

	return res
}

func (t translatedObject) getArray() []interface{} {
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

func (t translatedTimestamp) getObject() map[string]interface{} {
	key := t.field.Spec().Name
	return map[string]interface{}{
		key: Time(time.Now()),
	}
}

func (t translatedTimestamp) getArray() []interface{} {
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

func (t translatedGeoJSONPoint) getObject() map[string]interface{} {
	key := t.field.Spec().Name
	return t.getCoordinateObject(key, "Point", Array(Float64(0)))
}

func (t translatedGeoJSONPoint) getArray() []interface{} {
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

func (t translatedGeoJSONLineString) getObject() map[string]interface{} {
	key := t.field.Spec().Name
	return t.getCoordinateObject(key, "LineString", Array(Array(Float64(0))))
}

func (t translatedGeoJSONLineString) getArray() []interface{} {
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

func (t translatedGeoJSONPolygonSingleRing) getObject() map[string]interface{} {
	key := t.field.Spec().Name
	return t.getCoordinateObject(key, "Polygon", Array(Array(Array(Float64(0)))))
}

func (t translatedGeoJSONPolygonSingleRing) getArray() []interface{} {
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

func (t translatedGeoJSONPolygonMultipleRing) getObject() map[string]interface{} {
	key := t.field.Spec().Name
	return t.getCoordinateObject(key, "Polygon", Array(Array(Array(Float64(0))), Array(Array(Float64(0)))))
}

func (t translatedGeoJSONPolygonMultipleRing) getArray() []interface{} {
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

func (t translatedGeoJSONMultiPoint) getObject() map[string]interface{} {
	key := t.field.Spec().Name
	return t.getCoordinateObject(key, "MultiPoint", Array(Array(Float64(0), Float64(0))))
}

func (t translatedGeoJSONMultiPoint) getArray() []interface{} {
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

func (t translatedGeoJSONMultiLineString) getObejct() map[string]interface{} {
	key := t.field.Spec().Name
	return t.getCoordinateObject(key, "MultiLineString", Array(Array(Array(Float64(0), Float64(0)))))
}

func (t translatedGeoJSONMultiLineString) getArray() []interface{} {
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

func (t translatedGeoJSONMultiPolygon) getObject() map[string]interface{} {
	key := t.field.Spec().Name
	return t.getCoordinateObject(key, "MultiPolygon", Array(
		Array(Array(Array(Float64(0), Float64(0)))),
		Array(Array(Array(Float64(0), Float64(0)))),
	))
}

func (t translatedGeoJSONMultiPolygon) getArray() []interface{} {
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

func (t translatedGeoJSONGeometryCollection) getObject() map[string]interface{} {
	key := t.field.Spec().Name
	arrayFields := t.field.Spec().ArrayFields
	collections := Array()
	for _, _field := range *arrayFields {
		item := getTranslatedField(field.FromFieldSpec(&_field)).getObject()
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

func (t translatedLegacyCoordinate) getObject() map[string]interface{} {
	return map[string]interface{}{}
}

func (t translatedLegacyCoordinate) getArray() []interface{} {
	return Array()
}

// map field to correct translated field
func getTranslatedField(_field collection.Field) TranslatedFieldIf {
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
