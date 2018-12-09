package helper

import "time"

/*
	Utils and helper function
*/

// nowAsUnixMilli get the current unix time in millis
func NowAsUnixMilli() int64 {
	return time.Now().UnixNano() / 1e6
}
