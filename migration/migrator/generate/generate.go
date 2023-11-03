package generate

import (
	"context"
	"fmt"
	"time"

	dt "internal/data_type"

	"github.com/amirkode/go-mongr8/migration/migrator"
	"github.com/amirkode/go-mongr8/migration/migrator/writer"
	"github.com/amirkode/go-mongr8/migration/option"
	si "github.com/amirkode/go-mongr8/migration/translator/mongodb/schema_interpreter"
)

func Run(ctx *context.Context, actions dt.Pair[[]si.Action, []si.Action]) error {
	if len(actions.First) == 0 {
		fmt.Println("Migration files are already up to date")
		return nil
	}

	migrationID := time.Now().Format("20060102_150405")

	migration := migrator.Migration{
		ID:   migrationID,
		Desc: option.GetMigrationOptionFromContext(ctx).Desc,
		Up:   actions.First,
		Down: actions.Second,
	}

	return writer.Write(migration)
}
