package schema_interpreter

import (
	"fmt"

	"github.com/amirkode/go-mongr8/collection"
	"github.com/amirkode/go-mongr8/collection/field"
	"github.com/amirkode/go-mongr8/collection/index"
	"github.com/amirkode/go-mongr8/collection/metadata"
)

func (sas SubActionSchema) getMetadataDeclarationLiteral() string {
	res := fmt.Sprintf(`metadata.InitMetadata("%s")`, sas.Collection.Spec().Name)
		if sas.Collection.Spec().Options != nil {
		// check whether the collection is capped
		_, capped := (*sas.Collection.Spec().Options)[metadata.CollectionOptionCapped]
		if capped {
			res += fmt.Sprintf(`%s.Capped("%d")`, "\n", (*sas.Collection.Spec().Options)[metadata.CollectionOptionCappedSize])
		}
		// check whether the collection has expiration time
		_, hasExpiration := (*sas.Collection.Spec().Options)[metadata.CollectionOptionExpiredAfterSeconds]
		if hasExpiration {
			res += fmt.Sprintf(`%s.TTL("%d")`, "\n", (*sas.Collection.Spec().Options)[metadata.CollectionOptionExpiredAfterSeconds])
		}
	}

	return res
}

func (sas SubActionSchema) getFieldDeclarationLiteral(f collection.Field) string {
	var fieldLiteral func (f collection.Field) string
	fieldLiteral = func (f collection.Field) string {
		res := ""
		switch f.Spec().Type {
		case field.TypeString:
			res += fmt.Sprintf(`field.StringField("%s")`, f.Spec().Name)
		case field.TypeInt32:
			res += fmt.Sprintf(`field.Int32Field("%s")`, f.Spec().Name)
		case field.TypeInt64:
			res += fmt.Sprintf(`field.Int64Field("%s")`, f.Spec().Name)
		case field.TypeDouble:
			res += fmt.Sprintf(`field.DoubleField("%s")`, f.Spec().Name)
		case field.TypeBoolean:
			res += fmt.Sprintf(`field.BooleanField("%s")`, f.Spec().Name)
		case field.TypeArray:
			children := ""
			for _, child := range *f.Spec().ArrayFields {
				childField := field.FromFieldSpec(&child)
				children += fmt.Sprintf("%s,\n", fieldLiteral(childField))
			}
			res += fmt.Sprintf(`field.ArrayField("%s",%s%s)`, f.Spec().Name, "\n", children)
			// implement nested
		case field.TypeObject:
			children := ""
			for _, child := range *f.Spec().Object {
				childField := field.FromFieldSpec(&child)
				children += fmt.Sprintf("%s,\n", fieldLiteral(childField))
			}
			res += fmt.Sprintf(`field.ObjectField("%s",%s%s)`, f.Spec().Name, "\n", children)
			// implement nested
		case field.TypeTimestamp:
			res += fmt.Sprintf(`field.TimestampField("%s")`, f.Spec().Name)
		case field.TypeGeoJSONPoint:
			res += fmt.Sprintf(`field.GeoJSONPointField("%s")`, f.Spec().Name)
		case field.TypeGeoJSONLineString:
			res += fmt.Sprintf(`field.GeoJSONLineStringField("%s")`, f.Spec().Name)
		case field.TypeGeoJSONPolygonSingleRing:
			res += fmt.Sprintf(`field.GeoJSONPolygonSingleRingField("%s")`, f.Spec().Name)
		case field.TypeGeoJSONPolygonMultipleRing:
			res += fmt.Sprintf(`field.GeoJSONPolygonMultipleRingField("%s")`, f.Spec().Name)
		case field.TypeGeoJSONMultiPoint:
			res += fmt.Sprintf(`field.GeoJSONMultiPointField("%s")`, f.Spec().Name)
		case field.TypeGeoJSONMultiLineString:
			res += fmt.Sprintf(`field.GeoJSONMultiLineStringField("%s")`, f.Spec().Name)
		case field.TypeGeoJSONMultiPolygon:
			res += fmt.Sprintf(`field.GeoJSONMultiPolygonField("%s")`, f.Spec().Name)
		case field.TypeGeoJSONGeometryCollection:
			res += fmt.Sprintf(`field.GeoJSONGeometryCollectionField("%s")`, f.Spec().Name)
		case field.TypeLegacyCoordinateArray:
			res += fmt.Sprintf(`field.LegacyCoordinateArrayField("%s")`, f.Spec().Name)
		}

		return res
	}

	return fieldLiteral(f)
}

