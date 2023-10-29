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
	TypeSingleField       IndexType = "TypeSingleField"
	TypeCompound          IndexType = "TypeCompound"
	TypeText              IndexType = "TypeText"
	TypeGeopatial2dsphere IndexType = "TypeGeopatial2dsphere"
	TypeUnique            IndexType = "TypeUnique"
	TypePartial           IndexType = "TypePartial"
	TypeCollation         IndexType = "TypeCollation"
	// raw user defined index
	TypeRaw               IndexType = "TypeRaw"
)
