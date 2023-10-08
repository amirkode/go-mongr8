/*
Copyright (c) 2023 the go-mongr8 Authors and Contributors
[@see Authors file]

Licensed under the MIT License
(https://opensource.org/licenses/MIT)
*/
package index

type (
	IndexType string
)

const (
	TypeSingleField       IndexType = "single_field"
	TypeCompound          IndexType = "compound"
	TypeText              IndexType = "text"
	TypeGeopatial2dsphere IndexType = "geospatial_2dsphere"
	TypeUnique            IndexType = "unique"
	TypePartial           IndexType = "partial"
	TypeCollation         IndexType = "collation"
	// raw user defined index
	TypeRaw               IndexType = "raw"
)
