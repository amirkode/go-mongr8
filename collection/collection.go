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
