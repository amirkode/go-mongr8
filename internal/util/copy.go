package util

import (
	"github.com/mohae/deepcopy"
)

func DeepCopy[T any](v T) T {
	return deepcopy.Copy(v).(T)
}