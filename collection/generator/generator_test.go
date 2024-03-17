package generator

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestGetCollectionTemplateVar(t *testing.T) {
	Convey("Case 1: Normal snake case input", t, func() {
		templ, err := getCollectionTemplateVar("a_collection")
		if err != nil {
			panic(err)
		}

		Convey("Unexpected collection name", func() {
			So(templ.Collection, ShouldEqual, "a_collection")
		})
		Convey("Unexpected entity name", func() {
			So(templ.Entity, ShouldEqual, "ACollection")
		})
	})

	Convey("Case 2: Camel case input", t, func() {
		templ, err := getCollectionTemplateVar("CollectionOne")
		if err != nil {
			panic(err)
		}

		Convey("Unexpected collection name", func() {
			So(templ.Collection, ShouldEqual, "collectionone")
		})
		Convey("Unexpected entity name", func() {
			So(templ.Entity, ShouldEqual, "Collectionone")
		})
	})

	Convey("Case 3: Cluttered input", t, func() {
		templ, err := getCollectionTemplateVar("_Collection-One*&#")
		if err != nil {
			panic(err)
		}

		Convey("Unexpected collection name", func() {
			So(templ.Collection, ShouldEqual, "collection_one")
		})
		Convey("Unexpected entity name", func() {
			So(templ.Entity, ShouldEqual, "CollectionOne")
		})
	})
}

func TestGetCombinedCollectionsTemplateVar(t *testing.T) {
	// TODO: implement test
}

func TestGenerateMigrationTemplate(t *testing.T) {
	// TODO: implement test
}
