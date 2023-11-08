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

const (
	tplMongr8               = "mongr8_info"
	tplConfig               = "config"
	tplCombainedCollections = "combined_collections"
	tplCmdMain              = "cmd_main"
	tplCmdCall              = "cmd_call"
	tplMigrations           = "migrations"
)

// init mongr8 migration structure
// this will generate all required folders for mongr8
// project-root/
// ├── mongr8/
// |   ├── cmd/
// |       ├── contains migration commands
// |   ├── collection/
// |       ├── contains collection schema definitions
// |   ├── config/
// |       ├── contains some setup files
// |   ├── migration/
// |       ├── contains some migration files
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
	if err = initFolderStructure(*projectPath, applyRootDirValidation); err != nil {
		return err
	}

	// init templates
	tplPath := fmt.Sprintf("%s/migration/init/template.tpl", *packagePath)

	// init mongr8.info file
	if err = initMongr8Info(*projectPath, tplPath); err != nil {
		return err
	}

	// init config file
	if err = initConfig(*projectPath, tplPath); err != nil {
		return err
	}

	// init combined collections
	if err = initCombinedCollections(*packagePath, tplPath); err != nil {
		return err
	}

	// init migration sub action schemas
	if err = initMigrationSubActionSchemas(*packagePath, tplPath); err != nil {
		return err
	}

	// init cmds
	err = initCmd(*packagePath, tplPath)

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
		fmt.Sprintf("%s/cmd/apply", mainDir),
		fmt.Sprintf("%s/cmd/consolidate", mainDir),
		fmt.Sprintf("%s/cmd/generate", mainDir),
		fmt.Sprintf("%s/collection/no_edit", mainDir),
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

func initMongr8Info(projectPath, tplPath string) error {
	tplVar := struct {
		CreateDate string
		Version    string
	}{
		CreateDate: time.Now().Format("2006-01-02"),
		Version:    migration.Mongr8Version,
	}

	outputPath := fmt.Sprintf("%s/mongr8/mongr8.info", projectPath)

	return util.GenerateTemplate(tplMongr8, tplPath, outputPath, tplVar, false)
}

func initConfig(projectPath, tplPath string) error {
	tplVar := struct {
		CreateDate string
	}{
		CreateDate: time.Now().Format("2006-01-02"),
	}

	outputPath := fmt.Sprintf("%s/mongr8/config/config.go", projectPath)

	return util.GenerateTemplate(tplConfig, tplPath, outputPath, tplVar, true)
}

func initCombinedCollections(projectPath, tplPath string) error {
	tplVar := struct {
		CreateDate string
	}{
		CreateDate: time.Now().Format("2006-01-02"),
	}

	outputPath := fmt.Sprintf("%s/mongr8/collection/no_edit/combined_collections.go", projectPath)

	return util.GenerateTemplate(tplCombainedCollections, tplPath, outputPath, tplVar, true)
}

func initMigrationSubActionSchemas(projectPath, tplPath string) error {
	tplVar := struct {
		CreateDate string
	}{
		CreateDate: time.Now().Format("2006-01-02"),
	}

	outputPath := fmt.Sprintf("%s/mongr8/migration/base.go", projectPath)

	return util.GenerateTemplate(tplMigrations, tplPath, outputPath, tplVar, true)
}

func initCmd(projectPath, tplPath string) error {
	createDate := time.Now().Format("2006-01-02")
	moduleName := config.GetProjectRootModuleName(projectPath)

	// generate /mongr8/cmd/cmd.go
	tplCmdMainVar := struct {
		CreateDate string
		ModuleName string
	}{
		CreateDate: createDate,
		ModuleName: moduleName,
	}
	outputPath := fmt.Sprintf("%s/mongr8/cmd/cmd.go", projectPath)
	err := util.GenerateTemplate(tplCmdMain, tplPath, outputPath, tplCmdMainVar, true)
	if err != nil {
		return err
	}

	// generate /mongr8/cmd/[operation]/main.go
	outputs := []struct {
		operation string
		funcName  string
	}{
		{
			operation: "apply",
			funcName:  "CmdApplyMigration",
		},
		{
			operation: "consolidate",
			funcName:  "CmdConsolidateMigration",
		},
		{
			operation: "generate",
			funcName:  "CmdGenerateMigration",
		},
	}
	for _, output := range outputs {
		tplCmdCallVar := struct {
			CreateDate string
			ModuleName string
			FuncName   string
		}{
			CreateDate: createDate,
			ModuleName: moduleName,
			FuncName:   output.funcName,
		}
		outputPath = fmt.Sprintf("%s/mongr8/cmd/%s/main.go", projectPath, output.operation)
		err := util.GenerateTemplate(tplCmdCall, tplPath, outputPath, tplCmdCallVar, true)
		if err != nil {
			return err
		}
	}

	return nil
}
