{{ define "mongr8_info" }}
Create date: {{ .CreateDate}}
Created by: go-mongr8
Version: {{ .Version}}

go-mongr8, a lightweight yet robust package for MongoDB migration management.
Simplify the management of MongoDB schema changes.

Copyright (c) 2023 the go-mongr8 Author and Contributors
[@see Authors file]

Licensed under the MIT License
(https://opensource.org/licenses/MIT)
{{ end }}

{{ define "config" }}
/*
Create date: {{ .CreateDate}}
Created by: go-mongr8

Copyright (c) 2023 the go-mongr8 Author and Contributors
[@see Authors file]

Licensed under the MIT License
(https://opensource.org/licenses/MIT)
*/
package config

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
)

type Config struct {
	context  *context.Context
	database *mongo.Database
}

// global configuration
var conf = Config{}

// this provides a common context for any migration operation
func GlobalContext() context.Context {
	if conf.context == nil {
		ctx := context.Background()
		conf.context = &ctx
	}

	return *conf.context
}

// this provide a global database connection across the migration processes
func Database() mongo.Database {
	if conf.database == nil {
		// init database
		initDatabase()
	}

	return *conf.database
}

func initDatabase() {
	/*
	   do something here to init Config.Database
	   you can call an existing mongodb database instance from the project

	   here's simple example of direct database initialization
	   add import "go.mongodb.org/mongo-driver/mongo/options"
	   code example:
	   ctx := GlobalContext()
	   clientOpts := options.Client().ApplyURI(fmt.Sprintf("mongodb://%s", "db host goes here")).
	           SetDirect(true)

	   credential := options.Credential{
	       Username:   "username goes here",
	       Password:   "password goes here",
	       AuthSource: "auth db goes here",
	   }
	   clientOpts.SetAuth(credential)

	   client, err := mongo.Connect(ctx, clientOpts)
	   if err != nil {
	       db := client.Database("main db goes here")
	       conf.database = db
	   }
	*/
}
{{ end }}

{{ define "combined_collections" }}
/*
DOT NOT EDIT, THIS FILE IS GENERATED BY CODE GEN
Create date: {{ .CreateDate}}
Created by: go-mongr8

Copyright (c) 2023 the go-mongr8 Authors and Contributors
[@see Authors file]

Licensed under the MIT License
(https://opensource.org/licenses/MIT)
*/

package no_edit

import (
	"github.com/amirkode/go-mongr8/collection"
)

func GetAllCollections() []collection.Collection {
	res := []collection.Collection{}

	return res
}
{{ end }}

{{ define "cmd_main" }}
/*
DOT NOT EDIT, THIS FILE IS GENERATED BY CODE GEN
Create date: {{ .CreateDate}}
Created by: go-mongr8

Copyright (c) 2023 the go-mongr8 Authors and Contributors
[@see Authors file]

Licensed under the MIT License
(https://opensource.org/licenses/MIT)
*/
package cmd

import (
	"{{ .ModuleName}}/mongr8/collection/no_edit"
	"github.com/amirkode/go-mongr8/migration"
)

func CmdGenerateMigration() {
	collections := no_edit.GetAllCollections()
	migration := migration.NewMigration()
	err := migration.GenerateMigration(collections)
	if err != nil {

	}
}

func CmdApplyMigration() {
	migration := migration.NewMigration()
	err := migration.ApplyMigration()
	if err != nil {
		
	}
}

func CmdConsolidateMigration() {
	collections := no_edit.GetAllCollections()
	migration := migration.NewMigration()
	err := migration.ConsolidateMigration(collections)
	if err != nil {

	}
}
{{ end }}

{{ define "cmd_call"}}
/*
DOT NOT EDIT, THIS FILE IS GENERATED BY CODE GEN
Create date: {{ .CreateDate}}
Created by: go-mongr8

Copyright (c) 2023 the go-mongr8 Authors and Contributors
[@see Authors file]

Licensed under the MIT License
(https://opensource.org/licenses/MIT)
*/
package main

import (
	"{{ .ModuleName}}/mongr8/cmd"
)

func main() {
	cmd.{{ .FuncName}}()
}
{{ end }}