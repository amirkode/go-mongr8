package sync_strategy

// this contains mechanism to sync schema across all sources
// the main ideas of the synchronization are:
// - union of list
// - difference of intersection

import (
	"fmt"
	dt "internal/data_type"
	"internal/util"

	"github.com/amirkode/go-mongr8/collection"
	"github.com/amirkode/go-mongr8/collection/field"
)

type EntitySign int8

const (
	SignPlus  EntitySign = 1
	SignMinus EntitySign = -1
	// this additional sign means a conversion needed from previous entity
	// for now, the usecase is only for field entity
	// make it generic sign to cover future usecase in other entities
	SignConvert EntitySign = 0
)

type (
	// TODO: implement these function with pointer receiver
	operator[T any] interface {
		// intersection of two entities resulting modifications
		// TODO: rename [other] with [origin] and rename [this] with [incoming]
		Intersect(other T) *[]T
		// set entity sign (action direction)
		SetSign(sign EntitySign) T
		// set flag denoting whether the entity is formed by intersection
		SetIntersectionFlag(flag bool) T
		// if the current entity sign is convert
		// this function returns the original entity
		// that will be converted from
		// for now, it's probably has not yet used
		ConvertFrom() *T
		// a unique key required from Union and Intersection
		Key() string
	}

	SignedField struct {
		operator[SignedField]
		collection.Field
		convertFrom *SignedField
		Sign        EntitySign
	}

	SignedIndex struct {
		operator[SignedIndex]
		collection.Index
		Sign EntitySign
	}

	SignedMetadata struct {
		operator[SignedMetadata]
		collection.Metadata
		Sign EntitySign
	}

	SignedCollection struct {
		operator[SignedCollection]
		Metadata       SignedMetadata
		Fields         []SignedField
		Indexes        []SignedIndex
		Sign           EntitySign
		IsIntersection bool
	}
)

