/*
Copyright (c) 2023-present the go-mongr8 Authors and Contributors
[@see Authors file]

Licensed under the MIT License
(https://opensource.org/licenses/MIT)
*/
package field

import (
	"github.com/amirkode/go-mongr8/internal/util"
)

type (
	FieldType string
	FieldExtra string
)

// make sure constant name is exactly same as it's value
const (
	TypeString                      FieldType = "TypeString"
	TypeInt32                       FieldType = "TypeInt32"
	TypeInt64                       FieldType = "TypeInt64"
	TypeDouble                      FieldType = "TypeDouble"
	TypeBoolean                     FieldType = "TypeBoolean"
	TypeArray                       FieldType = "TypeArray"
	TypeObject                      FieldType = "TypeObject"
	TypeTimestamp                   FieldType = "TypeTimestamp"
	TypeGeoJSONPoint                FieldType = "TypeGeoJSONPoint"
	TypeGeoJSONLineString           FieldType = "TypeGeoJSONLineString"
	TypeGeoJSONPolygonSingleRing    FieldType = "TypeGeoJSONPolygonSingleRing"
	TypeGeoJSONPolygonMultipleRing  FieldType = "TypeGeoJSONPolygonMultipleRing"
	TypeGeoJSONMultiPoint           FieldType = "TypeGeoJSONMultiPoint"
	TypeGeoJSONMultiLineString      FieldType = "TypeGeoJSONMultiLineString"
	TypeGeoJSONMultiPolygon         FieldType = "TypeGeoJSONMultiPolygon"
	TypeGeoJSONGeometryCollection   FieldType = "TypeGeoJSONGeometryCollection"
	TypeLegacyCoordinateArray       FieldType = "TypeLegacyCoordinateArray"
	TypeLegacyCoordinateEmbeddedDoc FieldType = "TypeLegacyCoordinateEmbeddedDoc"
	// other types such as Decimal128, etc.
	// might be added in the future update

	// extra keys
	ExtraDrop FieldExtra = "drop"
)

func GetTypePointer(fieldType FieldType) *FieldType {
	return &fieldType
}

func (f FieldType) ToString() string {
	return string(f)
}

// implement /internal/util/Comparable
func (f FieldType) CompareTo(other FieldType) int {
	if string(f) > string(other) {
		return 1
	}

	if string(f) < string(other) {
		return -1
	}

	return 0
}

func (f FieldType) IsNumeric() bool {
	return util.InList(f, []FieldType{
		TypeInt32,
		TypeInt64,
		TypeDouble,
	})
}
