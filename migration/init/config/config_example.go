/*
Copyright (c) 2023 the go-mongr8 Authors and Contributors
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
