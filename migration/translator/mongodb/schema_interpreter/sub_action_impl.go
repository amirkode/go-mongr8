package schema_interpreter

import (
	"internal/convert"
	dt "internal/data_type"
	"internal/util"

	"github.com/amirkode/go-mongr8/collection"
	"github.com/amirkode/go-mongr8/collection/index"
	"github.com/amirkode/go-mongr8/migration/translator/dictionary"

	"go.mongodb.org/mongo-driver/bson"
)

func getIndexesMap(indexes []index.IndexField) map[string]interface{} {
	res := map[string]interface{}{}
	for _, index := range indexes {
		res[index.Key] = index.Value
	}

	return res
}

func indexesToBsonD(indexes []index.IndexField) bson.D {
	return convert.MapToBsonD(getIndexesMap(indexes))
}

func indexesToBsonM(indexes []index.IndexField) bson.M {
	return convert.MapToBsonM(getIndexesMap(indexes))
}

func getFieldsMap(fields []collection.Field) map[string]interface{} {
	res := map[string]interface{}{}
	for _, field := range fields {
		// Translate
		translated := dictionary.GetTranslatedField(field)
		res[field.Spec().Name] = translated.GetObject()
	}

	return ConvertValueTypeToRealType(res).(map[string]interface{})
}

func fieldsToBsonD(fields []collection.Field) bson.D {
	return convert.MapToBsonD(getFieldsMap(fields))
}

func fieldsToBsonM(fields []collection.Field) bson.M {
	return convert.MapToBsonM(getFieldsMap(fields))
}

func (sa SubAction) GetIndexesBsonD() []dt.Pair[bson.D, bson.D] {
	sa.validate()
	res := []dt.Pair[bson.D, bson.D]{}
	for _, index := range sa.ActionSchema.Indexes {
		indexes := indexesToBsonD(index.Spec().Fields)
		if index.Spec().Rules != nil {
			res = append(res, dt.NewPair(indexes, convert.MapToBsonD(*index.Spec().Rules)))
		} else {
			res = append(res, dt.NewPair(indexes, bson.D{}))
		}
	}

	return res
}

func (sa SubAction) GetIndexesBsonM() []dt.Pair[bson.M, bson.M] {
	sa.validate()
	res := []dt.Pair[bson.M, bson.M]{}
	for _, index := range sa.ActionSchema.Indexes {
		indexes := indexesToBsonM(index.Spec().Fields)
		if index.Spec().Rules != nil {
			res = append(res, dt.NewPair(indexes, convert.MapToBsonM(*index.Spec().Rules)))
		} else {
			res = append(res, dt.NewPair(indexes, bson.M{}))
		}
	}

	return res
}

func (sa SubAction) GetFieldsBsonD() bson.D {
	sa.validate()
	return fieldsToBsonD(sa.ActionSchema.Fields)
}

func (sa SubAction) GetFieldsBsonM() bson.M {
	sa.validate()
	return fieldsToBsonM(sa.ActionSchema.Fields)
}

func (sa SubAction) IsUp() bool {
	return util.InListEq(sa.Type, []SubActionType{
		SubActionTypeCreateCollection,
		SubActionTypeCreateField,
		SubActionTypeCreateIndex,
		SubActionTypeConvertField,
	})
}

func SubActionCreateCollection(schema SubActionSchema) *SubAction {
	return &SubAction{
		Type:         SubActionTypeCreateCollection,
		ActionSchema: schema,
		validate: func() {
			if len(schema.Fields) > 0 {
				panic("At least a field declared on collection creation")
			}
		},
	}
}

func SubActionCreateIndex(schema SubActionSchema) *SubAction {
	return &SubAction{
		Type:         SubActionTypeCreateIndex,
		ActionSchema: schema,
		validate: func() {
			if len(schema.Indexes) != 1 {
				panic("At least an index declared for index creation")
			}
		},
	}
}

func SubActionCreateField(schema SubActionSchema) *SubAction {
	return &SubAction{
		Type:         SubActionTypeCreateField,
		ActionSchema: schema,
		validate: func() {
			if len(schema.Fields) != 1 {
				panic("At least a field declared for field creation")
			}
		},
	}
}

