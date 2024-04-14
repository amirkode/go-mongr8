/*
Copyright (c) 2023-present the go-mongr8 Authors and Contributors
[@see Authors file]

Licensed under the MIT License
(https://opensource.org/licenses/MIT)
*/
package field

import (
	"testing"

	"github.com/amirkode/go-mongr8/internal/test"
)

func fieldsAreEqual(a, b *Spec) bool {
	if a.Name != b.Name {
		return false
	}

	if a.Type != b.Type {
		return false
	}

	if (a.Object == nil && b.Object != nil) ||
		(a.Object != nil && b.Object == nil) {
		return false
	}

	if (a.ArrayFields == nil && b.ArrayFields != nil) ||
		(a.ArrayFields != nil && b.ArrayFields == nil) {
		return false
	}

	if a.Object != nil {
		if len(*a.Object) != len(*b.Object) {
			return false
		}

		aObj := map[string]Spec{}
		bObj := map[string]Spec{}
		for _, obj := range *a.Object {
			aObj[obj.Name] = obj
		}
		for _, obj := range *b.Object {
			bObj[obj.Name] = obj
		}

		// check a on b
		for _, currA := range *a.Object {
			currB, ok := bObj[currA.Name]
			if !ok {
				return false
			}

			ok = fieldsAreEqual(&currA, &currB)
			if !ok {
				return false
			}
		}
	}

	if a.ArrayFields != nil {
		if len(*a.ArrayFields) != len(*b.ArrayFields) {
			return false
		}

		aArr := map[string]Spec{}
		bArr := map[string]Spec{}
		for _, arr := range *a.ArrayFields {
			aArr[arr.Name] = arr
		}
		for _, arr := range *b.ArrayFields {
			bArr[arr.Name] = arr
		}

		// check a on b
		for _, currA := range *a.ArrayFields {
			currB, ok := bArr[currA.Name]
			if !ok {
				return false
			}

			ok = fieldsAreEqual(&currA, &currB)
			if !ok {
				return false
			}
		}
	}

	return true
}

func TestStringField(t *testing.T) {
	// case 1: default
	case1Actual := StringField("name")
	case1Expected := baseField("name", TypeString)

	test.AssertTrue(t, fieldsAreEqual(case1Actual.spec, case1Expected.spec), "Case 1: Unexpected field value")
}

func TestInt32Field(t *testing.T) {
	// case 1: default
	case1Actual := Int32Field("name")
	case1Expected := baseField("name", TypeInt32)

	test.AssertTrue(t, fieldsAreEqual(case1Actual.spec, case1Expected.spec), "Case 1: Unexpected field value")
}

func TestInt64Field(t *testing.T) {
	// case 1: default
	case1Actual := Int64Field("name")
	case1Expected := baseField("name", TypeInt64)

	test.AssertTrue(t, fieldsAreEqual(case1Actual.spec, case1Expected.spec), "Case 1: Unexpected field value")
}

func TestDoubleField(t *testing.T) {
	// case 1: default
	case1Actual := DoubleField("name")
	case1Expected := baseField("name", TypeDouble)

	test.AssertTrue(t, fieldsAreEqual(case1Actual.spec, case1Expected.spec), "Case 1: Unexpected field value")
}

func TestBooleanField(t *testing.T) {
	// case 1: default
	case1Actual := BooleanField("name")
	case1Expected := baseField("name", TypeBoolean)

	test.AssertTrue(t, fieldsAreEqual(case1Actual.spec, case1Expected.spec), "Case 1: Unexpected field value")
}

func TestArrayField(t *testing.T) {
	// case 1: array of plain
	case1Actual := ArrayField("name", Int32Field("name"))
	case1Expected := Spec{
		Name: "name",
		Type: TypeArray,
		ArrayFields: &[]Spec{
			{
				Name: "name",
				Type: TypeInt32,
			},
		},
	}

	test.AssertTrue(t, fieldsAreEqual(case1Actual.spec, &case1Expected), "Case 1: Unexpected field value")

	// case 2: array of array
	case2Actual := ArrayField("name",
		ArrayField("child_name",
			StringField("child_child_name"),
		),
	)
	case2Expected := Spec{
		Name: "name",
		Type: TypeArray,
		ArrayFields: &[]Spec{
			{
				Name: "child_name",
				Type: TypeArray,
				ArrayFields: &[]Spec{
					{
						Name: "child_child_name",
						Type: TypeString,
					},
				},
			},
		},
	}

	test.AssertTrue(t, fieldsAreEqual(case2Actual.spec, &case2Expected), "Case 2: Unexpected field value")

	// case 3: array of object
	case3Actual := ArrayField("name",
		ObjectField("",
			StringField("name"),
			Int32Field("score"),
		),
	)
	case3Expected := Spec{
		Name: "name",
		Type: TypeArray,
		ArrayFields: &[]Spec{
			{
				Name: "",
				Type: TypeObject,
				Object: &[]Spec{
					{
						Name: "name",
						Type: TypeString,
					},
					{
						Name: "score",
						Type: TypeInt32,
					},
				},
			},
		},
	}

	test.AssertTrue(t, fieldsAreEqual(case3Actual.spec, &case3Expected), "Case 3: Unexpected field value")
}

