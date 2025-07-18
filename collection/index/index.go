/*
Copyright (c) 2023-present the go-mongr8 Authors and Contributors
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
		Key: name,
	}

	if len(value) == 1 {
		res.Value = value[0]
	}

	return res
}

func (f IndexField) NestedField(name string) IndexField {
	if f.Key != "" {
		f.Key += "."
	}

	f.Key += name

	return f
}

func (f IndexField) MustHaveValue(_for string) {
	if f.Value == nil {
		panic(fmt.Sprintf("%s: Value should be provided for %s", f.Key, _for))
	}
}

// We will convert this spec into MongoDB API
// indexModel := mongo.IndexModel{
// 	Keys:    keys,
// 	Options: opt,
// }
// collection.Indexes().CreateOne(ctx, indexModel)
// 
// Each type basically has two things:
// - Fields: a list of fields and its value, 
//   e.g: `{"name": 1, "age": -1}` or in the text index case `{"name": "text"}`
// - Rules: we generalize MongoDb options into "Rules"

type Spec struct {
	Type   IndexType
	Fields []IndexField
	Rules  *map[string]interface{}
	Name   *string
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

// Check whether an option is set
func (b *Spec) HasRule(option string) bool {
	if b.Rules == nil {
		return false
	}

	_, ok := (*b.Rules)[option]

	return ok
}

func (b *IndexSpec) Spec() *Spec {
	return b.spec
}

func (b *IndexSpec) InitRules() {
	if b.spec.Rules == nil {
		b.SetRules(map[string]interface{}{})
	}
}

func (b *IndexSpec) SetRules(rules map[string]interface{}) {
	b.spec.Rules = &rules
}

// Panic if index is raw type
func (b *IndexSpec) MustNotRaw() {
	if b.spec.Type == TypeRaw {
		panic("Index type must not be raw")
	}
}

// Set `sparse` option
// Only indexes documents that have particular field
func (b *IndexSpec) AsSparse() *IndexSpec {
	b.MustNotRaw()
	b.InitRules()

	(*b.spec.Rules)[OptionSparse] = true

	return b
}

// Set `background` option.
// Creates the index in the background so it doesn't block reads/writes.
func (b *IndexSpec) AsBackground() *IndexSpec {
	b.MustNotRaw()
	b.InitRules()

	(*b.spec.Rules)[OptionBackground] = true

	return b
}

// Set `unique` option.
// Adds uniqueness constraint to the field.
func (b *IndexSpec) AsUnique() *IndexSpec {
	b.MustNotRaw()
	b.InitRules()

	(*b.spec.Rules)[OptionUnique] = true

	return b
}

// Set `hidden` option.
// Creates an index that is hidden from the query optimizer.
func (b *IndexSpec) AsHidden() *IndexSpec {
	b.MustNotRaw()
	b.InitRules()

	(*b.spec.Rules)[OptionHidden] = true

	return b
}

// Set `partialFilterExpression` option.
// Indexes with particular filters.
func (b *IndexSpec) SetPartialExpression(partialExp map[string]interface{}) *IndexSpec {
	b.MustNotRaw()
	b.InitRules()

	(*b.spec.Rules)[OptionPartialFilterExp] = partialExp

	return b
}

// Set `expireAfterSeconds` option
// Adds TTL index to a timestamp field
func (b *IndexSpec) SetTTL(expireAfterSeconds int32) *IndexSpec {
	b.MustNotRaw()
	b.InitRules()

	(*b.spec.Rules)[OptionTTL] = expireAfterSeconds

	return b
}

// Set `collation` option
// Creates index with custom collation
func (b *IndexSpec) SetCollation(collation map[string]interface{}) *IndexSpec {
	b.MustNotRaw()
	b.InitRules()

	(*b.spec.Rules)[OptionCollation] = collation

	return b
}

func (b *IndexSpec) SetCustomIndexName(name string) *IndexSpec {
	b.spec.Name = &name
	return b
}

func baseIndex(_type IndexType, fields []IndexField, rules *map[string]interface{}) *IndexSpec {
	if len(fields) == 0 {
		panic("Index must have at least a field")
	}

	index := &IndexSpec{
		&Spec{
			Type:   _type,
			Fields: fields,
			Rules:  rules,
		},
	}

	return index
}

func FromIndexSpec(spec *Spec) *IndexSpec {
	return &IndexSpec{
		spec,
	}
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

// reference: https://www.mongodb.com/docs/manual/core/indexes/index-types/index-single/create-single-field-index/
func SingleFieldIndex(field IndexField) *IndexSpec {
	field.MustHaveValue("Single Field Index")
	return defaultIndex(TypeSingleField, []IndexField{field}, nil)
}

// reference: https://www.mongodb.com/docs/manual/core/indexes/index-types/index-compound/create-compound-index/
func CompoundIndex(fields ...IndexField) *IndexSpec {
	for _, _field := range fields {
		_field.MustHaveValue("Compound Index")
	}

	return defaultIndex(TypeCompound, fields, nil)
}

func TextIndex(field IndexField) *IndexSpec {
	field.Value = "text" // used for index name sufix
	// TODO: check other index too
	return defaultIndex(TypeText, []IndexField{field}, nil)
}

func Geospatial2dsphereIndex(field IndexField) *IndexSpec {
	return defaultIndex(TypeGeopatial2dsphere, []IndexField{field}, nil)
}

func HashedIndex(field IndexField) *IndexSpec {
	return defaultIndex(TypeHashed, []IndexField{field}, nil)
}

func RawIndex(fields map[string]interface{}, rules *map[string]interface{}) *IndexSpec {
	return customValueIndex(TypeRaw, fields, rules)
}
