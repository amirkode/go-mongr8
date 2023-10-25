package schema_interpreter

type SubActionType string

const (
	SubActionTypeCreateCollection SubActionType = "create_collection"
	SubActionTypeCreateIndex      SubActionType = "create_index"
	SubActionTypeCreateField      SubActionType = "create_field"
	SubActionTypeConvertField     SubActionType = "convert_field"
	SubActionTypeDropCollection   SubActionType = "drop_collection"
	SubActionTypeDropIndex        SubActionType = "drop_index"
	SubActionTypeDropField        SubActionType = "drop_field"
)
