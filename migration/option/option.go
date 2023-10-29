package option

import (
	"context"
	"flag"
	"fmt"
)

const (
	MigrationOptionKey = "migration-option"

	MigrationOptionArgUseSortedSchema     = "use-sorted-schema"
	MigrationOptionArgUseForceConversion  = "use-force-conversion"
	MigrationOptionArgUseSchemaValidation = "use-schema-validation"
	MigrationOptionArgDesc                = "desc"
)

type (
	MigrationOption struct {
		UseSortedSchema     bool
		UseForceConversion  bool
		UseSchemaValidation bool
		Desc                string
	}
)

func GetMigrationOptionFromArgs() MigrationOption {
	flag.VisitAll(func(f *flag.Flag) {
		fmt.Printf("%s: %s\n", f.Name, f.Value)
	})
	opt := MigrationOption{}
	flag.BoolVar(&opt.UseSortedSchema, MigrationOptionArgUseSortedSchema, false, "Define option for Sorted MongoDb Schema")
	flag.BoolVar(&opt.UseForceConversion, MigrationOptionArgUseForceConversion, false, "Define option for forced conversion on migration")
	flag.BoolVar(&opt.UseSchemaValidation, MigrationOptionArgUseSchemaValidation, false, "Define option for Schema Validation on migration")
	flag.StringVar(&opt.Desc, MigrationOptionArgDesc, "", "Define option for Schema Validation on migration")
	flag.Parse()

	return opt
}

func GetMigrationOptionFromContext(ctx *context.Context) MigrationOption {
	if ctx == nil {
		panic("Context must be provided to get the option")
	}

	opt := (*ctx).Value(MigrationOptionKey).(MigrationOption)

	return opt
}
