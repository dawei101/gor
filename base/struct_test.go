package base

import (
	"encoding/json"
	"fmt"
	"path"
	"runtime"
	"testing"
)

func TestStruct(t *testing.T) {

	_, filename, _, _ := runtime.Caller(0)
	cf := path.Join(path.Dir(filename), "test_config.yml")
	fmt.Printf("config file set to: %s\n", cf)
	SetConfigFile(cf)
	DefaultRooboConfig()

	data := map[string]interface{}{
		"key1": 1,
		"key2": "string",
		"key3": float64(1.0),
		"key4": []string{"1", "2", "3"},
		"key5": map[string]interface{}{
			"key6": map[string]string{"k": "v"},
		},
	}
	bt, _ := json.Marshal(data)
	raw := map[string]interface{}{}
	json.Unmarshal(bt, &raw)
	s := &Struct{Raw: raw}

	if i, ok := s.GetInt("key1"); !ok || i != 1 {
		t.Errorf("struct getInt not correct")
	}

	if j, ok := s.GetString("key2"); !ok || j != "string" {
		t.Errorf("struct getString not correct")
	}

	if k, ok := s.GetFloat("key3"); !ok || k != float64(1) {
		t.Errorf("struct GetFloat not correct")
	}

	if L, ok := s.GetSlice("key4"); !ok || len(L) == 0 {
		t.Errorf("struct GetSlice not correct")
	}

	if _, ok := s.Get("key5"); !ok {
		t.Errorf("struct Get not correct")
	}
	if string(s.JsonMarshal()) != string(bt) {
		t.Errorf("struct JsonMarshal not correct")
	}
}
