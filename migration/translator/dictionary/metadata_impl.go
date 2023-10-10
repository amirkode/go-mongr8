package dictionary

import (
	"github.com/amirkode/go-mongr8/collection/metadata"
)

func (dict Dictionary) GetOptions() *map[metadata.CollectionOption]interface{} {
	return dict.Collection.Collection().Spec().Options
}

// some other operations in the future