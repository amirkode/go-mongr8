package migrator

type (
	PreData interface {}

	Migrator struct {
		contracts interface {
			PrepData() PreData
			ExecOps()
			Rollback()
		}
	}
)

func (m Migrator) OnError() {
	m.contracts.Rollback()
}