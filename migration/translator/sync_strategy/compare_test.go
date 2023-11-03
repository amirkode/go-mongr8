package sync_strategy

import (
	"testing"

	"internal/test"
	"internal/util"

	"github.com/amirkode/go-mongr8/collection/field"
	"github.com/amirkode/go-mongr8/collection/index"
	"github.com/amirkode/go-mongr8/collection/metadata"
)

func TestSignedFieldGetKey(t *testing.T) {
	// test equal key equal to field name
	case1FieldName := "test_field"
	case1StringField := field.StringField(case1FieldName)
	case1SignedField := SignedField{
		Field: case1StringField,
	}

	test.AssertEqual(t, case1SignedField.Key(), case1FieldName, "Key is not equal to field name")
}

func TestSignedFieldIntersect(t *testing.T) {
	// test intersect on array with different children
	case1Field1 := SignedField{
		Field: field.ArrayField("array_field", 
			field.ObjectField("name not required", 
				field.BooleanField("state1"),
				field.BooleanField("state3"),
			),
		),
			// TODO: update this
			// for now, we do not support multipe types array
			// AddArrayField(field.StringField("string")),
			// AddArrayField(field.Int32Field("int32")).
			// AddArrayField(field.BooleanField("boolean")),
		Sign: SignPlus,
	}
	case1Field2 := SignedField{
		Field: field.ArrayField("array_field", 
			field.ObjectField("name not required", 
				field.BooleanField("state1"),
				field.BooleanField("state2"),
			),
		),
			// TODO: update this
			// for now, we do not support multipe types array
			// AddArrayField(field.StringField("string")).
			// AddArrayField(field.StringField("string_2")).
			// AddArrayField(field.Int32Field("int32")),
		Sign: SignPlus,
	}

	case1Intersection := case1Field1.Intersect(case1Field2)
	// there must be an intersection, since there's additional child
	// on arrayField1 (boolean field)
	test.AssertTrue(t, case1Intersection != nil && len(*case1Intersection) == 2, "Case 1: Intersection is not expected")
	// check on the intersection result
	for _, i := range *case1Intersection {
		if !util.InListEq(i.Sign, []EntitySign{SignPlus, SignMinus}) {
			msg := "Case 1: Intersection sign is not expected"
			t.Error(msg)
			panic(msg)
		}

		test.AssertEqual(t, i.Field.Spec().Type, field.TypeArray, "Case 1: Unexpected parent type")
		test.AssertTrue(t, i.Field.Spec().ArrayFields != nil, "Case 1: Array children not found")
		test.AssertEqual(t, (*i.Field.Spec().ArrayFields)[0].Type, field.TypeObject, "Case 1: Unexpected child type")
		test.AssertTrue(t, (*i.Field.Spec().ArrayFields)[0].Object != nil, "Case 1: Object children not found")
		test.AssertEqual(t, (*(*i.Field.Spec().ArrayFields)[0].Object)[0].Type, field.TypeBoolean, "Case 1: Unexpected object child type")
		test.AssertTrue(t, util.InListEq((*(*i.Field.Spec().ArrayFields)[0].Object)[0].Name, []string{"state2", "state3"}) , "Case 1: Unexpected object child name")
	}

	// test intersect on object with different children
	case2Field1 := SignedField{
		Field: field.ObjectField("object_field",
			field.StringField("string"),
			field.Int32Field("int32"),
			field.BooleanField("boolean"),
			field.ObjectField("nested_object",
				field.StringField("nested_string1"),
				field.StringField("nested_string2"),
			),
		),
		Sign: SignPlus,
	}
	case2Field2 := SignedField{
		Field: field.ObjectField("object_field",
			field.StringField("string"),
			field.StringField("string_2"),
			field.Int32Field("int32"),
			field.ObjectField("nested_object",
				field.StringField("nested_string1"),
				field.StringField("nested_string2"),
			),
		),
		Sign: SignPlus,
	}

	case2Intersection := case2Field1.Intersect(case2Field2)
	// there must be an intersection, since there's additional child
	// on arrayField1 (boolean field)
	test.AssertTrue(t, case2Intersection != nil && len(*case2Intersection) == 2, "Case 2: Object Intersection is not expected")
	// check on the intersection result
	for _, i := range *case2Intersection {
		if i.Sign == SignPlus {
			test.AssertEqual(t, i.Field.Spec().Type, field.TypeObject, "Case 2: Unexpected parent type")
			test.AssertTrue(t, i.Field.Spec().Object != nil, "Case 2: Object items not found")
			test.AssertEqual(t, (*i.Field.Spec().Object)[0].Type, field.TypeBoolean, "Case 2: Unexpected child type")
			test.AssertEqual(t, (*i.Field.Spec().Object)[0].Name, "boolean", "Case 2: Unexpected child name")
		} else if i.Sign == SignMinus {
			test.AssertEqual(t, i.Field.Spec().Type, field.TypeObject, "Case 2: Unexpected parent type")
			test.AssertTrue(t, i.Field.Spec().Object != nil, "Case 2: Object items not found")
			test.AssertEqual(t, (*i.Field.Spec().Object)[0].Type, field.TypeString, "Case 2: Unexpected child type")
			test.AssertEqual(t, (*i.Field.Spec().Object)[0].Name, "string_2", "Case 2: Unexpected child name")
		} else {
			t.Error("Case 2: Intersection sign is not expected")
		}
	}

	// test no difference in intersection found
	case3Field1 := SignedField{
		Field: field.StringField("string"),
	}
	case3Field2 := SignedField{
		Field: field.StringField("string"),
	}
	case3Intersection := case3Field1.Intersect(case3Field2)
	test.AssertTrue(t, case3Intersection != nil && len(*case3Intersection) == 0, "Case 3: Intersection is not expected")

	// test intersection with conversion
	case4Field1 := SignedField{
		Field: field.DoubleField("number"),
	}
	case4Field2 := SignedField{
		Field: field.Int64Field("number"),
	}
	case4Intersection := case4Field1.Intersect(case4Field2)

	test.AssertTrue(t, case4Intersection != nil && len(*case4Intersection) == 1, "Case 4: Intersection is not expected")
	test.AssertEqual(t, (*case4Intersection)[0].Sign, SignConvert, "Case 4: Intersection sign is not convert")
	test.AssertEqual(t, (*case4Intersection)[0].Spec().Type, field.TypeDouble, "Case 4: Conversion Type is not Double")
	test.AssertEqual(t, (*case4Intersection)[0].convertFrom.Spec().Type, field.TypeInt64, "Case 4: ConversionFrom Type is not Int64")
}