func (f SignedField) Intersect(other SignedField) *[]SignedField {
	if f.Key() != other.Key() {
		return nil
	}

	res := []SignedField{}
	// check recursively with DFS, e.g:
	// we want to check this
	// this:
	// {"field": {"field1": {"field2": "something"}}}
	// other:
	// {"field": {"field1": {"field3": "something"}}}
	// as we can see, there's a difference in field1's child
	// perform DFS to find those differences and
	// add up relevant signed field to "res"

	// restores field path returning a SignedField
	restorePath := func(path []dt.Pair[string, field.FieldType]) (SignedField, *collection.Field) {
		if len(path) < 1 {
			panic("Path length must be more than 1")
		}

		// init field
		headField := collection.FieldFromType(path[0].First, path[0].Second)
		prevField := &headField
		for i := 1; i < len(path); i++ {
			if !util.InList((*prevField).Spec().Type, []field.FieldType{
				field.TypeArray,
				field.TypeObject,
			}) {
				panic("Cannot add more child, previous field type is not Array or Object")
			}

			currField := collection.FieldFromType(path[i].First, path[i].Second)
			if path[i].Second == field.TypeArray {
				// set current field as previous field's array field
				(*prevField).Spec().ArrayFields = &[]field.Spec{*currField.Spec()}
			} else if path[i].Second == field.TypeObject {
				// set current field as previous field's object field
				(*prevField).Spec().Object = &[]field.Spec{*currField.Spec()}
			}

			prevField = &currField
		}

		signedField := SignedField{
			Field: headField,
			Sign:  SignPlus,
		}

		return signedField, prevField
	}

	// dfs function
	dfs := func(path []dt.Pair[string, field.FieldType], this SignedField, other SignedField) {
		// check type of both fields
		// assuming if one of the field is not an array or an object
		// we don't need to check the children (array/object)
		if this.Spec().Type != other.Spec().Type {
			// decide proper type conversion (Supported, Unsupported, or Undefined)
			if this.Spec().Type == field.TypeString ||
				(this.Spec().Type.IsNumeric() && other.Spec().Type.IsNumeric()) {
				// by default any type to string must be supported
				// for numeric to numeric conversion, there's an edge case
				// please see note on sync.go
				convertFrom, _ := restorePath(append(path, dt.NewPair(this.Spec().Name, other.Spec().Type)))
				convert, _ := restorePath(append(path, dt.NewPair(this.Spec().Name, this.Spec().Type)))
				convert.Sign = SignConvert
				convert.convertFrom = &convertFrom
				res = append(res, convert)
			} else if other.Spec().Type == field.TypeString {
				// string to any type conversion
				// this must be undefined conversion type
				// by default just perform drop and add

				// TODO: add condition for user-defined Force Conversion option

				// add plus action for "this"
				plus, _ := restorePath(append(path, dt.NewPair(this.Spec().Name, this.Spec().Type)))
				plus.Sign = SignPlus
				res = append(res, plus)
				// add minus action for "other"
				minus, _ := restorePath(append(path, dt.NewPair(this.Spec().Name, other.Spec().Type)))
				minus.Sign = SignMinus
				res = append(res, minus)
			} else {
				// this should be unsupported type conversion
				// add plus action for "this"
				plus, _ := restorePath(append(path, dt.NewPair(this.Spec().Name, this.Spec().Type)))
				plus.Sign = SignPlus
				res = append(res, plus)
				// add minus action for "other"
				minus, _ := restorePath(append(path, dt.NewPair(this.Spec().Name, other.Spec().Type)))
				minus.Sign = SignMinus
				res = append(res, minus)
			}

			return
		}

		// append path
		path = append(path, dt.NewPair(this.Spec().Name, this.Spec().Type))

		// - if the both types are same
		// continue find deeper on array type or object type
		// actually, we don't really call this DFS function to go deeper
		// instead, we call "Union" function that eventually will call another "Intersect" function
		// the behaviour might change in the future. So, we keep declaring this DFS function
		// - otherwise, we perform nothing
		switch this.Spec().Type {
		case field.TypeArray:
			if this.Spec().ArrayFields == nil || len(*this.Spec().ArrayFields) == 0 {
				panic("Array fields should not be nil or empty")
			}

			// if array fields length more than 1 (multiple types)
			// the number of items is 1, it's most likely fixed
			// TODO: in the future, we treat it more advance

			// just perform union
			// init this array fields
			thisFields := []SignedField{}
			for _, arrField := range *this.Spec().ArrayFields {
				// copy to another var to get a new reference
				currField := arrField
				thisFields = append(thisFields, SignedField{
					Field: field.FromFieldSpec(&currField),
					Sign:  SignPlus,
				})
			}
			// init other array fields
			otherArrFields := []SignedField{}
			if other.Spec().ArrayFields != nil {
				for _, arrField := range *other.Spec().ArrayFields {
					// copy to another var to get a new reference
					currField := arrField
					otherArrFields = append(otherArrFields, SignedField{
						Field: field.FromFieldSpec(&currField),
						Sign:  SignPlus,
					})
				}
			}

			union := Union(thisFields, otherArrFields)
			for _, u := range union {
				// create new instance of signed field each child
				curr, prevField := restorePath(path)
				// the sign of curr based on nested child u
				curr.Sign = u.Sign
				// check whether current field sign is convert
				// then, we need set the convertFrom in the current
				if u.Sign == SignConvert {
					convertFrom, cfPrevField := restorePath(path)
					cfCurrSpec := *u.convertFrom.Spec()
					(*cfPrevField).Spec().ArrayFields = &[]field.Spec{cfCurrSpec}
					// set curr.convertFrom
					convertFrom.Sign = SignConvert
					curr.convertFrom = &convertFrom
				}
				// join u.Collection as curr.Collection.ArrayFields
				currSpec := *u.Spec()
				(*prevField).Spec().ArrayFields = &[]field.Spec{currSpec}
				// eventually, add curr to res
				res = append(res, curr)
			}
		case field.TypeObject:
			// for object, we do the exact same operation as the array field
			// we separate all the (similar) parts because we might
			// encounter some differences in the future
			if this.Spec().Object == nil || len(*this.Spec().Object) == 0 {
				panic("Object fields should not be nil or empty")
			}

			// just perform union
			// init this array fields
			thisFields := []SignedField{}
			for _, objField := range *this.Spec().Object {
				// copy to another var to get a new reference
				currField := objField
				thisFields = append(thisFields, SignedField{
					Field: field.FromFieldSpec(&currField),
					Sign:  SignPlus,
				})
			}
			// init other array fields
			otherObjFields := []SignedField{}
			if other.Spec().Object != nil {
				for _, objField := range *other.Spec().Object {
					// copy to another var to get a new reference
					currField := objField
					otherObjFields = append(otherObjFields, SignedField{
						Field: field.FromFieldSpec(&currField),
						Sign:  SignPlus,
					})
				}
			}

			union := Union(thisFields, otherObjFields)
			for _, u := range union {
				// create new instance of signed field each child
				curr, prevField := restorePath(path)
				// the sign of curr based on nested child u
				curr.Sign = u.Sign
				// check whether current field sign is convert
				// then, we need set the convertFrom in the current
				if u.Sign == SignConvert {
					convertFrom, cfPrevField := restorePath(path)
					cfCurrSpec := *u.convertFrom.Spec()
					(*cfPrevField).Spec().Object = &[]field.Spec{cfCurrSpec}
					// set curr.convertFrom
					convertFrom.Sign = SignConvert
					curr.convertFrom = &convertFrom
				}
				// join u.Collection as curr.Collection.Object
				currSpec := *u.Spec()
				(*prevField).Spec().Object = &[]field.Spec{currSpec}
				// eventually, add curr to res
				res = append(res, curr)
			}
		}
	}

	// perform DFS
	dfs([]dt.Pair[string, field.FieldType]{}, f, other)

	// this will also return an emtpy slice
	// if there's no schema difference in deepest level
	return &res
}

