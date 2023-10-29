package writer

import (
	"fmt"
	"time"

	"internal/config"
	"internal/util"

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

func Write(migration migrator.Migration, suffix int) error {
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

	tplVar := struct {
		CreateDate      string
		MigrationSuffix int
		Migration       string
	}{
		CreateDate:      time.Now().Format("2006-01-02"),
		MigrationSuffix: suffix,
		Migration:       getMigrationLiteral(migration),
	}

	// init templates
	tplPath := fmt.Sprintf("%s/migration/migrator/writer/template.tpl", *packagePath)
	outputPath := fmt.Sprintf("%s/mongr8/migration/%s.go", *projectPath, migration.ID)

	return util.GenerateTemplate("migration", tplPath, outputPath, tplVar)
}
