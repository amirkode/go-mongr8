/*
Copyright (c) 2023-present the go-mongr8 Authors and Contributors
[@see Authors file]

Licensed under the MIT License
(https://opensource.org/licenses/MIT)
*/
package common

import "github.com/amirkode/go-mongr8/version"

const (
	MigrationHistoryCollection = "mongr8_migration_history"
)

func Mongr8Version() string {
	return version.Version
}
