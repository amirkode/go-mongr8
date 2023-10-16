/*
Copyright (c) 2023 the go-mongr8 Authors and Contributors
[@see Authors file]

Licensed under the MIT License
(https://opensource.org/licenses/MIT)
*/
package collection

import (
	"github.com/amirkode/go-mongr8/collection/metadata"
	"github.com/amirkode/go-mongr8/collection/field"
	"github.com/amirkode/go-mongr8/collection/index"
)

type (
	// stores basic informations of a collection
	// this also holds type of usage: collection or view
	Metadata interface {
		Spec() *metadata.Spec
	}

	Field interface {
		Spec() *field.Spec
	}

	Index interface {
		Spec() *index.Spec
	}

	// Collection entity manages a single MongoDB collection
	Collection interface {
		Collection()  Metadata
		Fields()      []Field
		Indexes()     []Index
	}
)

// // LookupField finds a field by key
// func (coll Collection) LookupField(key string) *field.Spec {
// 	fields := coll.Fields()
// 	for _, field := range fields {
// 		if field.Spec().Name == key {
// 			return field.Spec()
// 		}
// 	}
	
// 	return nil
// }

func FieldFromType(name string, _type field.FieldType) Field {
	switch _type {
	case field.TypeString:
		return field.StringField(name)
	case field.TypeInt32:
		return field.Int32Field(name)
	case field.TypeInt64:
		return field.Int64Field(name)
	case field.TypeDouble:
		return field.DoubleField(name)
	case field.TypeBoolean:
		return field.BooleanField(name)
	case field.TypeArray:
		return field.ArrayField(name)
	case field.TypeObject:
		return field.ObjectField(name)
	case field.TypeTimestamp:
		return field.TimestampField(name)
	case field.TypeGeoJSONPoint:
		return field.GeoJSONPointField(name)
	case field.TypeGeoJSONLineString:
		return field.GeoJSONLineStringField(name)
	case field.TypeGeoJSONPolygonSingleRing:
		return field.GeoJSONPolygonSingleRingField(name)
	case field.TypeGeoJSONPolygonMultipleRing:
		return field.GeoJSONPolygonMultipleRingField(name)
	case field.TypeGeoJSONMultiPoint:
		return field.GeoJSONMultiPointField(name)
	case field.TypeGeoJSONMultiLineString:
		return field.GeoJSONMultiLineStringField(name)
	case field.TypeGeoJSONMultiPolygon:
		return field.GeoJSONMultiPolygonField(name)
	case field.TypeGeoJSONGeometryCollection:
		return field.GeoJSONGeometryCollectionField(name)
	}

	return &field.FieldSpec{}
}