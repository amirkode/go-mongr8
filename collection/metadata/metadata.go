/*
Copyright (c) 2023 the go-mongr8 Authors and Contributors
[@see Authors file]

Licensed under the MIT License
(https://opensource.org/licenses/MIT)
*/
package metadata

type Spec struct {
	Name string
}

type MetadataSpec struct {
	spec *Spec
}

func (s *MetadataSpec) Spec() *Spec {
	return s.spec
}

func InitMetadata(name string) *MetadataSpec {
	res := &MetadataSpec{
		&Spec{
			Name: name,	
		},
	}
	return res
}
