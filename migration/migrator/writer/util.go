/*
Copyright (c) 2023-present the go-mongr8 Authors and Contributors
[@see Authors file]

Licensed under the MIT License
(https://opensource.org/licenses/MIT)
*/
package writer

import (
	"fmt"
	"strconv"

	"github.com/amirkode/go-mongr8/internal/config"
	"github.com/amirkode/go-mongr8/internal/validation"
)

const (
	baseMigrationPath = "mongr8/migration"
)

func getMigrationVarNames() ([]string, error) {
	res := []string{}
	rootPath, err := config.GetProjectRootDir()
	if err != nil {
		return res, err
	}

	path := fmt.Sprintf("%s/%s", *rootPath, baseMigrationPath)
	collectionFileNames := config.GetAllFileNames(path)
	for _, name := range collectionFileNames {
		if !validation.ValidateWithRegex(name, `^\d{8}_\d{6}.go$`) {
			continue
		}

		filePath := fmt.Sprintf("%s/%s", path, name)
		codeStr := config.GetFileContent(filePath)
		matchedStr := validation.FindWithRegex(codeStr, `Migration[0-9]+`)
		if matchedStr == "" {
			continue
		}

		res = append(res, matchedStr)
	}

	return res, nil
}

func getNextSuffix() (int, error) {
	res := 1
	migrationVarNames, err := getMigrationVarNames()
	if err != nil {
		return res, err
	}

	for _, varName := range migrationVarNames {
		currSuffix, _ := strconv.Atoi(varName[9:])
		if currSuffix >= res {
			res = currSuffix + 1
		}
	}

	return res, nil
}
