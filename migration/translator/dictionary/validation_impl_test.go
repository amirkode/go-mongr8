/*
Copyright (c) 2023 the go-mongr8 Authors and Contributors
[@see Authors file]

Licensed under the MIT License
(https://opensource.org/licenses/MIT)
*/
package dictionary

import (
	"strings"
	"testing"

	"github.com/amirkode/go-mongr8/internal/test"
	"github.com/amirkode/go-mongr8/collection"
	"github.com/amirkode/go-mongr8/collection/field"
	"github.com/amirkode/go-mongr8/collection/index"
	"github.com/amirkode/go-mongr8/collection/metadata"
	"github.com/amirkode/go-mongr8/migration/common"
)

func TestValidateCollections(t *testing.T) {
	// Case 1: invalid collection name with mongr8 migration history collection
	case1Err := validateCollections([]collection.Collection{
		collection.NewCollection(metadata.InitMetadata(common.MigrationHistoryCollection), []collection.Field{}, []collection.Index{}),
	})

	test.AssertTrue(t, case1Err != nil && strings.Contains(case1Err.Error(), common.MigrationHistoryCollection), "Case 1: Unxpected error")

	// Case 2: duplicate collection names
	case2Err := validateCollections([]collection.Collection{
		collection.NewCollection(metadata.InitMetadata("collection1"), []collection.Field{}, []collection.Index{}),
		collection.NewCollection(metadata.InitMetadata("collection1"), []collection.Field{}, []collection.Index{}),
	})

	test.AssertTrue(t, case2Err != nil && strings.Contains(case2Err.Error(), "Duplicate collection"), "Case 1: Unxpected error")

	// Case 1: collections are valid
	case3Err := validateCollections([]collection.Collection{
		collection.NewCollection(metadata.InitMetadata("collection1"), []collection.Field{}, []collection.Index{}),
		collection.NewCollection(metadata.InitMetadata("collection2"), []collection.Field{}, []collection.Index{}),
	})

	test.AssertTrue(t, case3Err == nil, "Case 1: Unxpected error")
}

func TestValidateID(t *testing.T) {
	// Case 1: Unallowed type object
	case1Err := validateID("collection_name", []collection.Field{field.ObjectField("_id")})
	
	test.AssertTrue(t, case1Err != nil && strings.Contains(case1Err.Error(), "ID field type invalid"), "Case 1: Unexpected error")
	
	// Case 2: Unallowed type geo json point
	case2Err := validateID("collection_name", []collection.Field{field.GeoJSONPointField("_id")})
	
	test.AssertTrue(t, case2Err != nil && strings.Contains(case2Err.Error(), "ID field type invalid"), "Case 2: Unexpected error")

	// Case 3: Allowed field string
	case3Err := validateID("collection_name", []collection.Field{field.StringField("_id")})
	
	test.AssertTrue(t, case3Err == nil, "Case 3: Unexpected error")

	// Case 4: Allowed field double
	case4Err := validateID("collection_name", []collection.Field{field.DoubleField("_id")})
	
	test.AssertTrue(t, case4Err == nil, "Case 4: Unexpected error")
}

func TestValidateFieldDuplication(t *testing.T) {
	// Case 1: duplicate two string fields
	case1Err := validateFieldDuplication("collection_name", []collection.Field{
		field.StringField("name"),
		field.StringField("name"),
	})

	test.AssertTrue(t, case1Err != nil && strings.Contains(case1Err.Error(), "duplicate field found"), "Case 1: Unexpected error")

	// Case 2: duplicate two array fields
	case2Err := validateFieldDuplication("collection_name", []collection.Field{
		field.ArrayField("name", field.StringField("")),
		field.ArrayField("name", field.StringField("")),
	})

	test.AssertTrue(t, case2Err != nil && strings.Contains(case2Err.Error(), "duplicate field found"), "Case 2: Unexpected error")

	// Case 3: duplicate two array of object fields
	case3Err := validateFieldDuplication("collection_name", []collection.Field{
		field.ArrayField("name", field.ObjectField("", 
			field.StringField("first_name"),
			field.StringField("first_name"),
		)),
	})

	test.AssertTrue(t, case3Err != nil && strings.Contains(case3Err.Error(), "duplicate field found"), "Case 3: Unexpected error")

	// Case 4: two array of object fields
	case4Err := validateFieldDuplication("collection_name", []collection.Field{
		field.ArrayField("name", field.ObjectField("", 
			field.StringField("first_name"),
			field.StringField("second_name"),
		)),
	})

	test.AssertTrue(t, case4Err == nil, "Case 4: Unexpected error")

	// Case 5: duplicate two string fields
	case5Err := validateFieldDuplication("collection_name", []collection.Field{
		field.StringField("name1"),
		field.StringField("name2"),
	})

	test.AssertTrue(t, case5Err == nil, "Case 5: Unexpected error")
}