func TestObjectField(t *testing.T) {
	// case 1: default
	case1Actual := ObjectField("name",
		StringField("name"),
		Int32Field("score"),
	)
	case1Expected := Spec{
		Name: "name",
		Type: TypeObject,
		Object: &[]Spec{
			{
				Name: "name",
				Type: TypeString,
			},
			{
				Name: "score",
				Type: TypeInt32,
			},
		},
	}

	test.AssertTrue(t, fieldsAreEqual(case1Actual.spec, &case1Expected), "Case 1: Unexpected field value")

	// case 2: object of array
	case2Actual := ObjectField("name",
		ArrayField("child_name",
			StringField("child_child_name"),
		),
	)
	case2Expected := Spec{
		Name: "name",
		Type: TypeObject,
		Object: &[]Spec{
			{
				Name: "child_name",
				Type: TypeArray,
				ArrayFields: &[]Spec{
					{
						Name: "child_child_name",
						Type: TypeString,
					},
				},
			},
		},
	}

	test.AssertTrue(t, fieldsAreEqual(case2Actual.spec, &case2Expected), "Case 2: Unexpected field value")

	// case 3: Object of object
	case3Actual := ObjectField("name",
		ObjectField("child_name",
			StringField("name"),
			Int32Field("score"),
		),
	)
	case3Expected := Spec{
		Name: "name",
		Type: TypeObject,
		Object: &[]Spec{
			{
				Name: "child_name",
				Type: TypeObject,
				Object: &[]Spec{
					{
						Name: "name",
						Type: TypeString,
					},
					{
						Name: "score",
						Type: TypeInt32,
					},
				},
			},
		},
	}

	test.AssertTrue(t, fieldsAreEqual(case3Actual.spec, &case3Expected), "Case 3: Unexpected field value")
}

func TestTimestampField(t *testing.T) {
	// case 1: default
	case1Actual := TimestampField("name")
	case1Expected := baseField("name", TypeTimestamp)

	test.AssertTrue(t, fieldsAreEqual(case1Actual.spec, case1Expected.spec), "Case 1: Unexpected field value")
}

func TestGeoJSONPointField(t *testing.T) {
	// case 1: default
	case1Actual := GeoJSONPointField("name")
	case1Expected := baseField("name", TypeGeoJSONPoint)

	test.AssertTrue(t, fieldsAreEqual(case1Actual.spec, case1Expected.spec), "Case 1: Unexpected field value")
}

func TestGeoJSONLineStringField(t *testing.T) {
	// case 1: default
	case1Actual := GeoJSONLineStringField("name")
	case1Expected := baseField("name", TypeGeoJSONLineString)

	test.AssertTrue(t, fieldsAreEqual(case1Actual.spec, case1Expected.spec), "Case 1: Unexpected field value")
}

func TestGeoJSONPolygonSingleRingField(t *testing.T) {
	// case 1: default
	case1Actual := GeoJSONPolygonSingleRingField("name")
	case1Expected := baseField("name", TypeGeoJSONPolygonSingleRing)

	test.AssertTrue(t, fieldsAreEqual(case1Actual.spec, case1Expected.spec), "Case 1: Unexpected field value")
}

func TestGeoJSONPolygonMultipleRingField(t *testing.T) {
	// case 1: default
	case1Actual := GeoJSONPolygonMultipleRingField("name")
	case1Expected := baseField("name", TypeGeoJSONPolygonMultipleRing)

	test.AssertTrue(t, fieldsAreEqual(case1Actual.spec, case1Expected.spec), "Case 1: Unexpected field value")
}

func TestGeoJSONMultiPointField(t *testing.T) {
	// case 1: default
	case1Actual := GeoJSONMultiPointField("name")
	case1Expected := baseField("name", TypeGeoJSONMultiPoint)

	test.AssertTrue(t, fieldsAreEqual(case1Actual.spec, case1Expected.spec), "Case 1: Unexpected field value")
}

func TestGeoJSONMultiLineStringField(t *testing.T) {
	// case 1: default
	case1Actual := GeoJSONMultiLineStringField("name")
	case1Expected := baseField("name", TypeGeoJSONMultiLineString)

	test.AssertTrue(t, fieldsAreEqual(case1Actual.spec, case1Expected.spec), "Case 1: Unexpected field value")
}

func TestGeoJSONMultiPolygonField(t *testing.T) {
	// case 1: default
	case1Actual := GeoJSONMultiPolygonField("name")
	case1Expected := baseField("name", TypeGeoJSONMultiPolygon)

	test.AssertTrue(t, fieldsAreEqual(case1Actual.spec, case1Expected.spec), "Case 1: Unexpected field value")
}

func TestGeoJSONGeometryCollectionField(t *testing.T) {
	// case 1: default
	case1Actual := GeoJSONGeometryCollectionField("name")
	case1Expected := baseField("name", TypeGeoJSONGeometryCollection)

	test.AssertTrue(t, fieldsAreEqual(case1Actual.spec, case1Expected.spec), "Case 1: Unexpected field value")
}

func TestLegacyCoordinateArrayField(t *testing.T) {
	// case 1: default
	case1Actual := LegacyCoordinateArrayField("name")
	case1Expected := baseField("name", TypeLegacyCoordinateArray)

	test.AssertTrue(t, fieldsAreEqual(case1Actual.spec, case1Expected.spec), "Case 1: Unexpected field value")
}