func TestSignedFieldUnion(t *testing.T) {
	// test union resulting plus and negative sign
	case1FieldArr1 := []SignedField{
		{
			Field: field.BooleanField("active"),
		},
		{
			Field: field.StringField("name"),
		},
	}
	case1FieldArr2 := []SignedField{
		{
			Field: field.ArrayField("scores",
				field.Int32Field("name not required"),
			),
		},
	}
	case1Union := Union(case1FieldArr1, case1FieldArr2)

	test.AssertTrue(t, len(case1Union) == 3, "Case 1: Union length is wrong")
	case1UnionHasPlus := false
	case1UnionHasMinus := false
	for _, u := range case1Union {
		if u.Sign == SignPlus {
			case1UnionHasPlus = true
		} else if u.Sign == SignMinus {
			case1UnionHasMinus = true
		}

		test.AssertTrue(t, util.InListEq(u.Spec().Name, []string{"active", "name", "scores"}), "Case 1: Invalid union field name")
	}
	test.AssertTrue(t, case1UnionHasPlus, "Case 1: Union does not have plus")
	test.AssertTrue(t, case1UnionHasMinus, "Case 1: Union does not have minus")

	// test union resulting empty list, because of not differences
	case2FieldArr1 := []SignedField{
		{
			Field: field.BooleanField("active"),
		},
		{
			Field: field.StringField("name"),
		},
	}
	case2FieldArr2 := []SignedField{
		{
			Field: field.BooleanField("active"),
		},
		{
			Field: field.StringField("name"),
		},
	}
	case2Union := Union(case2FieldArr1, case2FieldArr2)

	test.AssertTrue(t, len(case2Union) == 0, "Case 2: Union is not empty")
}

