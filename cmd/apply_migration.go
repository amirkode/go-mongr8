/*
Copyright (c) 2023-present the go-mongr8 Authors and Contributors
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

// applyMigrationCmd represents the applyMigration command
var applyMigrationCmd = &cobra.Command{
	Use:   "apply-migration",
	Short: "Apply all migrations",
	Long:  `Apply migration changes to MongoDB`,
	Run: func(cmd *cobra.Command, args []string) {
		projectPath, err := config.GetProjectRootDir()
		if err != nil {
			log.Printf("Error applying migration: %s", err.Error())
			return
		}

		// name := cmd.PersistentFlags().Lookup("name").Value.String()
		flags := []string{"run", "main.go"}
		for _, flag := range []string{
			option.MigrationOptionArgUseSortedSchema,
			option.MigrationOptionArgUseForceConversion,
			option.MigrationOptionArgUseSchemaValidation,
			option.MigrationOptionArgUseTransaction,
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

		migrationCmdPath := fmt.Sprintf("%s/mongr8/cmd/apply", *projectPath)
		migrationCmd := exec.Command("go", flags...)
		migrationCmd.Dir = migrationCmdPath
		output, err := migrationCmd.CombinedOutput()
		if err != nil {
			log.Printf("Error applying migration: %s: %s\n", err.Error(), output)
			return
		}

		// print original output
		fmt.Printf("%s", output)
	},
}

func init() {
	rootCmd.AddCommand(applyMigrationCmd)
}
