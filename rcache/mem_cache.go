package rcache

import lru "github.com/hashicorp/golang-lru"

type MemCache struct {
	Cache *lru.ARCCache
}

var Mem = &MemCache{}

func init() {
	arc, _ := lru.NewARC(1024)
	Mem.Cache = arc
}

func (m *MemCache) Get(key string) (value interface{}, ok bool) {
	return m.Cache.Get(key)
}

func (m *MemCache) Set(key string, value interface{}) {
	m.Cache.Add(key, value)
}

func (m *MemCache) Del(key string) {
	m.Cache.Remove(key)
}