func TestSignedIndexGetKey(t *testing.T) {
	// test two index keys are same
	case1Index1 := SignedIndex{
		Index: index.CompoundIndex(
			index.Field("name", 1),
			index.Field("age", -1),
		),
	}
	case1Index2 := SignedIndex{
		Index: index.CompoundIndex(
			index.Field("name", 1),
			index.Field("age", -1),
		),
	}

	test.AssertEqual(t, case1Index1.Key(), case1Index2.Key(), "Case 1: Keys are different")

	// test two index keys are different
	case2Index1 := SignedIndex{
		Index: index.CompoundIndex(
			index.Field("name", 1),
			index.Field("age", 1),
		),
	}
	case2Index2 := SignedIndex{
		Index: index.CompoundIndex(
			index.Field("name", 1),
			index.Field("age", -1),
		),
	}

	test.AssertNotEqual(t, case2Index1.Key(), case2Index2.Key(), "Case 1: Keys are same")
}

func TestSignedIndexInteract(t *testing.T) {
	// test intersection with same index keys
	case1Index1 := SignedIndex{
		Index: index.CompoundIndex(
			index.Field("name", 1),
			index.Field("age", -1),
		),
	}
	case1Index2 := SignedIndex{
		Index: index.CompoundIndex(
			index.Field("name", 1),
			index.Field("age", -1),
		),
	}

	case1Intersect := case1Index1.Intersect(case1Index2)

	test.AssertTrue(t, case1Intersect == nil, "Case 1: Intersection is not nil")

	// test intersection with different index keys
	case2Index1 := SignedIndex{
		Index: index.CompoundIndex(
			index.Field("name", 1),
			index.Field("age", -1),
		),
	}
	case2Index2 := SignedIndex{
		Index: index.CompoundIndex(
			index.Field("name", 1),
			index.Field("age", 1),
		),
	}

	case2Intersect := case2Index1.Intersect(case2Index2)

	test.AssertTrue(t, case2Intersect == nil, "Case 2: Intersection is not nil")
}

func TestSignedIndexUnion(t *testing.T) {
	// test union resulting plus and minus
	case1IndexArr1 := []SignedIndex{
		{
			Index: index.CompoundIndex(index.Field("name", 1)),
		},
	}
	case1IndexArr2 := []SignedIndex{
		{
			Index: index.CompoundIndex(index.Field("story", 1)),
		},
	}
	case1Union := Union(case1IndexArr1, case1IndexArr2)

	test.AssertTrue(t, len(case1Union) == 2, "Case 1: Union length is wrong")
	case1UnionHasPlus := false
	case1UnionHasMinus := false
	for _, u := range case1Union {
		if u.Sign == SignPlus {
			case1UnionHasPlus = true
		} else if u.Sign == SignMinus {
			case1UnionHasMinus = true
		}

		test.AssertTrue(t, util.InListEq(u.Key(), []string{
			case1IndexArr1[0].Key(),
			case1IndexArr2[0].Key(),
		}), "Case 1: Invalid union index key")
	}
	test.AssertTrue(t, case1UnionHasPlus, "Case 1: Union does not have plus")
	test.AssertTrue(t, case1UnionHasMinus, "Case 1: Union does not have minus")

	// test union on same keys
	case2IndexArr1 := []SignedIndex{
		{
			Index: index.CompoundIndex(index.Field("name", 1)),
		},
	}
	case2IndexArr2 := []SignedIndex{
		{
			Index: index.CompoundIndex(index.Field("name", 1)),
		},
	}
	case2Union := Union(case2IndexArr1, case2IndexArr2)

	test.AssertTrue(t, len(case2Union) == 0, "Case 2: Union is not empty")
}

func TestSignedMetadataGetKey(t *testing.T) {
	// test two keys are same
	case1Metadata1 := SignedMetadata{
		Metadata: metadata.InitMetadata("users").Capped(100000),
	}
	case1Metadata2 := SignedMetadata{
		Metadata: metadata.InitMetadata("users").Capped(100000),
	}

	test.AssertEqual(t, case1Metadata1.Key(), case1Metadata2.Key(), "Case 1: Keys are different")

	// test two keys are different
	case2Metadata1 := SignedMetadata{
		Metadata: metadata.InitMetadata("users").Capped(100000),
	}
	case2Metadata2 := SignedMetadata{
		Metadata: metadata.InitMetadata("users").Capped(50000),
	}

	test.AssertNotEqual(t, case2Metadata1.Key(), case2Metadata2.Key(), "Case 2: Keys are same")
}

