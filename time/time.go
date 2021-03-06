package time

import (
	"time"
)

// Timestamp return timestamp in second
func Timestamp() int64 {
	return time.Now().Unix()
}

// Unix return time from timestamp
func Unix(sec int64) time.Time {
	return time.Unix(sec, 0)
}

// Elapse time since the input timestamp
func Elapse(timestamp int64) time.Duration {
	return time.Duration(Timestamp()-timestamp) * time.Second
}

// Expired if timestamp is expired
func Expired(timestamp int64) bool {
	return Timestamp() > timestamp
}
