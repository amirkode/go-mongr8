/*
Copyright (c) 2023-present the go-mongr8 Authors and Contributors
[@see Authors file]

Licensed under the MIT License
(https://opensource.org/licenses/MIT)
*/
package field

import "fmt"

// TODO: add default value
type Spec struct {
	Name string
	// Type of the field
	Type FieldType

	// Array items, if current type is an array
	// this can be any field type
	ArrayFields *[]Spec

	// Children of object, if current type is an object
	Object *[]Spec

	// Nullable flag
	Nullable bool

	// Index required if schema is sorted
	// for adding/setting new field
	index int
}

type FieldSpec struct {
	spec *Spec
}

func (b *FieldSpec) Spec() *Spec {
	return b.spec
}

func (b *FieldSpec) addArrayField(fields ...*FieldSpec) *FieldSpec {
	if b.spec.ArrayFields == nil {
		// init slice of array fields
		arrFields := []Spec{}
		b.spec.ArrayFields = &arrFields
	}

	arr := *b.spec.ArrayFields
	for _, f := range fields {
		arr = append(arr, *f.Spec())
	}

	b.spec.ArrayFields = &arr

	return b
}

func (b *FieldSpec) addObjectFields(fields ...*FieldSpec) *FieldSpec {
	if b.spec.Object == nil {
		objects := []Spec{}
		b.spec.Object = &objects
	}

	objects := *b.spec.Object
	for _, o := range fields {
		objects = append(objects, *o.Spec())
	}

	b.spec.Object = &objects

	return b
}

func (b *FieldSpec) ObjectHasKey(key string) bool {
	if b.spec.Object != nil {
		for _, field := range *b.spec.Object {
			if field.Name == key {
				return true
			}
		}
	}

	return false
}

func (b *FieldSpec) SetNullable() *FieldSpec {
	b.spec.Nullable = true

	return b
}

func baseField(name string, fieldType FieldType) *FieldSpec {
	// already validated in translation level
	// if len(name) > 128 {
	// 	panic("Cannot have a field name more than 128 characters len")
	// }

	field := &FieldSpec{
		&Spec{
			Name:     name,
			Type:     fieldType,
			Nullable: false,
		},
	}

	return field
}

func FromFieldSpec(spec *Spec) *FieldSpec {
	return &FieldSpec{
		spec,
	}
}

func StringField(name string) *FieldSpec {
	return baseField(name, TypeString)
}

func Int32Field(name string) *FieldSpec {
	return baseField(name, TypeInt32)
}

func Int64Field(name string) *FieldSpec {
	return baseField(name, TypeInt64)
}

func DoubleField(name string) *FieldSpec {
	return baseField(name, TypeDouble)
}

func BooleanField(name string) *FieldSpec {
	return baseField(name, TypeBoolean)
}

func ArrayField(name string, fields ...*FieldSpec) *FieldSpec {
	// TODO: add support for array of different types, when the usecase is clear
	// for now, we only support for a single type, but we still support empty fields
	// since, there's a case that we field as a type representative only
	// yet will eventually be validated again at the translation level
	if len(fields) > 1 {
		panic(fmt.Sprintf("ArrayField needs 1 field must be declared"))
	}

	field := baseField(name, TypeArray)
	field.addArrayField(fields...)

	return field
}

func ObjectField(name string, fields ...*FieldSpec) *FieldSpec {
	field := baseField(name, TypeObject)
	field.addObjectFields(fields...)

	return field
}

func TimestampField(name string) *FieldSpec {
	return baseField(name, TypeTimestamp)
}

func GeoJSONPointField(name string) *FieldSpec {
	return baseField(name, TypeGeoJSONPoint)
}

func GeoJSONLineStringField(name string) *FieldSpec {
	return baseField(name, TypeGeoJSONLineString)
}

func GeoJSONPolygonSingleRingField(name string) *FieldSpec {
	return baseField(name, TypeGeoJSONPolygonSingleRing)
}

func GeoJSONPolygonMultipleRingField(name string) *FieldSpec {
	return baseField(name, TypeGeoJSONPolygonMultipleRing)
}

func GeoJSONMultiPointField(name string) *FieldSpec {
	return baseField(name, TypeGeoJSONMultiPoint)
}

func GeoJSONMultiLineStringField(name string) *FieldSpec {
	return baseField(name, TypeGeoJSONMultiLineString)
}

func GeoJSONMultiPolygonField(name string) *FieldSpec {
	return baseField(name, TypeGeoJSONMultiPolygon)
}

func GeoJSONGeometryCollectionField(name string) *FieldSpec {
	return baseField(name, TypeGeoJSONGeometryCollection)
}

func LegacyCoordinateArrayField(name string) *FieldSpec {
	return baseField(name, TypeLegacyCoordinateArray)
}

// add additional functions for legacy coordinate embedded document
type legacyCoordinateEmbeddedDocSpec struct {
	FieldSpec
	xIsSet bool
	yIsSet bool
}

func (s *legacyCoordinateEmbeddedDocSpec) setCoordinateField(name string, isX bool) {
	// check whether a coordinate key is already set
	if (isX && s.xIsSet) || (!isX && s.yIsSet) {
		key := "x"
		if !isX {
			key = "y"
		}

		panic(fmt.Sprintf("%s field is already set on Legacy Coordinate: %s", key, s.Spec().Name))
	}

	if s.ObjectHasKey(name) {
		panic(fmt.Sprintf("Key %s already exists in %s object", name, s.Spec().Name))
	}

	// update state
	if isX {
		s.xIsSet = true
	} else {
		s.yIsSet = true
	}

	s.addObjectFields(DoubleField(name))
}

func (s *legacyCoordinateEmbeddedDocSpec) SetCoordinateX(name string) *legacyCoordinateEmbeddedDocSpec {
	s.setCoordinateField(name, true)

	return s
}

func (s *legacyCoordinateEmbeddedDocSpec) SetCoordinateY(name string) *legacyCoordinateEmbeddedDocSpec {
	s.setCoordinateField(name, false)

	return s
}

func LegacyCoordinateEmbeddedDocField(name string) *legacyCoordinateEmbeddedDocSpec {
	baseField := baseField(name, TypeLegacyCoordinateEmbeddedDoc)
	return &legacyCoordinateEmbeddedDocSpec{
		FieldSpec: *baseField,
		xIsSet:    false,
		yIsSet:    false,
	}
}
