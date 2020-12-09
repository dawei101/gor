package rcache

/*
封装时间高性能的local cache功能。
直接使用即可
*/

import (
	gocache "github.com/patrickmn/go-cache"
	"time"
)

var local_c *gocache.Cache

type cache_data_holder struct {
	data interface{}
}

func init() {
	local_c = gocache.New(5*time.Minute, 10*time.Minute)
}

//
// 设置进程内存级缓存
//
func LocalCacheSet(key string, value interface{}, d time.Duration) {
	local_c.Set(key, &cache_data_holder{data: value}, d)
}

//
// 获取进程内存级缓存
//
func LocalCacheGet(key string) (value interface{}, ok bool) {
	h, ok := local_c.Get(key)
	if ok {
		value = h.(*cache_data_holder).data
	}
	return value, ok
}
