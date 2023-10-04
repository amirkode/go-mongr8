/*
Copyright (c) 2023 the go-mongr8 Author and Contributors
[@see Authors file]

Licensed under the MIT License
(https://opensource.org/licenses/MIT)
*/
package cmd

import (
	"fmt"

	migration_init "github.com/amirkode/go-mongr8/migration/init"
	"github.com/spf13/cobra"
)

// initMigrationCmd represents the migration command
var initMigrationCmd = &cobra.Command{
	Use:   "init-migration",
	Short: "Initialize migration components",
	Long: `Initialize all migration components in the main working project directory`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("init-migration called")
		fmt.Println(args)
		applyRootDirValidation := cmd.Flags().Lookup("apply-root-dir-validation").Value.String() == "true"
		fmt.Println(applyRootDirValidation)
		err := migration_init.InitMigration(applyRootDirValidation)
		if err != nil {
			fmt.Printf("Error initiating migration: %v", err)
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
