/*
Create date: {{ .CreateDate}}
Created by: go-mongr8

Copyright (c) 2023 the go-mongr8 Authors and Contributors
[@see Authors file]

Licensed under the MIT License
(https://opensource.org/licenses/MIT)
*/

package collection

import (
	"github.com/amirkode/go-mongr8/collection"
	"github.com/amirkode/go-mongr8/collection/metadata"
)

type {{ .Entity}} struct {
	collection.Collection
}

func ({{ .Entity}}) Metadata() collection.Metadata {
	return metadata.InitMetadata("{{ .Collection}}")
}

func ({{ .Entity}}) Fields() []collection.Field {
	return []collection.Field{
		// define fields here
		// examples:
		// import: "github.com/amirkode/go-mongr8/collection/field"
		// field.Int32Field("id"),
		// field.ArrayField("names").SetArrayField(field.StringField("")),
	}
}

func ({{ .Entity}}) Indexes() []collection.Index {
	return []collection.Index{
		// define indexes here
		// examples:
		// import: "github.com/amirkode/go-mongr8/collection/index"
	}
}