func TestValidateIndividualField(t *testing.T) {
	// Case 1: Field with exceeded max name length
	case1Err := validateIndividualField("collection_name", "", field.StringField("gvpzqwjlnbpptaiejcrpzzwjeqsoyxawhaprxnlbtbpiwzvrwvuqljajqpjxkjsrraxligwopgvhkzfkfajlrlefoujscbbfdemirnmbviolxpucrccrisiwcyloxuhtx"), false)

	test.AssertTrue(t, case1Err != nil && strings.Contains(case1Err.Error(), "field name more than 128 characters len"), "Case 1: Unexpected error")

	// Case 2: String field with empty name
	case2Err := validateIndividualField("collection_name", "", field.StringField(""), false)

	test.AssertTrue(t, case2Err != nil && strings.Contains(case2Err.Error(), "Field name must not be empty"), "Case 2: Unexpected error")

	// Case 3: Array of object field with exceeded max name length
	case3Err := validateIndividualField("collection_name", "", field.ArrayField("arr", 
		field.ObjectField("", field.StringField("gvpzqwjlnbpptaiejcrpzzwjeqsoyxawhaprxnlbtbpiwzvrwvuqljajqpjxkjsrraxligwopgvhkzfkfajlrlefoujscbbfdemirnmbviolxpucrccrisiwcyloxuhtx")),
	), false)

	test.AssertTrue(t, case3Err != nil && strings.Contains(case3Err.Error(), "field name more than 128 characters len"), "Case 3: Unexpected error")

	// Case 4: Array of object field with empty name
	case4Err := validateIndividualField("collection_name", "", field.ArrayField("arr", 
		field.ObjectField("", field.StringField("")),
	), false)

	test.AssertTrue(t, case4Err != nil && strings.Contains(case4Err.Error(), "Field name must not be empty"), "Case 4: Unexpected error")

	// Case 5: Array item field with empty name
	case5Err := validateIndividualField("collection_name", "", field.StringField(""), true)

	test.AssertTrue(t, case5Err == nil, "Case 5: Unexpected error")

	// Case 6: Array of object field
	case6Err := validateIndividualField("collection_name", "", field.ArrayField("arr", 
		field.ObjectField("", field.StringField("name")),
	), false)

	test.AssertTrue(t, case6Err == nil, "Case 6: Unexpected error")
}

func TestValidateFields(t *testing.T) {
	// Case 1: duplicate two string fields
	case1Err := validateFields("collection_name", []collection.Field{
		field.StringField("name"),
		field.StringField("name"),
	})

	test.AssertTrue(t, case1Err != nil && strings.Contains(case1Err.Error(), "duplicate field found"), "Case 1: Unexpected error")
	
	// Case 2: String field with empty name
	case2Err := validateFields("collection_name", []collection.Field{field.StringField("")})

	test.AssertTrue(t, case2Err != nil && strings.Contains(case2Err.Error(), "Field name must not be empty"), "Case 2: Unexpected error")

	// Case 3: Array of object field
	case3Err := validateFields("collection_name", []collection.Field{field.ArrayField("arr", 
		field.ObjectField("", field.StringField("name")),
	)})

	test.AssertTrue(t, case3Err == nil, "Case 6: Unexpected error")
}

func TestValidateIndexDuplication(t *testing.T) {
	// Case 1: duplicate single field indexes
	case1Err := validateIndexDuplication("collection_name", []collection.Index{
		index.SingleFieldIndex(index.Field("name", 1)),
		index.SingleFieldIndex(index.Field("name", 1)),
	})

	test.AssertTrue(t, case1Err != nil && strings.Contains(case1Err.Error(), "duplicate index found"), "Case 1: Unexpected error")

	// Case 2: duplicate compound field indexes
	case2Err := validateIndexDuplication("collection_name", []collection.Index{
		index.CompoundIndex(
			index.Field("name", 1),
			index.Field("age", -1),
		),
		index.CompoundIndex(
			index.Field("name", 1),
			index.Field("age", -1),
		),
	})

	test.AssertTrue(t, case2Err != nil && strings.Contains(case1Err.Error(), "duplicate index found"), "Case 2: Unexpected error")

	// Case 3: duplicate single field indexes with unique option
	case3Err := validateIndexDuplication("collection_name", []collection.Index{
		index.SingleFieldIndex(index.Field("name", 1)).AsUnique(),
		index.SingleFieldIndex(index.Field("name", 1)).AsUnique(),
	})

	test.AssertTrue(t, case3Err != nil && strings.Contains(case1Err.Error(), "duplicate index found"), "Case 3: Unexpected error")

	// Case 4: single field indexes with different options
	case4Err := validateIndexDuplication("collection_name", []collection.Index{
		index.SingleFieldIndex(index.Field("name", 1)).SetCollation(map[string]interface{}{"locale": "en_US"}),
		index.SingleFieldIndex(index.Field("name", 1)).AsUnique(),
	})

	test.AssertTrue(t, case4Err == nil, "Case 5: Unexpected error")

	// Case 5: single field indexes with different names
	case5Err := validateIndexDuplication("collection_name", []collection.Index{
		index.SingleFieldIndex(index.Field("name", 1)),
		index.SingleFieldIndex(index.Field("age", 1)),
	})

	test.AssertTrue(t, case5Err == nil, "Case 5: Unexpected error")
}

