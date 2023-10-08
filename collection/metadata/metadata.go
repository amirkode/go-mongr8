/*
Copyright (c) 2023 the go-mongr8 Authors and Contributors
[@see Authors file]

Licensed under the MIT License
(https://opensource.org/licenses/MIT)
*/
package metadata

import "fmt"

type Spec struct {
	Name    string
	Options *map[string]interface{}
	Type CollectionType
}

type MetadataSpec struct {
	spec *Spec
}

func (s *MetadataSpec) Spec() *Spec {
	return s.spec
}

func (s *MetadataSpec) Capped(size int) *MetadataSpec {
	if s.Spec().Options == nil {
		s.Spec().Options = &map[string]interface{}{}
	} else {
		panic(fmt.Sprintf("Cannot add capped option, another option already exists on collection: %s", s.Spec().Name))
	}

	(*s.Spec().Options)["capped"] = true
	(*s.Spec().Options)["size"] = true

	// set collection type to capped collection
	s.Spec().Type = TypeCappedCollection

	return s
}

func (s *MetadataSpec) TTL(expiredAfter int) *MetadataSpec {
	if s.Spec().Options == nil {
		s.Spec().Options = &map[string]interface{}{}
	} else {
		panic(fmt.Sprintf("Cannot add TTL option, another option already exists on collection: %s", s.Spec().Name))
	}

	(*s.Spec().Options)["expiredAfterSeconds"] = expiredAfter

	// set collection type to ttl collection
	s.Spec().Type = TypeTTLCollection

	return s
}

func (s *MetadataSpec) AsView() *MetadataSpec {
	// set collection type to view
	s.Spec().Type = TypeViewCollection

	return s
}

func InitMetadata(name string) *MetadataSpec {
	res := &MetadataSpec{
		&Spec{
			Name: name,
			Type: TypeDefaultCollection,
		},
	}
	return res
}
