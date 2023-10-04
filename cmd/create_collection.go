/*
Copyright (c) 2023 the go-mongr8 Author and Contributors
[@see Authors file]

Licensed under the MIT License
(https://opensource.org/licenses/MIT)
*/
package cmd

import (
	"fmt"

	"github.com/amirkode/go-mongr8/collection/generator"
	"github.com/spf13/cobra"
)

// createCollectionCmd represents the create-collection command
var createCollectionCmd = &cobra.Command{
	Use:   "create-collection",
	Short: "Create a new collection entity",
	Long:  `Create collection entity with `,
	Run: func(cmd *cobra.Command, args []string) {
		name := cmd.PersistentFlags().Lookup("name").Value.String()
		err := generator.GenerateMigrationTemplate(name)
		if err != nil {
			fmt.Println("Error creating collection: " + err.Error())
		}
	},
}

func init() {
	rootCmd.AddCommand(createCollectionCmd)
	
	createCollectionCmd.PersistentFlags().String("name", "", "MongoDB Collection Entity")
}
