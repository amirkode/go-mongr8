/*
Copyright (c) 2023 the go-mongr8 Authors and Contributors
[@see Authors file]

Licensed under the MIT License
(https://opensource.org/licenses/MIT)
*/
package field

type Spec struct {
	Name string
	// Type of the field
	Type FieldType
	// Array item, if current type is an array
	ArrayField *Spec
	// Children of object, if current type is an object
	Object *[]Spec
	// Nullable flag
	Nullable bool
}

type BaseSpec struct {
	spec *Spec
}

func (b *BaseSpec) Spec() *Spec {
	return b.spec
}

func (b *BaseSpec) SetArrayField(s *BaseSpec) *BaseSpec {
	b.spec.ArrayField = s.spec
	return b
}

func (b *BaseSpec) SetObject(object *[]BaseSpec) *BaseSpec {
	objects := make([]Spec, len(*object))
	for i, o := range *object {
		objects[i] = *o.spec
	}

	b.spec.Object = &objects

	return b
}

func (b *BaseSpec) IsNullable() *BaseSpec {
	b.spec.Nullable = true

	return b
}

func baseField(name string, fieldType FieldType) *BaseSpec {
	field := &BaseSpec{
		&Spec{
			Name: name,
			Type: fieldType,
			Nullable: false,
		},
	}

	return field
}

func StringField(name string) *BaseSpec {
	return baseField(name, TypeString)
}

func Int64Field(name string) *BaseSpec {
	return baseField(name, TypeInt64)
}

func Int32Field(name string) *BaseSpec {
	return baseField(name, TypeInt32)
}

func DoubleField(name string) *BaseSpec {
	return baseField(name, TypeDouble)
}

func ArrayField(name string) *BaseSpec {
	return baseField(name, TypeArray)
}

func ObjectField(name string) *BaseSpec {
	return baseField(name, TypeObject)
}
