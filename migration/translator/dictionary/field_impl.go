package dictionary

import (
	"time"

	"github.com/amirkode/go-mongr8/collection"
	"github.com/amirkode/go-mongr8/collection/field"
	"go.mongodb.org/mongo-driver/bson"
)

type (
	TranslatedField struct {
		TranslatedFieldIf
		field   collection.Field
		isArray bool
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

	translatedGeoJSONGeoCollection struct {
		translatedGeoJSON
	}

	translatedLegacyCoordinate struct {
		TranslatedField
	}
)

func (t TranslatedField) getItem() interface{} {
	if t.isArray {
		return t.getArray()
	}

	return t.getObject()
}

// translation for string field type
func newTranslatedString(field collection.Field) translatedString {
	return translatedString{
		TranslatedField: TranslatedField{
			field: field,
		},
	}
}

func (t translatedString) getObject() bson.M {
	key := t.field.Spec().Name
	return bson.M{
		key: "",
	}
}

func (t translatedString) getArray() bson.A {
	return bson.A{}
}

// translation for int32 field type
func newTranslatedInt32(field collection.Field) translatedInt32 {
	return translatedInt32{
		TranslatedField{
			field: field,
		},
	}
}

func (t translatedInt32) getObject() bson.M {
	key := t.field.Spec().Name
	return bson.M{
		key: int32(0),
	}
}

func (t translatedInt32) getArray() bson.A {
	return bson.A{}
}

// translation for int64 field type
func newTranslatedInt64(field collection.Field) translatedInt64 {
	return translatedInt64{
		TranslatedField{
			field: field,
		},
	}
}

func (t translatedInt64) getObject() bson.M {
	key := t.field.Spec().Name
	return bson.M{
		key: int64(0),
	}
}

func (t translatedInt64) getArray() bson.A {
	return bson.A{}
}

// translation for double field type
func newTranslatedDouble(field collection.Field) translatedDouble {
	return translatedDouble{
		TranslatedField{
			field: field,
		},
	}
}

func (t translatedDouble) getObject() bson.M {
	key := t.field.Spec().Name
	return bson.M{
		key: float64(0),
	}
}

func (t translatedDouble) getArray() bson.A {
	return bson.A{}
}

// translation for boolean field type
func newTranslatedBoolean(field collection.Field) translatedBoolean {
	return translatedBoolean{
		TranslatedField{
			field: field,
		},
	}
}

func (t translatedBoolean) getObject() bson.M {
	key := t.field.Spec().Name
	return bson.M{
		key: false,
	}
}

func (t translatedBoolean) getArray() bson.A {
	return bson.A{}
}

// translation for array field type
func newTranslatedArray(field collection.Field) translatedArray {
	return translatedArray{
		TranslatedField{
			field: field,
			isArray: true,
		},
	}
}

func (t translatedArray) getOject() bson.M {
	return bson.M{}
}

func (t translatedArray) getArray() bson.A {
	arrayFields := t.field.Spec().ArrayFields
	res := bson.A{}
	for _, _field := range *arrayFields {
		item := getTranslatedField(field.FromFieldSpec(&_field)).getItem()
		res = append(res, item)
	}

	return res
}

// translation for object field type
func newTranslatedObject(field collection.Field) translatedObject {
	return translatedObject{
		TranslatedField{
			field: field,
		},
	}
}

func (t translatedObject) getObject() bson.M {
	res := bson.M{}
	for _, _field := range *t.field.Spec().Object {
		res[_field.Name] = getTranslatedField(field.FromFieldSpec(&_field)).getItem()
	}

	return res
}

func (t translatedObject) getArray() bson.A {
	return bson.A{}
}

// translation for timestamp field type
func newTranslatedTimestamp(field collection.Field) translatedTimestamp {
	return translatedTimestamp{
		TranslatedField{
			field: field,
		},
	}
}

func (t translatedTimestamp) getObject() bson.M {
	key := t.field.Spec().Name
	return bson.M{
		key: time.Now(),
	}
}

func (t translatedTimestamp) getArray() bson.A {
	return bson.A{}
}

// Geo JSON Section
func (t translatedGeoJSON) getCoordinateObject(key, _type string, child interface {}) bson.M {
	return bson.M{
		key: bson.M{
			"type": _type,
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

func (t translatedGeoJSONPoint) getObject() bson.M {
	key := t.field.Spec().Name
	return t.getCoordinateObject(key, "Point", bson.A{float64(0)})
}

func (t translatedGeoJSONPoint) getArray() bson.A {
	return bson.A{}
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

func (t translatedGeoJSONLineString) getObject() bson.M {
	key := t.field.Spec().Name
	return t.getCoordinateObject(key, "LineString", bson.A{bson.A{float64(0)}})
}

// TODO: continue

// map field to correct translated field
func  getTranslatedField(_field collection.Field) TranslatedFieldIf {
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
	}

	return TranslatedField{}
}