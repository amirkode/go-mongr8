/*
Copyright (c) 2023 the go-mongr8 Author and Contributors
[@see Authors file]

Licensed under the MIT License
(https://opensource.org/licenses/MIT)
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// generateMigrationCmd represents the generate-migration command
var generateMigrationCmd = &cobra.Command{
	Use:   "generate-migration",
	Short: "Generate migration files",
	Long: `Generate migration files based on defined collections`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("generateMigration called")
	},
}

func init() {
	rootCmd.AddCommand(generateMigrationCmd)
}