func (sas SubActionSchema) getIndexDeclarationLiteral(idx collection.Index) string {
	res := ""

	fieldsToMap := func() map[string]interface{} {
		fMap := map[string]interface{}{}
		for _, f := range idx.Spec().Fields {
			fMap[f.Key] = f.Value
		}

		return fMap
	}
	rulesToLiteral := func() string {
		if idx.Spec().Rules != nil {
			return AnyToLiteral(*idx.Spec().Rules)
		}

		return ""
	}

	switch idx.Spec().Type {
	case index.TypeSingleField:
		res += fmt.Sprintf(`index.SingleFieldIndex(index.Field("%s", %s))`, 
			idx.Spec().Fields[0].Key, 
			anyToLiteralString(idx.Spec().Fields[0].Value),
		)
	case index.TypeCompound:
		children := ""
		for _, f := range idx.Spec().Fields {
			children += fmt.Sprintf(`index.Field("%s", %s),%s`,
				f.Key,
				anyToLiteralString(f.Value),
				"\n",
			)
		}
		res += fmt.Sprintf(`index.CompoundIndex(%s%s)`, 
			"\n",
			children,
		)
	case index.TypeText:
		res += fmt.Sprintf(`index.TextIndex(index.Field("%s"))`,
			idx.Spec().Fields[0].Key,
		)
	case index.TypeGeopatial2dsphere:
		res += fmt.Sprintf(`index.Geospatial2dsphereIndex(index.Field("%s"))`,
			idx.Spec().Fields[0].Key,
		)
	case index.TypeUnique:
		res += fmt.Sprintf(`index.UniqueIndex(index.Field("%s", %s))`, 
			idx.Spec().Fields[0].Key, 
			anyToLiteralString(idx.Spec().Fields[0].Value),
		)
	case index.TypePartial:
		// then convert the map to literal as required in the arguments
		fArgs := AnyToLiteral(fieldsToMap())
		res += fmt.Sprintf(`index.PartialIndex(%s)`,
			fArgs,
		)
	case index.TypeCollation:
		res += fmt.Sprintf(`index.CollationIndex(index.Field("%s", %s))`, 
			idx.Spec().Fields[0].Key, 
			anyToLiteralString(idx.Spec().Fields[0].Value),
		)
	case index.TypeRaw:
		fArgs := AnyToLiteral(fieldsToMap())
		ruleArgs := rulesToLiteral()
		if ruleArgs == "" {
			res += fmt.Sprintf(`index.RawIndex(%s, nil)`,
				fArgs,
			)
		} else {
			res += fmt.Sprintf(`index.RawIndex(%s, &%s)`,
				fArgs,
				ruleArgs,
			)
		}
	}

	return res
}

func (sas SubActionSchema) GetLiteralInstance(prefix string, isArrayItem bool) string {
	res := ""
	if !isArrayItem {
		res += fmt.Sprintf("%sSubActionSchema", prefix)
	}

	res += "{\n"
	res += fmt.Sprintf("Collection: %s,\n", sas.getMetadataDeclarationLiteral())
	res += "Fields: []collection.Field{\n"
	// fill fields
	for _, f := range sas.Fields {
		res += fmt.Sprintf("%s,\n", sas.getFieldDeclarationLiteral(f))
	}

	res += "},\n"
	res += "Indexes: []collection.Index{\n"
	// fill indexes
	for _, idx := range sas.Indexes {
		res += fmt.Sprintf("%s,\n", sas.getIndexDeclarationLiteral(idx))
	}

	res += "},\n"
	// set convertFrom if exists
	if sas.FieldConvertFrom != nil {
		res += fmt.Sprintf("FieldConvertFrom: field.GetTypePointer(field.%s),\n", sas.FieldConvertFrom.ToString())
	}

	res += "}"

	return res
}