func (f SignedField) ConvertFrom() *SignedField {
	return f.convertFrom
}

// get deepest type of a signed field
// assuming there's only a single way
func (f SignedField) FieldDeepestType() field.FieldType {
	var dfs func(_field collection.Field) field.FieldType
	dfs = func(_field collection.Field) field.FieldType {
		switch _field.Spec().Type {
		case field.TypeArray:
			return dfs(field.FromFieldSpec(&(*_field.Spec().ArrayFields)[0]))
		case field.TypeObject:
			return dfs(field.FromFieldSpec(&(*_field.Spec().Object)[0]))
		}

		return _field.Spec().Type
	}

	return dfs(f.Field)
}

// set deepest type of a signed field
// assuming there's only a single way
func (f SignedField) SetFieldDeepestType(toType field.FieldType) {
	var dfs func(_field *field.Spec)
	dfs = func(_field *field.Spec) {
		switch (*_field).Type {
		case field.TypeArray:
			dfs(&(*_field.ArrayFields)[0])
			return
		case field.TypeObject:
			dfs(&(*_field.Object)[0])
			return
		}

		_field.Type = toType
	}

	dfs(f.Spec())
}

func (f SignedField) SetSign(sign EntitySign) SignedField {
	f.Sign = sign
	return f
}

func (f SignedField) SetIntersectionFlag(flag bool) SignedField {
	return f
}

func (f SignedField) Key() string {
	// field name as key
	return f.Spec().Name
}

// this returns the difference between two index
// if both share the same name and type
// cases that might happen:
//  1. index1 -> {"name": 1}
//     index2 -> {"name": "text"}
//     returns list of [dropping index1,  add index2]
//  2. index1 -> {"name": 1}
func (f SignedIndex) Intersect(other SignedIndex) *[]SignedIndex {
	// intersect only happens when keys of both indexes are same
	// but, there's at at least 1 difference at their's nested level
	// this case won't happen since, each kay is constructed
	// by whole structure of the index
	// so, for now, all intersects in the "Union" process
	// will be ignored
	return nil
}

func (f SignedIndex) ConvertFrom() *SignedIndex { return nil }

func (f SignedIndex) SetSign(sign EntitySign) SignedIndex {
	f.Sign = sign
	return f
}

func (f SignedIndex) SetIntersectionFlag(flag bool) SignedIndex {
	return f
}

func (f SignedIndex) Key() string {
	return f.Spec().GetKey()
}

func (f SignedMetadata) Intersect(other SignedMetadata) *[]SignedMetadata {
	return nil
}

func (f SignedMetadata) ConvertFrom() *SignedMetadata { return nil }