func TestValidateIndexWithFields(t *testing.T) {
	// Case 1: index with empty fields
	case1Err := validateIndexWithFields("collection_name", []collection.Field{}, collection.IndexFromSpec(&index.Spec{}))

	test.AssertTrue(t, case1Err != nil && strings.Contains(case1Err.Error(), "Index Fields cannot be empty"), "Case 1: Unexpected error")

	// Case 2: index field does not present in collection field - String field
	case2Err := validateIndexWithFields("collection_name", []collection.Field{field.StringField("name")}, index.SingleFieldIndex(index.Field("address", 1)))

	test.AssertTrue(t, case2Err != nil && strings.Contains(case2Err.Error(), "index key is invalid"), "Case 2: Unexpected error")

	// Case 3: index field does not present in collection field - Array of Object
	case3Err := validateIndexWithFields("collection_name", 
		[]collection.Field{field.ArrayField("values", field.ObjectField("", field.Int32Field("score")))},
		index.SingleFieldIndex(index.Field("values.name", 1)),
	)

	test.AssertTrue(t, case3Err != nil && strings.Contains(case3Err.Error(), "index key is invalid"), "Case 3: Unexpected error")

	// Case 4: index field does not present in collection field - Object of Object
	case4Err := validateIndexWithFields("collection_name", 
		[]collection.Field{field.ObjectField("field", field.ObjectField("child_field", field.Int32Field("child_child_field")))},
		index.SingleFieldIndex(index.Field("field.child_field.name", 1)),
	)

	test.AssertTrue(t, case4Err != nil && strings.Contains(case4Err.Error(), "index key is invalid"), "Case 4: Unexpected error")

	// Case 5: index field presents in collection field - String field
	case5Err := validateIndexWithFields("collection_name", []collection.Field{field.StringField("name")}, index.SingleFieldIndex(index.Field("name", 1)))

	test.AssertTrue(t, case5Err == nil, "Case 5: Unexpected error")

	// Case 6: index field presents in collection field - Array of Object
	case6Err := validateIndexWithFields("collection_name", 
		[]collection.Field{field.ArrayField("values", field.ObjectField("", field.Int32Field("score")))},
		index.SingleFieldIndex(index.Field("values.score", 1)),
	)

	test.AssertTrue(t, case6Err != nil, "Case 6: Unexpected error")

	// Case 7: index field presents in collection field - Object of Object
	case7Err := validateIndexWithFields("collection_name", 
		[]collection.Field{field.ObjectField("field", field.ObjectField("child_field", field.Int32Field("child_child_field")))},
		index.SingleFieldIndex(index.Field("field.child_field.child_child_field", 1)),
	)

	test.AssertTrue(t, case7Err == nil, "Case 7: Unexpected error")

	// Case 8: partial expression key does not present in field
	case8Err := validateIndexWithFields("collection_name", 
		[]collection.Field{field.StringField("name"), field.StringField("address")}, 
		index.SingleFieldIndex(index.Field("address", 1)).SetPartialExpression(map[string]interface{}{"wrong_field": "some value"}),
	)

	test.AssertTrue(t, case8Err != nil && strings.Contains(case8Err.Error(), "Partial filter key is invalid"), "Case 8: Unexpected error")	

	// Case 9: partial expression key presents in field
	case9Err := validateIndexWithFields("collection_name", 
		[]collection.Field{field.StringField("name"), field.StringField("address")}, 
		index.SingleFieldIndex(index.Field("address", 1)).SetPartialExpression(map[string]interface{}{"name": "some value"}),
	)

	test.AssertTrue(t, case9Err == nil, "Case 9: Unexpected error")	

	// Case 10: TTL option with no timestamp field
	case10Err := validateIndexWithFields("collection_name", 
		[]collection.Field{field.StringField("name"), field.StringField("address")}, 
		index.SingleFieldIndex(index.Field("address", 1)).SetTTL(3600),
	)

	test.AssertTrue(t, case10Err != nil && strings.Contains(case10Err.Error(), "Timestamp field must exist in TTL index"), "Case 10: Unexpected error")	

	// Case 11: TTL option with timestamp field
	case11Err := validateIndexWithFields("collection_name", 
		[]collection.Field{field.StringField("name"), field.TimestampField("updated_at")}, 
		index.CompoundIndex(index.Field("name", 1), index.Field("updated_at", 1)).SetTTL(3600),
	)

	test.AssertTrue(t, case11Err == nil, "Case 11: Unexpected error")	
}

func TestValidateIndexes(t *testing.T) {
	// TODO: implement this
}