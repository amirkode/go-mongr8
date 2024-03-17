package generator

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/amirkode/go-mongr8/internal/config"
	"github.com/amirkode/go-mongr8/internal/util"

	. "github.com/smartystreets/goconvey/convey"
)

func testSetupCollectionFolder() {
	rootPath, err := config.GetProjectRootDir()
	if err != nil {
		panic(err)
	}

	path := fmt.Sprintf("%s/%s/no_edit", *rootPath, baseCollectionPath)
	if config.DoesPathExist(path) {
		return
	}

	if err := os.MkdirAll(path, os.ModePerm); err != nil {
		panic(err)
	}

	// init combined collection
	tplVar := struct {
		CreateDate string
	}{
		CreateDate: time.Now().Format("2006-01-02"),
	}

	tplPath := fmt.Sprintf("%s/migration/init/template.tpl", *rootPath)
	outputPath := fmt.Sprintf("%s/mongr8/collection/no_edit/combined_collections.go", *rootPath)
	err = util.GenerateTemplate(tplCombainedCollections, tplPath, outputPath, tplVar, true)
	if err != nil {
		panic(err)
	}
}

func testSetupCollection(collectionName string) {
	testSetupCollectionFolder()

	rootPath, err := config.GetProjectRootDir()
	if err != nil {
		panic(err)
	}

	tplPath := fmt.Sprintf("%s/collection/generator/template.tpl", *rootPath)
	collectionName = util.ToSnakeCase(collectionName)
	templateVar := &CollectionTemplateVar{
		CreateDate: time.Now().Format("2006-01-02"),
		Entity:     util.ToCapitalizedCamelCase(collectionName),
		Collection: collectionName,
	}

	// genenrate collection
	outputPath := fmt.Sprintf("%s/mongr8/collection/%s.go", *rootPath, collectionName)
	err = util.GenerateTemplate(tplCollection, tplPath, outputPath, templateVar, true)
	if err != nil {
		panic(err)
	}

	// generate combined collections
	combinedCollsTemplateVar, err := getCombinedCollectionsTemplateVar(*rootPath)
	if err != nil {
		panic(err)
	}

	outputPath = fmt.Sprintf("%s/mongr8/collection/no_edit/combined_collections.go", *rootPath)
	err = util.GenerateTemplate(tplCombainedCollections, tplPath, outputPath, combinedCollsTemplateVar, true)
	if err != nil {
		panic(err)
	}
}

func testTearDown() {
	rootPath, err := config.GetProjectRootDir()
	if err != nil {
		panic(err)
	}

	path := fmt.Sprintf("%s/%s", *rootPath, baseCollectionPath)
	if err := os.RemoveAll(path); err != nil {
		panic(err)
	}
}

func TestGetAllCollectionStructs(t *testing.T) {
	Convey("Case 1: Normal", t, func() {
		testTearDown()
		testSetupCollection("first_collection")
		testSetupCollection("second_collection")
		
		collectionStructs := getAllCollectionStructs()
		Convey("Unexpected collection structs length", func() {
			So(len(collectionStructs), ShouldEqual, 2)
		})

		expectedStructs := map[string]bool{
			"FirstCollection": true,
			"SecondCollection": true,
		}

		Convey("Unexpected struct name", func() {
			for _, name := range collectionStructs {
				_, ok := expectedStructs[name]
				So(ok, ShouldBeTrue)
			}
		})
	})
}

func TestGetCollectionStructName(t *testing.T) {
	Convey("Case 1: Normal", t, func() {
		testTearDown()
		rootPath, err := config.GetProjectRootDir()
		if err != nil {
			panic(err)
		}

		testSetupCollection("first_collection")
		testSetupCollection("second_collection")

		Convey("Unexpected struct name", func() {
			coll1Path := fmt.Sprintf("%s/%s/first_collection.go", *rootPath, baseCollectionPath)
			name1, err := getCollectionStructName(coll1Path)
			if err != nil {
				panic(err)
			}

			So(*name1, ShouldEqual, "FirstCollection")

			coll2Path := fmt.Sprintf("%s/%s/second_collection.go", *rootPath, baseCollectionPath)
			name2, err := getCollectionStructName(coll2Path)
			if err != nil {
				panic(err)
			}

			So(*name2, ShouldEqual, "SecondCollection")
		})
	})
}

func TestGetAllCollectionNames(t *testing.T) {
	Convey("Case 1: Normal snake case input", t, func() {
		testTearDown()
		testSetupCollection("first_collection")
		testSetupCollection("second_collection")
		
		collectionNames := getAllCollectionNames()
		Convey("Unexpected collection names length", func() {
			So(len(collectionNames), ShouldEqual, 2)
		})

		expectedNames := map[string]bool{
			"first_collection": true,
			"second_collection": true,
		}

		Convey("Unexpected collection name", func() {
			for _, name := range collectionNames {
				_, ok := expectedNames[name]
				So(ok, ShouldBeTrue)
			}
		})
	})

	Convey("Case 2: Camel case input", t, func() {
		testTearDown()
		testSetupCollection("FirstCollection")
		testSetupCollection("SecondCollection")
		
		collectionNames := getAllCollectionNames()
		Convey("Unexpected collection names length", func() {
			So(len(collectionNames), ShouldEqual, 2)
		})

		expectedNames := map[string]bool{
			"firstcollection": true,
			"secondcollection": true,
		}

		Convey("Unexpected collection name", func() {
			for _, name := range collectionNames {
				_, ok := expectedNames[name]
				So(ok, ShouldBeTrue)
			}
		})
	})
}

func TestGetCollectionName(t *testing.T) {
	Convey("Case 1: Normal snake case input", t, func() {
		testTearDown()
		rootPath, err := config.GetProjectRootDir()
		if err != nil {
			panic(err)
		}

		testSetupCollection("first_collection")
		testSetupCollection("second_collection")

		Convey("Unexpected collection name", func() {
			coll1Path := fmt.Sprintf("%s/%s/first_collection.go", *rootPath, baseCollectionPath)
			name1, err := getCollectionName(coll1Path)
			if err != nil {
				panic(err)
			}

			So(*name1, ShouldEqual, "first_collection")

			coll2Path := fmt.Sprintf("%s/%s/second_collection.go", *rootPath, baseCollectionPath)
			name2, err := getCollectionName(coll2Path)
			if err != nil {
				panic(err)
			}

			So(*name2, ShouldEqual, "second_collection")
		})
	})

	Convey("Case 2: Camel case input", t, func() {
		testTearDown()
		rootPath, err := config.GetProjectRootDir()
		if err != nil {
			panic(err)
		}

		testSetupCollection("FirstCollection")
		testSetupCollection("SecondCollection")

		Convey("Unexpected collection name", func() {
			coll1Path := fmt.Sprintf("%s/%s/firstcollection.go", *rootPath, baseCollectionPath)
			name1, err := getCollectionName(coll1Path)
			if err != nil {
				panic(err)
			}

			So(*name1, ShouldEqual, "firstcollection")

			coll2Path := fmt.Sprintf("%s/%s/secondcollection.go", *rootPath, baseCollectionPath)
			name2, err := getCollectionName(coll2Path)
			if err != nil {
				panic(err)
			}

			So(*name2, ShouldEqual, "secondcollection")
		})
	})
}
