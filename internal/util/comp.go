/*
Copyright (c) 2023-present the go-mongr8 Authors and Contributors
[@see Authors file]

Licensed under the MIT License
(https://opensource.org/licenses/MIT)
*/
package util

type Comparable[T any] interface {
	// the values are -1 (less than), 0 (same), or 1 (greater than)
	CompareTo(T) int
}

func InList[T Comparable[T]](what T, in []T) bool {
	for _, i := range in {
		if i.CompareTo(what) == 0 {
			return true
		}
	}

	return false
}

func InListEq[T comparable](what T, in []T) bool {
	for _, i := range in {
		if what == i {
			return true
		}
	}

	return false
}

func NotInList[T comparable](what T, in []T) bool {
	for _, i := range in {
		if what == i {
			return false
		}
	}

	return true
}
