/*
Copyright (c) 2023 the go-mongr8 Authors and Contributors
[@see Authors file]

Licensed under the MIT License
(https://opensource.org/licenses/MIT)
*/
package index

type IndexField struct {
	Key   string
	Value interface{}
}

type Spec struct {
	Type   IndexType
	Fields []IndexField
	Rules  *map[string]interface{}
	Sparse bool
}

type BaseSpec struct {
	spec *Spec
}

func (b *BaseSpec) Spec() *Spec {
	return b.spec
}

func (b *BaseSpec) SetRules(rules map[string]interface{}) {
	b.spec.Rules = &rules
}

func (b *BaseSpec) SetSparse(sparse bool) *BaseSpec {
	b.spec.Sparse = sparse

	return b
}

func baseIndex(_type IndexType, fields []IndexField, rules *map[string]interface{}) *BaseSpec {
	index := &BaseSpec{
		&Spec{
			Type:   _type,
			Fields: fields,
			Rules:  rules,
		},
	}

	return index
}

func defaultIndex(_type IndexType, fields []string, rules *map[string]interface{}) *BaseSpec {
	indexFields := make([]IndexField, len(fields))
	for i, field := range fields {
		indexFields[i] = IndexField{
			Key: field,
			Value: int(1),
		}
	}

	return baseIndex(_type, indexFields, rules)
}

func customValueIndex(_type IndexType, fields map[string]interface{}, rules *map[string]interface{}) *BaseSpec {
	indexFields := make([]IndexField, len(fields))
	i := 0
	for key, value := range fields {
		indexFields[i] = IndexField{
			Key: key,
			Value: value,
		}
		i++
	}

	return baseIndex(_type, indexFields, rules)
}

func SingleFieldIndex(field string) *BaseSpec {
	return defaultIndex(TypeSingleField, []string{field}, nil)
}

func CompoundIndex(fields []string) *BaseSpec {
	return defaultIndex(TypeCompound, fields, nil)
}

func TextIndex(field string) *BaseSpec {
	return defaultIndex(TypeText, []string{field}, nil)
}

func Geospatial2dsphereIndex(field string) *BaseSpec {
	return defaultIndex(TypeGeopatial2dsphere, []string{field}, nil)
}

func UniqueIndex(field string) *BaseSpec {
	return defaultIndex(TypeUnique, []string{field}, nil)
}

// partial index custom spec
type partialIndexSpec struct {
	BaseSpec
}

func (s *partialIndexSpec) SetPartialExpression(partialExp map[string]interface{}) *partialIndexSpec {
	s.SetRules(partialExp)

	return s
}

func PartialIndex(fields map[string]interface{}) *partialIndexSpec {
	baseSpec := customValueIndex(TypePartial, fields, nil)
	return &partialIndexSpec{
		*baseSpec,
	}
}

// collation index custom spec
type collationIndexSpec struct {
	BaseSpec
}

func (s *collationIndexSpec) SetCollation(collation map[string]interface{}) {
	s.SetRules(collation)
}

func CollationIndex(field string) *collationIndexSpec {
	baseSpec := defaultIndex(TypeCollation, []string{field}, nil)
	return &collationIndexSpec{
		*baseSpec,
	}
}

func RawIndex(fields map[string]interface{}, rules *map[string]interface{}) *BaseSpec {
	return customValueIndex(TypeRaw, fields, rules)
}