/*
Copyright (c) 2023 the go-mongr8 Authors and Contributors
[@see Authors file]

Licensed under the MIT License
(https://opensource.org/licenses/MIT)
*/
package data_type

type Pair[T any, U any] struct {
	First  T
	Second U
}

func NewPair[T any, U any](first T, second U) Pair[T, U] {
	return Pair[T, U]{
		First:  first,
		Second: second,
	}
}

type Set[T any] struct {
	items map[interface{}]bool
}

func (s Set[T]) Exists(item interface{}) bool {
	_, ok := s.items[item]
	return ok
}

func (s Set[T]) Insert(item interface{}) {
	s.items[item] = true
}

func (s Set[T]) Erase(item interface{}) {
	delete(s.items, item)
}
