/*
Copyright (c) 2023-present the go-mongr8 Authors and Contributors
[@see Authors file]

Licensed under the MIT License
(https://opensource.org/licenses/MIT)
*/
package metadata

import "fmt"

type Spec struct {
	Name    string
	Options *map[CollectionOption]interface{}
	Type    CollectionType
}

type MetadataSpec struct {
	spec *Spec
}

func (s *MetadataSpec) Spec() *Spec {
	return s.spec
}

func (s *MetadataSpec) initOptions() {
	if s.Spec().Options == nil {
		s.Spec().Options = &map[CollectionOption]interface{}{}
	}
}

func (s *MetadataSpec) Capped(size int64) *MetadataSpec {
	s.initOptions()

	_, found := (*s.Spec().Options)[CollectionOptionCapped]
	if found {
		panic(fmt.Sprintf("Cannot add capped option, another option already exists on collection: %s", s.Spec().Name))
	}

	(*s.Spec().Options)[CollectionOptionCapped] = true
	(*s.Spec().Options)[CollectionOptionCappedSize] = size

	return s
}

func (s *MetadataSpec) TTL(expiredAfter int64) *MetadataSpec {
	s.initOptions()

	_, found := (*s.Spec().Options)[CollectionOptionExpiredAfterSeconds]
	if found {
		panic(fmt.Sprintf("Cannot add TTL option, another option already exists on collection: %s", s.Spec().Name))
	}

	(*s.Spec().Options)[CollectionOptionExpiredAfterSeconds] = expiredAfter

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
