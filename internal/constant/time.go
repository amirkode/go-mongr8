package constant

import (
	"time"
)

func MinTimevalue() time.Time {
	return time.Date(1970, time.January, 1, 0, 0, 0, 0, time.UTC)
}