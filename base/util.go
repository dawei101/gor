package base

import (
	"time"
)

func MicroTimestamp(t time.Time) int64 {
	return t.UnixNano() / int64(time.Millisecond)
}

func Timestamp(t time.Time) int64 {
	return t.UnixNano() / int64(time.Second)
}
