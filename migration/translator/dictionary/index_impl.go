package dictionary

import (
	"fmt"

	"internal/convert"

	"github.com/amirkode/go-mongr8/collection"
	"go.mongodb.org/mongo-driver/bson"
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

func (t TranslatedIndex) getFieldsObject() bson.M {
	res := bson.M{}
	for _, field := range t.index.Spec().Fields {
		res[field.Key] = field.Value
	}

	return res
}

func (t TranslatedIndex) getSparseObjectOrNil() *bson.M {
	if t.index.Spec().Sparse {
		return &bson.M{
			"sparse": true,
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

func (t translatedSingleField) getObject() bson.M {
	t.hasAtLeastFieldsLengthValidation(1)
	return t.getFieldsObject()
}

func (t translatedSingleField) getRules() *bson.M { 
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

func (t translatedCompound) getObject() bson.M {
	t.hasAtLeastFieldsLengthValidation(2)	
	return t.getFieldsObject()
}

func (t translatedCompound) getRules() *bson.M { 
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

func (t translatedText) getObject() bson.M {
	t.hasAtLeastFieldsLengthValidation(1)
	field := t.index.Spec().Fields[0]
	return bson.M{
		field.Key: "text",
	}
}

func (t translatedText) getRules() *bson.M {
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

func (t translatedGeospatial2dsphere) getOjbect() bson.M {
	t.hasAtLeastFieldsLengthValidation(1)
	field := t.index.Spec().Fields[0]
	return bson.M{
		field.Key: "2dsphere",
	}
}

func (t translatedGeospatial2dsphere) getRules() *bson.M {
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

func (t translatedUnique) getObject() bson.M {
	t.hasAtLeastFieldsLengthValidation(1)
	return t.getFieldsObject()
}

func (t translatedUnique) getRules() *bson.M {
	res := bson.M{
		"unique": true,
	}
	if t.index.Spec().Sparse {
		res["sparse"] = true
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

func (t translatedPartial) getObject() bson.M {
	t.hasAtLeastFieldsLengthValidation(1)
	return t.getFieldsObject()
}

func (t translatedPartial) getRules() *bson.M {
	t.mustProvideRulesValidation()

	rules := bson.M{
		"partialFilterExpression": convert.MapToBson(*t.index.Spec().Rules),
	}
	if t.index.Spec().Sparse {
		rules["sparse"] = true
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

func (t translatedCollation) getObject() bson.M {
	t.hasAtLeastFieldsLengthValidation(1)
	return t.getFieldsObject()	
}

func (t translatedCollation) getRules() *bson.M {
	t.mustProvideRulesValidation()
	
	rules := bson.M{
		"collation": convert.MapToBson(*t.index.Spec().Rules),
	}
	if t.index.Spec().Sparse {
		rules["sparse"] = true
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

func (t translatedRaw) getObject() bson.M {
	t.hasAtLeastFieldsLengthValidation(1)

	return t.getFieldsObject()
}

func (t translatedRaw) getRules() *bson.M {
	if t.index.Spec().Rules == nil {
		return nil
	}

	rules := convert.MapToBson(*t.index.Spec().Rules)
	if t.index.Spec().Sparse {
		rules["sparse"] = true
	}

	return &rules
}