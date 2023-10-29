package schema_interpreter

import (
	"fmt"
)

func (a Action) GetLiteralInstance(prefix string, isArrayItem bool) string {
	res := ""
	if !isArrayItem {
		res += fmt.Sprintf("%sAction", prefix)
	}

	res += "{\n"
	res += fmt.Sprintf(`ActionKey: "%s",%s`, a.ActionKey, "\n")
	res += fmt.Sprintf("SubActions: []%sSubAction{\n", prefix)
	// fill sub actions
	for _, sa := range a.SubActions {
		res += fmt.Sprintf("%s,\n", sa.GetLiteralInstance(prefix, true))
	}

	res += "},\n"
	res += "}"

	return res
}
