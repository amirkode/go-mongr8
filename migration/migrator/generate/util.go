package generate

import (
	"fmt"
	"strconv"

	"internal/config"
	"internal/validation"
)

const (
	baseMigrationPath = "mongr8/migration"
)

func getMigrationVarNames(filePath string) (*string, error) {
	codeStr := config.GetFileContent(filePath)
	matchedStr := validation.FindWithRegex(codeStr, `Migration[0-9]+`)
	if matchedStr == "" {
		return nil, fmt.Errorf("No valid migration was found")
	}

	return &matchedStr, nil
}

func getNextSuffix() (int, error) {
	res := 1
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
		varName, err := getMigrationVarNames(filePath)
		if err != nil {
			fmt.Println(err)
			continue
		}

		// migration
		currSuffix, _ := strconv.Atoi((*varName)[9:len(*varName)])
		if currSuffix >= res {
			res = currSuffix + 1
		}
	}

	return res, nil
}