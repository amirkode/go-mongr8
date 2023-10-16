/*
Copyright (c) 2023 the go-mongr8 Authors and Contributors
[@see Authors file]

Licensed under the MIT License
(https://opensource.org/licenses/MIT)
*/
package field

import (
	"internal/util"
)

type (
	FieldType string
)

const (
	TypeString                      FieldType = "string"
	TypeInt32                       FieldType = "int32"
	TypeInt64                       FieldType = "int64"
	TypeDouble                      FieldType = "double64"
	TypeBoolean                     FieldType = "boolean"
	TypeArray                       FieldType = "array"
	TypeObject                      FieldType = "object"
	TypeTimestamp                   FieldType = "timestamp"
	TypeGeoJSONPoint                FieldType = "geo_json_point"
	TypeGeoJSONLineString           FieldType = "geo_json_line_string"
	TypeGeoJSONPolygonSingleRing    FieldType = "geo_json_polygon_single_ring"
	TypeGeoJSONPolygonMultipleRing  FieldType = "geo_json_polygon_multiple_ring"
	TypeGeoJSONMultiPoint           FieldType = "geo_json_multi_point"
	TypeGeoJSONMultiLineString      FieldType = "geo_json_multi_line_string"
	TypeGeoJSONMultiPolygon         FieldType = "geo_json_multi_polygon"
	TypeGeoJSONGeometryCollection   FieldType = "geo_json_geometry_collection"
	TypeLegacyCoordinateArray       FieldType = "legacy_coordinate_array"
	TypeLegacyCoordinateEmbeddedDoc FieldType = "legacy_coordinate_embedded_doc"
	// other types such as Decimal128, etc.
	// might be added in the future update
)

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
