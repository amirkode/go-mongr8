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

// consolidateMigrationCmd represents the consolidate-migration command
var consolidateMigrationCmd = &cobra.Command{
	Use:   "consolidate-migration",
	Short: "Consolidate migration with current database schema",
	Long: `This command will consolidate current migration files with current schema/data in database`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("consolidateMigration called")
	},
}

func init() {
	rootCmd.AddCommand(consolidateMigrationCmd)
}
