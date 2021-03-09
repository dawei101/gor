package rcache

import "time"

type Cache interface {
	Get(key string) (value interface{}, ok bool)
	Set(key string, value interface{}, d time.Duration)
	Del(key string)
}