func (f SignedMetadata) SetSign(sign EntitySign) SignedMetadata {
	f.Sign = sign
	return f
}

func (f SignedMetadata) SetIntersectionFlag(flag bool) SignedMetadata {
	return f
}

func (f SignedMetadata) Key() string {
	// combination of name and all options
	key := string(f.Spec().Name)
	key += string(f.Spec().Type)
	if f.Spec().Options != nil {
		key += fmt.Sprintf("%v", *f.Spec().Options)
	}

	return key
}

func (f SignedCollection) Intersect(other SignedCollection) *[]SignedCollection {
	if f.Key() != other.Key() {
		return nil
	}

	// panic if there's any metadata difference
	// once a collection is declared, it's expected
	// not to be modfied on any option
	if f.Metadata.Key() != other.Metadata.Key() {
		panic(fmt.Sprintf("Metadata should be modified on collection %s", f.Metadata.Spec().Name))
	}

	res := []SignedCollection{}
	// get union of fields and push as individual SignedCollection(s)
	signedFields := Union(f.Fields, other.Fields)
	for _, signedField := range signedFields {
		curr := signedField
		res = append(res, SignedCollection{
			Metadata:       f.Metadata,
			Fields:         []SignedField{curr},
			Sign:           curr.Sign,
			IsIntersection: true,
		})
	}
	// get union of indexes and push as individual SignedCollection(s)
	signedIndexes := Union(f.Indexes, other.Indexes)
	for _, signedIndex := range signedIndexes {
		curr := signedIndex
		res = append(res, SignedCollection{
			Metadata:       f.Metadata,
			Indexes:        []SignedIndex{curr},
			Sign:           curr.Sign,
			IsIntersection: true,
		})
	}

	return &res
}

func (f SignedCollection) SetSign(sign EntitySign) SignedCollection {
	f.Sign = sign
	return f
}

func (f SignedCollection) SetIntersectionFlag(flag bool) SignedCollection {
	f.IsIntersection = flag
	return f
}

func (f SignedCollection) Key() string {
	return f.Metadata.Spec().Name
}

func (f SignedCollection) GetFields() []collection.Field {
	res := make([]collection.Field, len(f.Fields))
	for index, signedField := range f.Fields {
		res[index] = signedField.Field
	}

	return res
}

func (f SignedCollection) GetIndexes() []collection.Index {
	res := make([]collection.Index, len(f.Indexes))
	for index, signedIndex := range f.Indexes {
		res[index] = signedIndex.Index
	}

	return res
}

func Union[T operator[T]](source1 []T, source2 []T) []T {
	// assuming all items in each source are unique (combination of name and type)
	// we need to find items in source1 those are not in source2 and vice versa
	source1MappedByKey := map[string]T{}
	source2MappedByKey := map[string]T{}
	for _, item := range source1 {
		source1MappedByKey[item.Key()] = item
	}
	for _, item := range source2 {
		source2MappedByKey[item.Key()] = item
	}

	res := []T{}

	// find items in source2 that need to be deleted
	for _, item := range source2 {
		_, ok := source1MappedByKey[item.Key()]
		if !ok {
			// reverse sign to negative
			new := item
			new = new.SetSign(SignMinus)
			res = append(res, new)
		}
	}

	// find items in source1 that need to be added
	for _, item := range source1 {
		_, ok := source2MappedByKey[item.Key()]
		if !ok {
			// make sure of postive sign
			new := item
			new = new.SetSign(SignPlus)
			res = append(res, new)
		}
	}

	// find intersections between source1 and source2
	for _, s1 := range source1 {
		s2, ok := source2MappedByKey[s1.Key()]
		if ok {
			curr := s1.Intersect(s2)
			if curr == nil {
				continue
			}

			for _, c := range *curr {
				res = append(res, c)
			}
		}
	}

	return res
}

var _1 []SignedField = Union([]SignedField{}, []SignedField{})
var _2 []SignedIndex = Union([]SignedIndex{}, []SignedIndex{})
var _3 []SignedMetadata = Union([]SignedMetadata{}, []SignedMetadata{})
var _4 []SignedCollection = Union([]SignedCollection{}, []SignedCollection{})
