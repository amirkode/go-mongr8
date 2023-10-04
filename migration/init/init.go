/*
Copyright (c) 2023 the go-mongr8 Author and Contributors
[@see Authors file]

Licensed under the MIT License
(https://opensource.org/licenses/MIT)
*/
package init

import (
	"fmt"
	"os"
	"time"

	"internal/config"
	"internal/util"

	"github.com/amirkode/go-mongr8/migration"
)

// init mongr8 migration structure
// this will generate all required folders for mongr8
// project-root/
// ├── mongr8/
// |   ├── collection/
// |       ├── contains collection definitions
// |   ├── config/
// |       ├── contains some setup files
// |   ├── migration/
// |       ├── contains some setup files
func InitMigration(applyRootDirValidation bool) error {
	/// projectPath should be the root project directory
	projectPath, err := config.GetProjectRootDir()
	if err != nil {
		if !applyRootDirValidation {
			// if validation is not required
			// just take current working directory
			wd, err := os.Getwd()
			if err != nil {
				return err
			}

			// set current working directory to currDirr
			projectPath = &wd
		} else {
			return fmt.Errorf("You're not woking in any go project.")
		}
	}
	// packagePath should be the root go-mongr8 package directory
	packagePath, err := config.GetPackageDir()
	if err != nil {
		return err
	}

	// init folder structure
	if err = initFolderStructure(*projectPath, applyRootDirValidation); err != nil { return err }

	// init mongr8.info file
	if err = initMongr8Info(*projectPath, *packagePath); err != nil { return err }

	// init config file
	err = initConfig(*projectPath, *packagePath)

	// TODO: might add something in the future

	return err
}

func initFolderStructure(projectPath string, applyRootDirValidation bool) error {
	mainDir := fmt.Sprintf("%s/mongr8", projectPath)
	mongr8InfoDir := fmt.Sprintf("%s/mongr8.info", mainDir)
	if config.DoesPathExist(mongr8InfoDir) {
		return fmt.Errorf("The mongr8.lock file was already iniated. Please delete this file to continue.")
	}

	childrenDir := []string{
		fmt.Sprintf("%s/collection", mainDir),
		fmt.Sprintf("%s/migration", mainDir),
		fmt.Sprintf("%s/config", mainDir),
	}

	// init all directories
	for _, dir := range childrenDir {
		err := os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			return err
		}
	}

	return nil
}

func initMongr8Info(projectPath, packagePath string) error {
	tplVar := struct {
		CreateDate string
		Version    string
	}{
		CreateDate: time.Now().Format("2006-01-02"),
		Version:    migration.Mongr8Version,
	}

	tplPath := fmt.Sprintf("%s/migration/init/mongr8.tpl", packagePath)
	outputPath := fmt.Sprintf("%s/mongr8/mongr8.info", projectPath)

	return util.GenerateTemplate(tplPath, outputPath, tplVar)
}

func initConfig(projectPath, packagePath string) error {
	tplVar := struct {
		CreateDate string
	}{
		CreateDate: time.Now().Format("2006-01-02"),
	}

	tplPath := fmt.Sprintf("%s/migration/init/config/config.tpl", packagePath)
	outputPath := fmt.Sprintf("%s/mongr8/config/config.go", projectPath)

	return util.GenerateTemplate(tplPath, outputPath, tplVar)
}