func TestSignedMetadataIntersect(t *testing.T) {
	// test two keys are same
	case1Metadata1 := SignedMetadata{
		Metadata: metadata.InitMetadata("users").Capped(100000),
	}
	case1Metadata2 := SignedMetadata{
		Metadata: metadata.InitMetadata("users").Capped(100000),
	}
	case1Intersection := case1Metadata1.Intersect(case1Metadata2)

	test.AssertTrue(t, case1Intersection == nil, "Case 1: Intersection is not nil")

	// test two keys are different
	case2Metadata1 := SignedMetadata{
		Metadata: metadata.InitMetadata("users").Capped(100000),
	}
	case2Metadata2 := SignedMetadata{
		Metadata: metadata.InitMetadata("users").Capped(50000),
	}
	case2Intersection := case2Metadata1.Intersect(case2Metadata2)

	test.AssertTrue(t, case2Intersection == nil, "Case 2: Intersection is not nil")
}

func TestSignedMetadataUnion(t *testing.T) {
	// test union resulting plus and minus
	case1MetadataArr1 := []SignedMetadata{
		{
			Metadata: metadata.InitMetadata("users").Capped(100000),
		},
	}
	case1MetadataArr2 := []SignedMetadata{
		{
			Metadata: metadata.InitMetadata("users").Capped(50000),
		},
	}
	case1Union := Union(case1MetadataArr1, case1MetadataArr2)

	test.AssertTrue(t, len(case1Union) == 2, "Case 1: Union length is wrong")
	case1UnionHasPlus := false
	case1UnionHasMinus := false
	for _, u := range case1Union {
		if u.Sign == SignPlus {
			case1UnionHasPlus = true
		} else if u.Sign == SignMinus {
			case1UnionHasMinus = true
		}

		test.AssertTrue(t, util.InListEq(u.Key(), []string{
			case1MetadataArr1[0].Key(),
			case1MetadataArr2[0].Key(),
		}), "Case 1: Invalid union index key")
	}
	test.AssertTrue(t, case1UnionHasPlus, "Case 1: Union does not have plus")
	test.AssertTrue(t, case1UnionHasMinus, "Case 1: Union does not have minus")

	// test union on same keys
	case2MetadataArr1 := []SignedMetadata{
		{
			Metadata: metadata.InitMetadata("users").Capped(100000),
		},
	}
	case2MetadataArr2 := []SignedMetadata{
		{
			Metadata: metadata.InitMetadata("users").Capped(100000),
		},
	}
	case2Union := Union(case2MetadataArr1, case2MetadataArr2)

	test.AssertTrue(t, len(case2Union) == 0, "Case 2: Union is not empty")
}

func TestSignedCollectionGetKey(t *testing.T) {
	// test two keys are same
	case1Collection1 := SignedCollection{
		Metadata: SignedMetadata{
			Metadata: metadata.InitMetadata("users"),
		},
	}
	case1Collection2 := SignedCollection{
		Metadata: SignedMetadata{
			Metadata: metadata.InitMetadata("users"),
		},
	}

	test.AssertEqual(t, case1Collection1.Key(), case1Collection2.Key(), "Case 1: Keys are different")

	// test two keys are different
	case2Collection1 := SignedCollection{
		Metadata: SignedMetadata{
			Metadata: metadata.InitMetadata("users"),
		},
	}
	case2Collection2 := SignedCollection{
		Metadata: SignedMetadata{
			Metadata: metadata.InitMetadata("customers"),
		},
	}

	test.AssertNotEqual(t, case2Collection1.Key(), case2Collection2.Key(), "Case 2: Keys are same")
}

