/*
Copyright (c) 2023-present the go-mongr8 Authors and Contributors
[@see Authors file]

Licensed under the MIT License
(https://opensource.org/licenses/MIT)
*/
package generator

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/amirkode/go-mongr8/collection"
	"github.com/amirkode/go-mongr8/internal/config"
)

const (
	mongr8Path = "mongr8"
	baseCollectionPath = mongr8Path + "/collection"
)

func LoadCollection() []collection.Collection {
	collections := []collection.Collection{}

	return collections
}

func getAllCollectionStructs() []string {
	res := []string{}
	rootPath, err := config.GetProjectRootDir()
	if err != nil {
		return res
	}

	path := fmt.Sprintf("%s/%s", *rootPath, baseCollectionPath)
	collectionFileNames := config.GetAllFileNames(path)
	for _, name := range collectionFileNames {
		filePath := fmt.Sprintf("%s/%s", path, name)
		structName, err := getCollectionStructName(filePath)
		if err != nil {
			fmt.Println(err)
			continue
		}

		res = append(res, *structName)
	}

	return res
}

func getCollectionStructName(filePath string) (*string, error) {
	source := config.GetFileContent(filePath)
	pattern := `type [a-zA-Z0-9]+ struct`
	re := regexp.MustCompile(pattern)
	matchedStr := re.FindString(source)
	if matchedStr == "" {
		return nil, fmt.Errorf("no valid collection struct was found")
	}

	// get the second word
	res := strings.Split(matchedStr, " ")[1]

	return &res, nil
}

func getAllCollectionNames() []string {
	res := []string{}
	rootPath, err := config.GetProjectRootDir()
	if err != nil {
		return res
	}

	path := fmt.Sprintf("%s/%s", *rootPath, baseCollectionPath)
	collectionFileNames := config.GetAllFileNames(path)
	for _, name := range collectionFileNames {
		filePath := fmt.Sprintf("%s/%s", path, name)
		structName, err := getCollectionName(filePath)
		if err != nil {
			fmt.Println(err)
			continue
		}

		res = append(res, *structName)
	}

	return res
}

func getCollectionName(filePath string) (*string, error) {
	source := config.GetFileContent(filePath)
	pattern := `metadata.InitMetadata\("[a-zA-Z0-9_]+"\)`
	re := regexp.MustCompile(pattern)
	matchedStr := re.FindString(source)
	if matchedStr == "" {
		return nil, fmt.Errorf("no valid collection name was found")
	}

	// get anything inside the double quotes
	res := strings.Split(matchedStr, "\"")[1]

	return &res, nil
}