func SubActionConvertField(schema SubActionSchema) *SubAction {
	return &SubAction{
		Type:         SubActionTypeConvertField,
		ActionSchema: schema,
		validate: func() {
			if len(schema.Indexes) != 1 {
				panic("At least a field declared for conversion")
			}

			if schema.FieldConvertFrom == nil {
				panic("FieldConvertFrom must not be nil for conversion")
			}
		},
	}
}

func SubActionDropCollection(schema SubActionSchema) *SubAction {
	return &SubAction{
		Type:         SubActionTypeDropCollection,
		ActionSchema: schema,
		validate: func() {
			// nothing to validate
		},
	}
}

func SubActionDropIndex(schema SubActionSchema) *SubAction {
	return &SubAction{
		Type:         SubActionTypeDropIndex,
		ActionSchema: schema,
		validate: func() {
			// expecting exactly 1 index to drop by schema
			if len(schema.Indexes) != 1 {
				panic("At least an index declared for dropping index")
			}
		},
	}
}

func SubActionDropField(schema SubActionSchema) *SubAction {
	return &SubAction{
		Type:         SubActionTypeDropIndex,
		ActionSchema: schema,
		validate: func() {
			// expecting exactly 1 field to drop by field name
			if len(schema.Fields) != 1 {
				panic("At least an index declared for dropping field")
			}
		},
	}
}

/*
var ctx context.Context
var db *mongo.Database

// init subaction in MongoDB manner
// we assume all pre-validations have passed

func NewCreateCollectionSubAction(collectionName string, opts map[metadata.CollectionOption]interface{}) SubAction {
	stmt := fmt.Sprintf("%s := options.CreateCollectionOptions{}\n", VarNameCreateOptions)
	for key, value := range opts {
		if key == metadata.CollectionOptionCapped {
			stmt += fmt.Sprintf("%s.SetCapped(%s)\n", VarNameCreateOptions, anyToLiteralString(value))
		} else if key == metadata.CollectionOptionSize {
			stmt += fmt.Sprintf("%s.SetSizeInBytes(%s)\n", VarNameCreateOptions, anyToLiteralString(value))
		} else if key == metadata.CollectionOptionExpiredAfterSeconds {
			stmt += fmt.Sprintf("%s.SetExpireAfterSeconds(%s)\n", VarNameCreateOptions, anyToLiteralString(value))
		}
	}

	stmt += fmt.Sprintf(`if %s := %s.CreateCollection(%s, "%s", &%s); %s != nil { return %s }\n`,
		VarNameError,
		VarNameDatabase,
		VarNameContext,
		collectionName,
		VarNameCreateOptions,
		VarNameError,
		VarNameError,
	)

	return SubAction{
		GetStatement: func() string {
			return stmt
		},
		Type: SubActionCreateCollection,
	}
}

func NewInsertSingleSubAction(initCollection bool, collectionName string, payload map[string]interface{}) SubAction {
	stmt := ""
	// if it requires to initialize collection instance
	if initCollection {
		stmt += fmt.Sprintf(`%s := %s.Collection("%s")`, VarNameCollection, VarNameDatabase, collectionName) + "\n"
	}

	stmt += fmt.Sprintf(`count, _ := %s.CountDocuments(%s, bson.M{"_id": %s})`, VarNameCollection, VarNameContext, anyToLiteralString(payload["_id"])) + "\n"
	stmt += fmt.Sprintf(`if count > 0 { return fmt.Errorf("Collection already exists, cannot insert initial data") }`) + "\n"
	stmt += fmt.Sprintf(`if _, %s := %s.InsertOne(%s, %s); %s != nil { return %s }`,
		VarNameError,
		VarNameCollection,
		VarNameContext,
		toLiteralStringBsonMap(payload),
		VarNameError,
		VarNameError,
	)

	return SubAction{
		GetStatement: func() string {
			return stmt
		},
		Type: SubActionInsertOne,
	}
}
*/
