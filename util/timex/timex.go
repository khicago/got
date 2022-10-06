package timex

import (
	"github.com/khicago/got/util/typer"
	"golang.org/x/exp/constraints"
	"time"
)

type (
	Time = time.Time
)

// GetMSStamp
// returns the millisecond timestamp
func GetMSStamp() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

func GetFloat64Time(t time.Time) float64 {
	return float64(t.Nanosecond())/float64(time.Second) + float64(t.Unix())
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
