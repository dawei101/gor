package base

import (
	"encoding/json"
	"sync"
)

type Struct struct {
	Raw map[string]interface{}
	sync.RWMutex
}

func NewStruct(data map[string]interface{}) *Struct {
	return &Struct{Raw: data}
}

func (r *Struct) DataAssignTo(val interface{}) {
	r.RLock()
	defer r.RUnlock()
	d, _ := json.Marshal(r.Raw)
	json.Unmarshal(d, val)
}

func (r *Struct) GetInt(key string) (int, bool) {
	r.RLock()
	defer r.RUnlock()
	if v, ok := r.Raw[key]; ok {
		if tv, tok := v.(float64); tok {
			return int(tv), true
		}
	}
	return 0, false
}

func (r *Struct) GetBool(key string) (interface{}, bool) {
	r.RLock()
	defer r.RUnlock()
	if v, ok := r.Raw[key]; ok {
		if tv, tok := v.(float64); tok {
			return int(tv), true
		}
	}
	return 0, false
}

func (r *Struct) GetString(key string) (string, bool) {
	r.RLock()
	defer r.RUnlock()
	if v, ok := r.Raw[key]; ok {
		if tv, tok := v.(string); tok {
			return tv, true
		}
	}
	return "", false
}

func (r *Struct) GetFloat(key string) (float64, bool) {
	r.RLock()
	defer r.RUnlock()
	if v, ok := r.Raw[key]; ok {
		if tv, tok := v.(float64); tok {
			return tv, true
		}
	}
	return 0, false
}

func (r *Struct) GetSlice(key string) ([]interface{}, bool) {
	r.RLock()
	defer r.RUnlock()
	if v, ok := r.Raw[key]; ok {
		if tv, tok := v.([]interface{}); tok {
			return tv, true
		}
	}
	return []interface{}{}, false
}

func (r *Struct) Get(key string) (interface{}, bool) {
	r.RLock()
	defer r.RUnlock()
	val, ok := r.Raw[key]
	return val, ok
}

func (r *Struct) GetStruct(key string) (*Struct, bool) {
	r.RLock()
	defer r.RUnlock()
	if v, ok := r.Raw[key]; ok {
		if tv, tok := v.(map[string]interface{}); tok {
			return &Struct{Raw: tv}, true
		}
	}
	return nil, false
}

func (r *Struct) Set(key string, val interface{}) {
	r.Lock()
	defer r.Unlock()
	r.Raw[key] = val
}

func (r *Struct) Keys() []string {
	r.Lock()
	defer r.Unlock()
	keys := []string{}
	for k, _ := range r.Raw {
		keys = append(keys, k)
	}
	return keys
}

func (r *Struct) JsonMarshal() []byte {
	r.RLock()
	defer r.RUnlock()
	d, _ := json.Marshal(r.Raw)
	return d
}
