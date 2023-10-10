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

func Field(name string, value interface{}) IndexField {
	return IndexField{
		Key:   name,
		Value: value,
	}
}

func (f IndexField) NestedField(name string) IndexField {
	f.Key += "." + name

	return f
}

type Spec struct {
	Type   IndexType
	Fields []IndexField
	Rules  *map[string]interface{}
	Sparse bool
}

type IndexSpec struct {
	spec *Spec
}

func (b *IndexSpec) Spec() *Spec {
	return b.spec
}

func (b *IndexSpec) SetRules(rules map[string]interface{}) {
	b.spec.Rules = &rules
}

func (b *IndexSpec) SetSparse(sparse bool) *IndexSpec {
	b.spec.Sparse = sparse

	return b
}

func baseIndex(_type IndexType, fields []IndexField, rules *map[string]interface{}) *IndexSpec {
	index := &IndexSpec{
		&Spec{
			Type:   _type,
			Fields: fields,
			Rules:  rules,
		},
	}

	return index
}

func defaultIndex(_type IndexType, fields []IndexField, rules *map[string]interface{}) *IndexSpec {
	// Some operation here (?)
	return baseIndex(_type, fields, rules)
}

func customValueIndex(_type IndexType, fields map[string]interface{}, rules *map[string]interface{}) *IndexSpec {
	indexFields := make([]IndexField, len(fields))
	i := 0
	for key, value := range fields {
		indexFields[i] = IndexField{
			Key:   key,
			Value: value,
		}
		i++
	}

	return baseIndex(_type, indexFields, rules)
}

func SingleFieldIndex(field IndexField) *IndexSpec {
	return defaultIndex(TypeSingleField, []IndexField{field}, nil)
}

func CompoundIndex(fields ...IndexField) *IndexSpec {
	return defaultIndex(TypeCompound, fields, nil)
}

func TextIndex(field IndexField) *IndexSpec {
	return defaultIndex(TypeText, []IndexField{field}, nil)
}

func Geospatial2dsphereIndex(field IndexField) *IndexSpec {
	return defaultIndex(TypeGeopatial2dsphere, []IndexField{field}, nil)
}

func UniqueIndex(field IndexField) *IndexSpec {
	return defaultIndex(TypeUnique, []IndexField{field}, nil)
}

// partial index custom spec
type partialIndexSpec struct {
	IndexSpec
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
	IndexSpec
}

func (s *collationIndexSpec) SetCollation(collation map[string]interface{}) {
	s.SetRules(collation)
}

func CollationIndex(field IndexField) *collationIndexSpec {
	baseSpec := defaultIndex(TypeCollation, []IndexField{field}, nil)
	return &collationIndexSpec{
		*baseSpec,
	}
}

func RawIndex(fields map[string]interface{}, rules *map[string]interface{}) *IndexSpec {
	return customValueIndex(TypeRaw, fields, rules)
}
