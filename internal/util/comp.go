package util

type Comparable[T any] interface {
	// the values are -1 (less than), 0 (same), or 1 (greater than)
	CompareTo(T) int
}

func InList[T Comparable[T]](what T, in []T) bool {
	for _, i := range in {
		if i.CompareTo(what) == 0 {
			return true
		}
	}

	return false
}
