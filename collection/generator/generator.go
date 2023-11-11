/*
Copyright (c) 2023 the go-mongr8 Authors and Contributors
[@see Authors file]

Licensed under the MIT License
(https://opensource.org/licenses/MIT)
*/
package generator

import (
	"errors"
	"fmt"
	"time"

	"github.com/amirkode/go-mongr8/internal/config"
	"github.com/amirkode/go-mongr8/internal/util"
)

const (
	tplCollection           = "collection"
	tplCombainedCollections = "combined_collections"
)

type CollectionTemplateVar struct {
	CreateDate string
	Entity     string
	Collection string
}

type CombinedCollectionsTemplateVar struct {
	CreateDate  string
	ModuleName  string
	Collections []string
}

func getCollectionTemplateVar(collectionName string) (*CollectionTemplateVar, error) {
	if len(collectionName) == 0 {
		return nil, errors.New("An empty string provided")
	}

	createDate := time.Now().Format("2006-01-02")
	entityName := util.ToCapitalizedCamelCase(collectionName)
	templateVar := &CollectionTemplateVar{
		CreateDate: createDate,
		Entity:     entityName,
		Collection: collectionName,
	}

	return templateVar, nil
}

func getCombinedCollectionsTemplateVar(rootPath string) (*CombinedCollectionsTemplateVar, error) {
	createDate := time.Now().Format("2006-01-02")
	moduleName := config.GetProjectRootModuleName(rootPath)

	existingCollections := getAllCollectionStructs()
	collections := []string{}
	for _, coll := range existingCollections {
		collections = append(collections, fmt.Sprintf("Instance%s", coll))
	}

	templateVar := &CombinedCollectionsTemplateVar{
		CreateDate:  createDate,
		ModuleName:  moduleName,
		Collections: collections,
	}

	return templateVar, nil
}

func GenerateMigrationTemplate(collectionName string) error {
	collTemplateVar, err := getCollectionTemplateVar(collectionName)
	if err != nil {
		return err
	}

	rootPath, err := config.GetProjectRootDir()
	packagePath, err := config.GetPackageDir()
	if err != nil {
		return err
	}

	tplPath := fmt.Sprintf("%s/collection/generator/template.tpl", *packagePath)

	// genenrate collection
	outputPath := fmt.Sprintf("%s/mongr8/collection/%s.go", *rootPath, collectionName)
	err = util.GenerateTemplate(tplCollection, tplPath, outputPath, collTemplateVar, true)

	// generate combined collections
	combinedCollsTemplateVar, err := getCombinedCollectionsTemplateVar(*rootPath)
	if err != nil {
		return err
	}

	outputPath = fmt.Sprintf("%s/mongr8/collection/no_edit/combined_collections.go", *rootPath)

	return util.GenerateTemplate(tplCombainedCollections, tplPath, outputPath, combinedCollsTemplateVar, true)
}
