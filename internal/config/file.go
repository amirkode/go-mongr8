package config

import (
	"fmt"
	"os"
	"strings"

	"io/ioutil"
)

func GetFileContent(path string) string {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println("Error:", err)
		return ""
	}

	return string(content)
}

// this gets the project module name with provided root path
func GetProjectRootModuleName(rootPath string) string {
	content := GetFileContent(fmt.Sprintf("%s/go.mod", rootPath))
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "module ") {
			moduleName := strings.TrimSpace(strings.TrimPrefix(line, "module"))
			return moduleName
		}
	}

	return ""
}

// get all names inside a directory
func GetAllFileNames(path string) []string {
	res := []string{}
	dir, err := os.Open(path)
	if err != nil {
		return res
	}
	defer dir.Close()

	// read all entries inside the current dir
	allEntries, err := dir.ReadDir(0)
	if err != nil {
		return res
	}

	// list all files
	for _, entry := range allEntries {
		if entry.IsDir() {
			continue
		}

		if strings.HasSuffix(entry.Name(), ".go") {
			res = append(res, entry.Name())
		}
	}

	return res
}