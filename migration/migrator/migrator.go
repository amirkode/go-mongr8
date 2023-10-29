package migrator

import (
	si "github.com/amirkode/go-mongr8/migration/translator/mongodb/schema_interpreter"
)

type (
	// Migration entity
	Migration struct {
		ID   string
		Desc string
		Up   []si.Action
		Down []si.Action
	}

	MigratorIf interface {
		PrepData()
		ExecOps()
		Rollback()
	}

	Migrator struct {
		MigratorIf
	}
)

func (m Migrator) OnError() {
	m.Rollback()
}
