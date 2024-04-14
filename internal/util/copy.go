/*
Copyright (c) 2023-present the go-mongr8 Authors and Contributors
[@see Authors file]

Licensed under the MIT License
(https://opensource.org/licenses/MIT)
*/
package util

import (
	"github.com/mohae/deepcopy"
)

func DeepCopy[T any](v T) T {
	return deepcopy.Copy(v).(T)
}
