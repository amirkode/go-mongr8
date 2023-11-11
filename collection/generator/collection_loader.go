/*
Copyright (c) 2023 the go-mongr8 Authors and Contributors
[@see Authors file]

Licensed under the MIT License
(https://opensource.org/licenses/MIT)
*/
package generator

import (
	"fmt"
	"internal/config"
	"regexp"
	"strings"

	"github.com/amirkode/go-mongr8/collection"
)

const (
	baseCollectionPath = "mongr8/collection"
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
	codeStr := config.GetFileContent(filePath)
	pattern := `type [a-zA-Z0-9]+ struct`
	re := regexp.MustCompile(pattern)
	matchedStr := re.FindString(codeStr)
	if matchedStr == "" {
		return nil, fmt.Errorf("No valid collection struct was found")
	}

	// get the second word
	res := strings.Split(matchedStr, " ")[1]

	return &res, nil
}