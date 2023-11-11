/*
Copyright (c) 2023 the go-mongr8 Authors and Contributors
[@see Authors file]

Licensed under the MIT License
(https://opensource.org/licenses/MIT)
*/
package dictionary

import (
	"fmt"

	"github.com/amirkode/go-mongr8/collection"
	"github.com/amirkode/go-mongr8/collection/index"
)

type (
	TranslatedIndex struct {
		TranslatedIndexIf
		index collection.Index
	}

	translatedSingleField struct {
		TranslatedIndex
	}

	translatedCompound struct {
		TranslatedIndex
	}

	translatedText struct {
		TranslatedIndex
	}

	translatedGeospatial2dsphere struct {
		TranslatedIndex
	}

	translatedHashed struct {
		TranslatedIndex
	}

	translatedRaw struct {
		TranslatedIndex
	}
)

func (t TranslatedIndex) hasAtLeastFieldsLengthValidation(minLength int) {
	if len(t.index.Spec().Fields) < minLength {
		panic(fmt.Sprintf("Provided field array must at least have %d length in the index definition", minLength))
	}
}

func (t TranslatedIndex) mustProvideRulesValidation() {
	if t.index.Spec().Rules == nil {
		panic(fmt.Sprintf("Rules must be provided for index type: %v", t.index.Spec().Type))
	}
}

func (t TranslatedIndex) getFieldsObject() map[string]interface{} {
	res := map[string]interface{}{}
	for _, field := range t.index.Spec().Fields {
		res[field.Key] = ConvertAnyToValueType(field.Value)
	}

	return res
}

func (t TranslatedIndex) getRules() *map[string]interface{} {
	if t.index.Spec().Rules == nil {
		return nil
	}

	rules := ConvertAnyToValueType(*t.index.Spec().Rules).(map[string]interface{})

	return &rules
}

// translation for single field index
func newTranslatedSingleFieldIndex(index collection.Index) translatedSingleField {
	return translatedSingleField{
		TranslatedIndex{
			index: index,
		},
	}
}

func (t translatedSingleField) GetObject() map[string]interface{} {
	t.hasAtLeastFieldsLengthValidation(1)
	return t.getFieldsObject()
}

func (t translatedSingleField) GetRules() *map[string]interface{} {
	return t.getRules()
}

// translation for compound index
func newTranslatedCompoundIndex(index collection.Index) translatedCompound {
	return translatedCompound{
		TranslatedIndex{
			index: index,
		},
	}
}

func (t translatedCompound) GetObject() map[string]interface{} {
	t.hasAtLeastFieldsLengthValidation(2)
	return t.getFieldsObject()
}

func (t translatedCompound) GetRules() *map[string]interface{} {
	return t.getRules()
}

// translation for text index
func newTranslatedTextIndex(index collection.Index) translatedText {
	return translatedText{
		TranslatedIndex{
			index: index,
		},
	}
}

func (t translatedText) GetObject() map[string]interface{} {
	t.hasAtLeastFieldsLengthValidation(1)
	field := t.index.Spec().Fields[0]
	return map[string]interface{}{
		field.Key: String("text"),
	}
}

func (t translatedText) GetRules() *map[string]interface{} {
	return t.getRules()
}

// translation for geospatial: 2dspehere index
func newTranslatedGeospatial2dsphereIndex(index collection.Index) translatedGeospatial2dsphere {
	return translatedGeospatial2dsphere{
		TranslatedIndex{
			index: index,
		},
	}
}

func (t translatedGeospatial2dsphere) GetObject() map[string]interface{} {
	t.hasAtLeastFieldsLengthValidation(1)
	field := t.index.Spec().Fields[0]
	return map[string]interface{}{
		field.Key: String("2dsphere"),
	}
}

func (t translatedGeospatial2dsphere) GetRules() *map[string]interface{} {
	return t.getRules()
}

// translation for hashed index
func newTranslatedHashedIndex(index collection.Index) translatedHashed {
	return translatedHashed{
		TranslatedIndex{
			index: index,
		},
	}
}

func (t translatedHashed) GetObject() map[string]interface{} {
	t.hasAtLeastFieldsLengthValidation(1)
	field := t.index.Spec().Fields[0]
	return map[string]interface{}{
		field.Key: String("hashed"),
	}
}

func (t translatedHashed) GetRules() *map[string]interface{} {
	return t.getRules()
}

// translation for raw definition index
func newTranslatedRawIndex(index collection.Index) translatedRaw {
	return translatedRaw{
		TranslatedIndex{
			index: index,
		},
	}
}

func (t translatedRaw) GetObject() map[string]interface{} {
	t.hasAtLeastFieldsLengthValidation(1)
	return t.getFieldsObject()
}

func (t translatedRaw) GetRules() *map[string]interface{} {
	return t.getRules()
}

// map index to correct translated index
func GetTranslatedIndex(_index collection.Index) TranslatedIndexIf {
	switch _index.Spec().Type {
	case index.TypeSingleField:
		return newTranslatedSingleFieldIndex(_index)
	case index.TypeCompound:
		return newTranslatedCompoundIndex(_index)
	case index.TypeText:
		return newTranslatedTextIndex(_index)
	case index.TypeGeopatial2dsphere:
		return newTranslatedGeospatial2dsphereIndex(_index)
	case index.TypeHashed:
		return newTranslatedHashedIndex(_index)
	case index.TypeRaw:
		return newTranslatedRawIndex(_index)
	}

	return TranslatedIndex{}
}
