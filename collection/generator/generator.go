/*
Copyright (c) 2023 the go-mongr8 Author and Contributors
[@see Authors file]

Licensed under the MIT License
(https://opensource.org/licenses/MIT)
*/
package generator

import (
	"errors"
	"fmt"
	"time"

	"internal/config"
	"internal/util"
)

type TemplateVar struct {
	CreateDate string
	Entity string
	Collection string
}

func getTemplateVar(collectionName string) (*TemplateVar, error) {
	if len(collectionName) == 0 {
		return nil, errors.New("An empty string provided")
	}
	
	createDate := time.Now().Format("2006-01-02")
	entityName := util.ToCapitalizedCamelCase(collectionName)
	templateVar := &TemplateVar{
		CreateDate: createDate,
		Entity: entityName,
		Collection: collectionName,
	}

	return templateVar, nil
}

func GenerateMigrationTemplate(collectionName string) error {
	templateVar, err := getTemplateVar(collectionName)
	if err != nil {
		return err
	}

	rootPath, err := config.GetProjectRootDir()
	packagePath, err := config.GetPackageDir()
	if err != nil {
		return err
	}

	tplPath := fmt.Sprintf("%s/collection/generator/template.tpl", *packagePath)
	outputPath := fmt.Sprintf("%s/mongr8/collection/%s.go", *rootPath, collectionName)

	return util.GenerateTemplate(tplPath, outputPath, templateVar)
}