package dictionary

import (
	"fmt"

	"github.com/amirkode/go-mongr8/collection"
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

	translatedUnique struct {
		TranslatedIndex
	}

	translatedPartial struct {
		TranslatedIndex
	}

	translatedCollation struct {
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

func (t TranslatedIndex) getSparseObjectOrNil() *map[string]interface{} {
	if t.index.Spec().Sparse {
		return &map[string]interface{}{
			"sparse": Boolean(true),
		}
	}

	return nil
}

// translation for single field index
func newTranslatedSingleField(index collection.Index) translatedSingleField {
	return translatedSingleField{
		TranslatedIndex{
			index: index,
		},
	}
}

func (t translatedSingleField) getObject() map[string]interface{} {
	t.hasAtLeastFieldsLengthValidation(1)
	return t.getFieldsObject()
}

func (t translatedSingleField) getRules() *map[string]interface{} {
	return t.getSparseObjectOrNil()
}

// translation for compound index
func newTranslatedCompound(index collection.Index) translatedCompound {
	return translatedCompound{
		TranslatedIndex{
			index: index,
		},
	}
}

func (t translatedCompound) getObject() map[string]interface{} {
	t.hasAtLeastFieldsLengthValidation(2)
	return t.getFieldsObject()
}

func (t translatedCompound) getRules() *map[string]interface{} {
	return t.getSparseObjectOrNil()
}

// translation for text index
func newTranslatedText(index collection.Index) translatedText {
	return translatedText{
		TranslatedIndex{
			index: index,
		},
	}
}

func (t translatedText) getObject() map[string]interface{} {
	t.hasAtLeastFieldsLengthValidation(1)
	field := t.index.Spec().Fields[0]
	return map[string]interface{}{
		field.Key: String("text"),
	}
}

func (t translatedText) getRules() *map[string]interface{} {
	return t.getSparseObjectOrNil()
}

// translation for geospatial: 2dspehere index
func newGeospatial2dsphere(index collection.Index) translatedGeospatial2dsphere {
	return translatedGeospatial2dsphere{
		TranslatedIndex{
			index: index,
		},
	}
}

func (t translatedGeospatial2dsphere) getOjbect() map[string]interface{} {
	t.hasAtLeastFieldsLengthValidation(1)
	field := t.index.Spec().Fields[0]
	return map[string]interface{}{
		field.Key: String("2dsphere"),
	}
}

func (t translatedGeospatial2dsphere) getRules() *map[string]interface{} {
	return t.getSparseObjectOrNil()
}

// translation for unique index
func newTranslatedUnique(index collection.Index) translatedUnique {
	return translatedUnique{
		TranslatedIndex{
			index: index,
		},
	}
}

func (t translatedUnique) getObject() map[string]interface{} {
	t.hasAtLeastFieldsLengthValidation(1)
	return t.getFieldsObject()
}

func (t translatedUnique) getRules() *map[string]interface{} {
	res := map[string]interface{}{
		"unique": Boolean(true),
	}
	if t.index.Spec().Sparse {
		res["sparse"] = Boolean(true)
	}

	return &res
}

// translation for partial index
func newTranslatedPartial(index collection.Index) translatedPartial {
	return translatedPartial{
		TranslatedIndex{
			index: index,
		},
	}
}

func (t translatedPartial) getObject() map[string]interface{} {
	t.hasAtLeastFieldsLengthValidation(1)
	return t.getFieldsObject()
}

func (t translatedPartial) getRules() *map[string]interface{} {
	t.mustProvideRulesValidation()

	rules := map[string]interface{}{
		"partialFilterExpression": ConvertAnyToValueType(*t.index.Spec().Rules),
	}
	if t.index.Spec().Sparse {
		rules["sparse"] = Boolean(true)
	}

	return &rules
}

// translation for collation index
func newCollation(index collection.Index) translatedCollation {
	return translatedCollation{
		TranslatedIndex{
			index: index,
		},
	}
}

func (t translatedCollation) getObject() map[string]interface{} {
	t.hasAtLeastFieldsLengthValidation(1)
	return t.getFieldsObject()
}

func (t translatedCollation) getRules() *map[string]interface{} {
	t.mustProvideRulesValidation()

	rules := map[string]interface{}{
		"collation": ConvertAnyToValueType(*t.index.Spec().Rules),
	}
	if t.index.Spec().Sparse {
		rules["sparse"] = Boolean(true)
	}

	return &rules
}

// translation for raw definition index
func newRaw(index collection.Index) translatedRaw {
	return translatedRaw{
		TranslatedIndex{
			index: index,
		},
	}
}

func (t translatedRaw) getObject() map[string]interface{} {
	t.hasAtLeastFieldsLengthValidation(1)

	return t.getFieldsObject()
}

func (t translatedRaw) getRules() *map[string]interface{} {
	if t.index.Spec().Rules == nil {
		return nil
	}

	rules := ConvertAnyToValueType(*t.index.Spec().Rules).(map[string]interface{})
	if t.index.Spec().Sparse {
		rules["sparse"] = Boolean(true)
	}

	return &rules
}