func TestSignedCollectionIntersect(t *testing.T) {
	// test intersection with same index keys
	case1Collection1 := SignedCollection{
		Metadata: SignedMetadata{
			Metadata: metadata.InitMetadata("users"),
		},
		Indexes: []SignedIndex{
			{
				Index: index.CompoundIndex(
					index.Field("name", 1),
					index.Field("age", -1),
				),
			},
		},
	}
	case1Collection2 := SignedCollection{
		Metadata: SignedMetadata{
			Metadata: metadata.InitMetadata("users"),
		},
		Indexes: []SignedIndex{
			{
				Index: index.CompoundIndex(
					index.Field("name", 1),
					index.Field("age", -1),
				),
			},
		},
	}

	case1Intersect := case1Collection1.Intersect(case1Collection2)

	test.AssertTrue(t, case1Intersect != nil && len(*case1Intersect) == 0, "Case 1: Intersection is not empty")

	// test intersection with some differences
	case2Collection1 := SignedCollection{
		Metadata: SignedMetadata{
			Metadata: metadata.InitMetadata("users"),
		},
		Fields: []SignedField{
			{
				Field: field.StringField("name"),
			},
			{
				Field: field.Int32Field("age"),
			},
			{
				Field: field.ObjectField("body_info",
					field.DoubleField("height"),
					field.DoubleField("weight"),
				),
			},
		},
		Indexes: []SignedIndex{
			{
				Index: index.CompoundIndex(
					index.Field("name", 1),
					index.Field("age", -1),
				),
			},
		},
	}
	case2Collection2 := SignedCollection{
		Metadata: SignedMetadata{
			Metadata: metadata.InitMetadata("users"),
		},
		Fields: []SignedField{
			{
				Field: field.StringField("name"),
			},
			{
				Field: field.Int32Field("age"),
			},
			{
				Field: field.ObjectField("body_info",
					field.DoubleField("height"),
					field.Int64Field("weight"),
					field.TimestampField("updated_at"),
				),
			},
		},
		Indexes: []SignedIndex{
			{
				Index: index.CompoundIndex(
					index.Field("name", 1),
					index.Field("age", 1),
				),
			},
		},
	}

	case2Intersect := case2Collection1.Intersect(case2Collection2)
	test.AssertTrue(t, case2Intersect != nil && len(*case2Intersect) == 4, "Case 2: Intersection is not expected")
	for _, i := range *case2Intersect {
		// intersection flag should be true
		test.AssertTrue(t, i.IsIntersection, "Case 2: Intersection flag is not true")
		// the number of fields or indexes is either 0 or 1
		test.AssertTrue(t, len(i.Fields) == 1 || len(i.Indexes) == 1, "Case 2: Fields and Indexes are empty")
		if len(i.Fields) > 0 {
			test.AssertEqual(t, i.Fields[0].Spec().Type, field.TypeObject, "Case 2: field type is not object")
			test.AssertTrue(t, i.Fields[0].Spec().Object != nil && len(*i.Fields[0].Spec().Object) == 1, "Case 2: object item count is not 1")
			switch (*i.Fields[0].Spec().Object)[0].Name {
			case "weight":
				test.AssertEqual(t, i.Fields[0].Sign, SignConvert, "Case 2: weight field sign is not convert")
				test.AssertEqual(t, (*i.Fields[0].Spec().Object)[0].Type, field.TypeDouble, "Case 2: weight field conversion type is not Double")
				test.AssertEqual(t, (*i.Fields[0].convertFrom.Spec().Object)[0].Type, field.TypeInt64, "Case 2: ConversionFrom type is not Int64")
			case "updated_at":
				test.AssertEqual(t, i.Fields[0].Sign, SignMinus, "Case 2: updated_at field sign is not minus")
			default:
				t.Errorf("Case 2: Field %s is not valid", i.Metadata.Spec().Name)
			}
		} else if len(i.Indexes) > 0 {
			test.AssertEqual(t, i.Indexes[0].Spec().Type, index.TypeCompound, "Case 2: Index type is not valid")
			test.AssertEqual(t, len(i.Indexes[0].Index.Spec().Fields), 2, "Case 2: Index field count is not 2")
			test.AssertTrue(t, util.InListEq(i.Indexes[0].Sign, []EntitySign{
				SignPlus,
				SignMinus,
			}), "Case 2: Index entity sign is not valid")

			for _, indexField := range i.Indexes[0].Index.Spec().Fields {
				if !util.InListEq(indexField.Key, []string{"name", "age"}) {
					t.Errorf("Case 2: Index field %s is not valid", indexField.Key)
				}
			}
		}
	}
}

func TestSignedCollectionUnion(t *testing.T) {
	// TODO: implement this
}