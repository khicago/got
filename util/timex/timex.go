package timex

import "time"

// GetMSStamp
// returns the millisecond timestamp
func GetMSStamp() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}
