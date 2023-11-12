/*
Copyright (c) 2023 the go-mongr8 Authors and Contributors
[@see Authors file]

Licensed under the MIT License
(https://opensource.org/licenses/MIT)
*/
package cmd

import (
	"log"

	"github.com/amirkode/go-mongr8/collection/generator"
	"github.com/spf13/cobra"
)

// createCollectionCmd represents the create-collection command
var createCollectionCmd = &cobra.Command{
	Use:   "create-collection",
	Short: "Create a new collection entity",
	Long:  `Create collection entity with `,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			log.Println("Error creating collection: Invalid collection name")
			return
		}

		// name := cmd.PersistentFlags().Lookup("name").Value.String()
		name := args[0]
		err := generator.GenerateMigrationTemplate(name)
		if err != nil {
			log.Println("Error creating collection: " + err.Error())
		} else {
			log.Printf("Collection '%s' has been created", name)
		}
	},
}

func init() {
	rootCmd.AddCommand(createCollectionCmd)
	
	// createCollectionCmd.PersistentFlags().String("name", "", "MongoDB Collection Entity")
}
