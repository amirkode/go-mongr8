package api_interpreter

type ActionType string

const (
	// Action type consists two conditions
	ActionAdd    ActionType = "add"
	ActionRemove ActionType = "remove"

	// Comment for now
	// ActionCreateCollection ActionType = "create_collection"
	// ActionCreateIndex      ActionType = "create_index"
	// ActionCreateField      ActionType = "create_field"
	// ActionDropCollection   ActionType = "drop_collection"
	// ActionDropIndex        ActionType = "drop_index"
	// ActionDropField        ActionType = "drop_field"
)

type SubActionType string

const (
	SubActionCreateCollection SubActionType = "Create Collection"
	SubActionInsertOne        SubActionType = "Insert One"
	SubActionCreateIndex      SubActionType = "Create Index"
	SubActionUpdateManySet    SubActionType = "Update Many Set"
	SubActionUpdateManyUnset  SubActionType = "Update Many Unset"
)
