package timex

import (
	"github.com/khicago/got/util/typer"
	"golang.org/x/exp/constraints"
	"time"
)

// GetMSStamp
// returns the millisecond timestamp
func GetMSStamp() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

func AfterMax(duration ...time.Duration) <-chan time.Time {
	return time.After(typer.SliceMax(duration))
}

func AfterMin(duration ...time.Duration) <-chan time.Time {
	return time.After(typer.SliceMin(duration))
}

func Second[T constraints.Integer | constraints.Float](v T) time.Duration {
	return time.Duration(float64(v) * float64(time.Second))
}
