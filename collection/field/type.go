/*
Copyright (c) 2023 the go-mongr8 Authors and Contributors
[@see Authors file]

Licensed under the MIT License
(https://opensource.org/licenses/MIT)
*/
package field

type (
	FieldType string
)

const (
	TypeString           FieldType = "string"
	TypeInt32            FieldType = "int32"
	TypeInt64            FieldType = "int64"
	TypeDouble           FieldType = "double64"
	TypeBoolean          FieldType = "boolean"
	TypeArray            FieldType = "array"
	TypeObject           FieldType = "object"
	TypeTimestamp        FieldType = "timestamp"
	TypeGeoJSON          FieldType = "geo_json"
	TypeLegacyCoordinate FieldType = "legacy_coordinate"
	// other types such as Decimal128, etc.
	// might be added in the future update
)
