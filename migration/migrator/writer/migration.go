/*
Copyright (c) 2023-present the go-mongr8 Authors and Contributors
[@see Authors file]

Licensed under the MIT License
(https://opensource.org/licenses/MIT)
*/
package writer

import (
	"fmt"
	"time"

	"github.com/amirkode/go-mongr8/internal/config"
	"github.com/amirkode/go-mongr8/internal/util"
	"github.com/amirkode/go-mongr8/migration/migrator"
)

func getMigrationLiteral(migration migrator.Migration) string {
	literalUpActions := ""
	for _, action := range migration.Up {
		literalUpActions += fmt.Sprintf("%s,\n", action.GetLiteralInstance("si.", true))
	}

	literalDownActions := ""
	for _, action := range migration.Down {
		literalDownActions += fmt.Sprintf("%s,\n", action.GetLiteralInstance("si.", true))
	}

	res := fmt.Sprintf(`migrator.Migration{
		ID: "%s",
		Desc: "%s",
		Up: []si.Action{
			%s
		},
		Down: []si.Action{
			%s
		},
	}`, migration.ID, migration.Desc, literalUpActions, literalDownActions)

	// fmt.Println(res)

	return res
}

func Write(migration migrator.Migration) error {
	suffix, err := getNextSuffix()
	if err != nil {
		return err
	}

	/// projectPath should be the root project directory
	projectPath, err := config.GetProjectRootDir()
	if err != nil {
		return err
	}
	// packagePath should be the root go-mongr8 package directory
	packagePath, err := config.GetPackageDir()
	if err != nil {
		return err
	}

	// check whether current migration uses field or/and index
	useField := false
	useIndex := false
	for _, action := range migration.Up {
		for _, subAction := range action.SubActions {
			if !useField {
				useField = len(subAction.ActionSchema.Fields) > 0
			}
			if !useIndex {
				useIndex = len(subAction.ActionSchema.Indexes) > 0
			}
			// all states are found
			if useField && useIndex {
				break
			}
		}
		// all states are found
		if useField && useIndex {
			break
		}
	}

	tplVar := struct {
		CreateDate      string
		MigrationSuffix int
		Migration       string
		UseField        bool
		UseIndex        bool
	}{
		CreateDate:      time.Now().Format("2006-01-02"),
		MigrationSuffix: suffix,
		Migration:       getMigrationLiteral(migration),
		UseField:        useField,
		UseIndex:        useIndex,
	}

	// init templates
	tplPath := fmt.Sprintf("%s/migration/migrator/writer/template.tpl", *packagePath)
	outputPath := fmt.Sprintf("%s/mongr8/migration/%s.go", *projectPath, migration.ID)

	err = util.GenerateTemplate("migration", tplPath, outputPath, tplVar, true)
	if err != nil {
		return err
	}

	// updated migration variable names
	migrationVarNames, err := getMigrationVarNames()
	if err != nil {
		return err
	}

	baseTplVar := struct {
		CreateDate string
		Migrations []string
	}{
		CreateDate: time.Now().Format("2006-01-02"),
		Migrations: migrationVarNames,
	}

	// init templates
	tplPath = fmt.Sprintf("%s/migration/migrator/writer/template.tpl", *packagePath)
	outputPath = fmt.Sprintf("%s/mongr8/migration/base.go", *projectPath)

	return util.GenerateTemplate("migrations", tplPath, outputPath, baseTplVar, true)
}
