/*
Copyright (c) 2023 the go-mongr8 Authors and Contributors
[@see Authors file]

Licensed under the MIT License
(https://opensource.org/licenses/MIT)
*/
package cmd

import (
	"log"
	"os/exec"

	migration_init "github.com/amirkode/go-mongr8/migration/init"
	"github.com/spf13/cobra"
)

// initMigrationCmd represents the migration command
var initMigrationCmd = &cobra.Command{
	Use:   "init-migration",
	Short: "Initialize migration components",
	Long: `Initialize all migration components in the main working project directory`,
	Run: func(cmd *cobra.Command, args []string) {
		err := migration_init.InitMigration()
		// add go-mongr8 to current project
		output, err := exec.Command("go", "get", "-u", "github.com/amirkode/go-mongr8@latest").CombinedOutput()
		if err != nil {
			log.Printf("Error initiating migration: %s: %s\n", err.Error(), output)
			return
		}

		err = migration_init.InitMigration()
		if err != nil {
			log.Printf("Error initiating migration: %v", err)
		} else {
			log.Println("Migration folder 'mongr8' has been initiated")
		}

		// tidy everything up
		output, err = exec.Command("go", "mod", "tidy").CombinedOutput()
		if err != nil {
			log.Printf("Error initiating migration: %s: %s\n", err.Error(), output)
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(initMigrationCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	initMigrationCmd.PersistentFlags().String("type", "", "MongoDB Migration type")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	
	initMigrationCmd.Flags().BoolP("apply-root-dir-validation", "v", false, "Help message for toggle")
}
