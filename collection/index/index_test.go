/*
Copyright (c) 2023-present the go-mongr8 Authors and Contributors
[@see Authors file]

Licensed under the MIT License
(https://opensource.org/licenses/MIT)
*/
package index

import (
	"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestSingleFieldIndex(t *testing.T) {
	Convey("Case 1: Normal", t, func() {
		// Case 1: normal
		case1Index := SingleFieldIndex(Field("name", 1))

		Convey("Unexpected index type", func() {
			So(case1Index.spec.Type, ShouldEqual, TypeSingleField)
		})
	})

	Convey("Case 2: index field with no value", t, func() {
		defer func() {
			if r := recover(); r != nil {
				Convey("Unexpected panic", func() {
					So(fmt.Sprintf("%v", r), ShouldContainSubstring, "Value should be provided")
				})
			}
		}()

		SingleFieldIndex(Field("name"))
	})
}

func TestCompoundIndex(t *testing.T) {
	Convey("Case 1: normal", t, func() {
		case1Index := CompoundIndex(
			Field("name", 1),
			Field("age", -1),
		)

		Convey("Unexpected index type", func() {
			So(case1Index.spec.Type, ShouldEqual, TypeCompound)
		})
	})

	Convey("Case 2: index field with no value", t, func() {
		defer func() {
			if r := recover(); r != nil {
				Convey("Unexpected panic", func() {
					So(fmt.Sprintf("%v", r), ShouldContainSubstring, "Value should be provided")
				})
			}
		}()

		CompoundIndex(
			Field("name"),
			Field("age"),
		)
	})

	Convey("Case 3: index with no field", t, func() {
		defer func() {
			if r := recover(); r != nil {
				Convey("Unexpected panic", func() {
					So(fmt.Sprintf("%v", r), ShouldContainSubstring, "Index must have at least a field")
				})
			}
		}()

		CompoundIndex()
	})
}
