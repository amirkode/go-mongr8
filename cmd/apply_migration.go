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

// applyMigrationCmd represents the applyMigration command
var applyMigrationCmd = &cobra.Command{
	Use:   "apply-migration",
	Short: "Apply all migrations",
	Long: `Apply migration changes to MongoDB`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("applyMigration called")
	},
}

func init() {
	rootCmd.AddCommand(applyMigrationCmd)
}
