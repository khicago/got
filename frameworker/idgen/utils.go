package idgen

import "time"

func getIDPrefix(t time.Time, digitsClearExtra int) int64 {
	prefix := int64(
		(float64(t.Unix())+float64(t.Nanosecond())/float64(time.Second))*
			float64(1<<32-ControlIDDigits), // 18
	) << ControlIDDigits // 14
	if digitsClearExtra > 0 {
		prefix &= ^((1<<digitsClearExtra - 1) << ControlIDDigits)
	}
	return prefix
}
