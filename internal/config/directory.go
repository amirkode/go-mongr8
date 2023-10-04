/*
Copyright (c) 2023 the go-mongr8 Authors and Contributors
[@see Authors file]

Licensed under the MIT License
(https://opensource.org/licenses/MIT)
*/
package config

import (
	"fmt"
	"os"
	"strings"

	"os/exec"
	"path/filepath"
)

// this function retrieves the root directory of go-mongr8
// and is used to get some project template drectories
// this function could replaced with a better approach in the future
func GetPackageDir() (*string, error) {
	// hard-coded package name, should be moved or passed as parameter in the future
	pkgName := "github.com/amirkode/go-mongr8"

	// get path using golang command "go list {{.Dir}} packageName"
	cmd := exec.Command("go", "list", "{{.Dir}}", pkgName)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Error getting directory: %v\n", err)
		// get project root directory instead, it might be for internal testing
		return GetProjectRootDir()
	}
	
	path := strings.TrimSpace(string(output))

	return &path, err
}


// this function retrieves the root directory of working project
func GetProjectRootDir() (*string, error) {
	// start from the current working directory
	currDir, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	fmt.Println("curr dir: " + currDir)

	// continue moving up the directory tree until we find a marker file or reach the root
	for {
		// check if go.mod exists in the current directory
		if _, err := os.Stat(filepath.Join(currDir, "go.mod")); err == nil {
			return &currDir, nil
		}

		// move up one directory
		parentDir := filepath.Dir(currDir)

		// reached the device root directory
		if parentDir == currDir {
			return nil, fmt.Errorf("Project root directory not found")
		}

		currDir = parentDir
	}
}

func DoesPathExist(path string) bool {
	_, err := os.Stat(path)

	return err == nil
}