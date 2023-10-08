package action

type ActionType string

const (
	ActionCreateCollection ActionType = "create_collection"
	ActionCreateIndex      ActionType = "create_index"
	ActionCreateField      ActionType = "create_field"
	ActionDropCollection   ActionType = "drop_collection"
	ActionDropIndex        ActionType = "drop_index"
	ActionDropField        ActionType = "drop_field"
)
