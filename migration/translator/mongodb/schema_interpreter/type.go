package schema_interpreter

type SubActionType string

// make sure constant name is exactly same as it's value
const (
	SubActionTypeCreateCollection SubActionType = "SubActionTypeCreateCollection"
	SubActionTypeCreateIndex      SubActionType = "SubActionTypeCreateIndex"
	SubActionTypeCreateField      SubActionType = "SubActionTypeCreateField"
	SubActionTypeConvertField     SubActionType = "SubActionTypeConvertField"
	SubActionTypeDropCollection   SubActionType = "SubActionTypeDropCollection"
	SubActionTypeDropIndex        SubActionType = "SubActionTypeDropIndex"
	SubActionTypeDropField        SubActionType = "SubActionTypeDropField"
)

func (sat SubActionType) ToString() string {
	return string(sat)
}
