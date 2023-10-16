package sync_strategy

import (
	"internal/test"
	"testing"

	"github.com/amirkode/go-mongr8/collection/field"
)

func TestSignedFieldGetKey(t *testing.T) {
	fieldName := "test_field"
	stringField := field.StringField(fieldName)
	signedField := SignedField{
		Field: stringField,
	}

	test.AssertEqual(t, signedField.Key(), fieldName, "Key is not equal to field name")
}

func TestSignedFieldIntersect(t *testing.T) {
	// test intersect on array with different children
	arrayField1 := SignedField{
		Field: field.ArrayField("array_field").
			AddArrayField(field.StringField("string")).
			AddArrayField(field.Int32Field("int32")).
			AddArrayField(field.BooleanField("boolean")),
		Sign: SignPlus,
	}
	arrayField2 := SignedField{
		Field: field.ArrayField("array_field").
			AddArrayField(field.StringField("string")).
			AddArrayField(field.StringField("string_2")).
			AddArrayField(field.Int32Field("int32")),
		Sign: SignPlus,
	}
	
	intersection := arrayField1.Intersect(arrayField2)
	// there must be an intersection, since there's additional child
	// on arrayField1 (boolean field)
	test.AssertTrue(t, intersection != nil && len(*intersection) == 2, "Intersection is not found")
	// check on the intersection result
	for _, i := range *intersection {
		if i.Sign == SignPlus {
			test.AssertEqual(t, i.Field.Spec().Type, field.TypeArray, "Unexpected parent type")
			test.AssertTrue(t, i.Field.Spec().ArrayFields != nil, "Array children not found")
			test.AssertEqual(t, (*i.Field.Spec().ArrayFields)[0].Type, field.TypeBoolean, "Unexpected child type")
			test.AssertEqual(t, (*i.Field.Spec().ArrayFields)[0].Name, "boolean", "Unexpected child name")
		} else if i.Sign == SignMinus {
			test.AssertEqual(t, i.Field.Spec().Type, field.TypeArray, "Unexpected parent type")
			test.AssertTrue(t, i.Field.Spec().ArrayFields != nil, "Array children not found")
			test.AssertEqual(t, (*i.Field.Spec().ArrayFields)[0].Type, field.TypeString, "Unexpected child type")
			test.AssertEqual(t, (*i.Field.Spec().ArrayFields)[0].Name, "string_2", "Unexpected child name")
		} else {
			t.Error("Intersection sign is not expected")
		}
	}
}