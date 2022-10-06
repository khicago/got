package idgen

import (
	"github.com/khicago/got/util/timex"
)

func getIDPrefix(t timex.Time, digitsClearExtra int) int64 {
	prefix := int64(timex.GetFloat64Time(t)*(1<<(32-ControlIDDigits))) << ControlIDDigits // total move 32 bits
	if digitsClearExtra > 0 {
		prefix &= ^((1<<digitsClearExtra - 1) << ControlIDDigits) // erase more bits to free up space for custom segments
	}
	return prefix
}
