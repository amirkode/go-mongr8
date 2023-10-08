/*
Copyright (c) 2023 the go-mongr8 Authors and Contributors
[@see Authors file]

Licensed under the MIT License
(https://opensource.org/licenses/MIT)
*/
package field

import "fmt"

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
}

type BaseSpec struct {
	spec *Spec
}

func (b *BaseSpec) Spec() *Spec {
	return b.spec
}

func (b *BaseSpec) AddArrayField(s *BaseSpec) *BaseSpec {
	if b.spec.ArrayFields == nil {
		// init slice of array fields
		arrFields := []Spec{}
		b.spec.ArrayFields = &arrFields
	}

	currArr := *b.spec.ArrayFields
	currArr = append(currArr, *s.Spec())

	b.spec.ArrayFields = &currArr

	return b
}

func (b *BaseSpec) AddObjectFields(fields *[]BaseSpec) *BaseSpec {
	if b.spec.Object == nil {
		objects := []Spec{}
		b.spec.Object = &objects
	}

	objects := *b.spec.Object
	for _, o := range *fields {
		objects = append(objects, *o.Spec())
	}

	b.spec.Object = &objects

	return b
}

func (b *BaseSpec) ObjectHasKey(key string) bool {
	if b.spec.Object != nil {
		for _, field := range *b.spec.Object {
			if field.Name == key {
				return true
			}
		}
	}

	return false
}

func (b *BaseSpec) SetNullable() *BaseSpec {
	b.spec.Nullable = true

	return b
}

func baseField(name string, fieldType FieldType) *BaseSpec {
	field := &BaseSpec{
		&Spec{
			Name:     name,
			Type:     fieldType,
			Nullable: false,
		},
	}

	return field
}

func FromFieldSpec(spec *Spec) *BaseSpec {
	return &BaseSpec{
		spec,
	}
}

func StringField(name string) *BaseSpec {
	return baseField(name, TypeString)
}

func Int64Field(name string) *BaseSpec {
	return baseField(name, TypeInt64)
}

func Int32Field(name string) *BaseSpec {
	return baseField(name, TypeInt32)
}

func DoubleField(name string) *BaseSpec {
	return baseField(name, TypeDouble)
}

func ArrayField(name string) *BaseSpec {
	return baseField(name, TypeArray)
}

func ObjectField(name string) *BaseSpec {
	return baseField(name, TypeObject)
}

func TimestampField(name string) *BaseSpec {
	return baseField(name, TypeTimestamp)
}

func GeoJSONPointField(name string) *BaseSpec {
	return baseField(name, TypeGeoJSONPoint)
}

func GeoJSONLineStringField(name string) *BaseSpec {
	return baseField(name, TypeGeoJSONLineString)
}

func GeoJSONPolygonSingleRingField(name string) *BaseSpec {
	return baseField(name, TypeGeoJSONPolygonSingleRing)
}

func GeoJSONPolygonMultipleRingField(name string) *BaseSpec {
	return baseField(name, TypeGeoJSONPolygonMultipleRing)
}

func GeoJSONMultiPointField(name string) *BaseSpec {
	return baseField(name, TypeGeoJSONMultiPoint)
}

func GeoJSONMultiLineStringField(name string) *BaseSpec {
	return baseField(name, TypeGeoJSONMultiLineString)
}

func GeoJSONMultiPolygonField(name string) *BaseSpec {
	return baseField(name, TypeGeoJSONMultiPolygon)
}

func GeoJSONGeometryCollectionField(name string) *BaseSpec {
	return baseField(name, TypeGeoJSONGeometryCollection)
}

func LegacyCoordinateArrayField(name string) *BaseSpec {
	return baseField(name, TypeLegacyCoordinateArray)
}

// add additional functions for legacy coordinate embedded document
type legacyCoordinateEmbeddedDocSpec struct {
	BaseSpec
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

	s.AddObjectFields(&[]BaseSpec{
		*DoubleField(name),
	})
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
		BaseSpec: *baseField,
		xIsSet: false,
		yIsSet: false,
	}
}
