/*
Copyright (c) 2023 the go-mongr8 Authors and Contributors
[@see Authors file]

Licensed under the MIT License
(https://opensource.org/licenses/MIT)
*/
package cmd

import (
	"fmt"
	"log"
	"os/exec"

	"github.com/amirkode/go-mongr8/internal/config"

	"github.com/amirkode/go-mongr8/migration/option"
	"github.com/spf13/cobra"
)

// generateMigrationCmd represents the generate-migration command
var generateMigrationCmd = &cobra.Command{
	Use:   "generate-migration",
	Short: "Generate migration files",
	Long: `Generate migration files based on defined collections`,
	Run: func(cmd *cobra.Command, args []string) {
		projectPath, err := config.GetProjectRootDir()
		if err != nil {
			log.Printf("Error generating migration: %s", err.Error())
			return
		}

		// name := cmd.PersistentFlags().Lookup("name").Value.String()
		flags := []string{"run", "main.go"}
		for _, flag := range []string{
			option.MigrationOptionArgUseSortedSchema,
			option.MigrationOptionArgUseForceConversion,
			option.MigrationOptionArgUseSchemaValidation,
			option.MigrationOptionArgDesc,
		} {
			currFlag := cmd.PersistentFlags().Lookup(flag)
			if currFlag != nil {
				value := currFlag.Value.String()
				flags = append(flags, fmt.Sprintf("-%s", flag))
				// if not boolean
				if value != "true" {
					flags = append(flags, value)
				}
			}
		}

		migrationCmdPath := fmt.Sprintf("%s/mongr8/cmd/generate", *projectPath)
		migrationCmd := exec.Command("go", flags...)
		migrationCmd.Dir = migrationCmdPath
		output, err := migrationCmd.CombinedOutput()
		if err != nil {
			log.Printf("Error generating migration: %s: %s\n", err.Error(), output)
			return
		}

		// print original output
		fmt.Printf("%s\n", output)
	},
}

func init() {
	rootCmd.AddCommand(generateMigrationCmd)

	generateMigrationCmd.PersistentFlags().Bool(option.MigrationOptionArgUseSortedSchema, true, "Use sorted schema on migration")
	generateMigrationCmd.PersistentFlags().Bool(option.MigrationOptionArgUseForceConversion, true, "Force on type convertion on migration")
	generateMigrationCmd.PersistentFlags().Bool(option.MigrationOptionArgUseSchemaValidation, true, "Apply schema validation on migration")
	generateMigrationCmd.PersistentFlags().String(option.MigrationOptionArgDesc, "", "Description for current migration")	
}
