/*
Copyright (c) 2023 the go-mongr8 Authors and Contributors
[@see Authors file]

Licensed under the MIT License
(https://opensource.org/licenses/MIT)
*/
package index

import (
	"fmt"
	"reflect"
)

type IndexField struct {
	Key   string
	Value interface{}
}

func Field(name string, value ...interface{}) IndexField {
	if len(value) > 1 {
		panic("Index field value at most declared once")
	}

	res := IndexField{
		Key:   name,
	}

	if len(value) == 1 {
		res.Value = value[0]
	}

	return res
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
	// TODO: implement custom index name
	Name  *string
}

type IndexSpec struct {
	spec *Spec
}

// this is used for compare two index whether
// they both have the exact same structure
func (s *Spec) GetKey() string {
	// we can have unique index by comparing entire structure
	// with these combinations:
	// - type
	// - index fields
	// - rules
	// - sparse option
	// key := fmt.Sprintf("%v", f.Spec())
	key := string(s.Type)
	key += fmt.Sprintf("%v", s.Sparse)

	for _, field := range s.Fields {
		key += field.Key
		key += fmt.Sprintf("%v", field.Value)
	}

	if s.Rules != nil {
		key += fmt.Sprintf("%v", *s.Rules)
	}

	return key
}

func (s *Spec) GetName() string {
	if s.Name != nil {
		return *s.Name
	}
	
	res := ""
	for _, field := range s.Fields {
		if res == "" {
			res = fmt.Sprintf("%s_%v", field.Key, field.Value)
		} else {
			res = fmt.Sprintf("%s_%s_%v", res, field.Key, field.Value)
		}
	}

	var getPath func(curr interface{}) string
	getPath = func(curr interface{}) string {
		path := ""
		if reflect.TypeOf(curr).Kind() == reflect.Map &&
			reflect.TypeOf(curr).Key().Kind() == reflect.String &&
			reflect.TypeOf(curr).Elem().Kind() == reflect.Interface {
			mp := curr.(map[string]interface{})
			for key, val := range mp {
				res = fmt.Sprintf("%s_%s%s", res, key, getPath(val))
			}
		} else if reflect.TypeOf(curr).Kind() == reflect.Slice &&
			reflect.TypeOf(curr).Elem().Kind() == reflect.Interface {
			arr := curr.([]interface{})
			for _, val := range arr {
				res = fmt.Sprintf("%s%s", res, getPath(val))
			}
		} else {
			res = fmt.Sprintf("%s_%v", res, curr)
		}

		return path
	}

	if s.Rules != nil {
		res += getPath(*s.Rules)
	}

	return res
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

func (b *IndexSpec) SetCustomIndexName(name string) *IndexSpec {
	b.spec.Name = &name
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

func NewIndexField(key string, value interface{}) IndexField {
	return IndexField{
		Key:   key,
		Value: value,
	}
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

// TODO: apply this generally
func UniqueIndex(field IndexField) *IndexSpec {
	rules := map[string]interface{}{
		"unique": true,
	}
	return defaultIndex(TypeUnique, []IndexField{field}, &rules)
}

// partial index custom spec
type partialIndexSpec struct {
	IndexSpec
}

func (s *partialIndexSpec) SetPartialExpression(partialExp map[string]interface{}) *partialIndexSpec {
	rules := map[string]interface{}{
		"partialFilterExpression": partialExp,
	}

	s.SetRules(rules)

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

func (s *collationIndexSpec) SetCollation(collation map[string]interface{}) *collationIndexSpec {
	rules := map[string]interface{}{
		"collation": collation,
	}

	s.SetRules(rules)

	return s
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
