package time

import (
	"time"
)

// Timestamp returns the current Unix timestamp in seconds.
func Timestamp() int64 {
	return time.Now().Unix()
}

// Unix returns the time for a Unix timestamp in seconds.
func Unix(sec int64) time.Time {
	return time.Unix(sec, 0)
}

// Elapse returns the time elapsed since the input timestamp (in seconds),
// truncated to whole seconds.
func Elapse(timestamp int64) time.Duration {
	return time.Since(Unix(timestamp)).Truncate(time.Second)
}

// Expired reports whether the current time is past the timestamp.
func Expired(timestamp int64) bool {
	return Timestamp() > timestamp
}
