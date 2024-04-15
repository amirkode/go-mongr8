/*
Copyright (c) 2023-present the go-mongr8 Authors and Contributors
[@see Authors file]

Licensed under the MIT License
(https://opensource.org/licenses/MIT)
*/
package index

type (
	IndexType string
)

const (
	// Index Types
	TypeSingleField       IndexType = "TypeSingleField"
	TypeCompound          IndexType = "TypeCompound"
	TypeText              IndexType = "TypeText"
	TypeGeopatial2dsphere IndexType = "TypeGeopatial2dsphere"
	TypeHashed            IndexType = "TypeHashedIndex"
	TypeRaw               IndexType = "TypeRaw"
	// Index Options
	OptionSparse           = "sparse"
	OptionBackground       = "background"
	OptionUnique           = "unique"
	OptionHidden           = "hidden"
	OptionPartialFilterExp = "partialFilterExpression"
	OptionTTL              = "expireAfterSeconds"
	OptionCollation        = "collation"
)
