package rconfig

import (
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"strings"
	"sync"

	"github.com/dawei101/gor/base"
)

var configs = sync.Map{}

// 最顶层设置
type RConfig struct {
	DevMode bool `yaml:"devMode"`
}

type Config struct {
	*base.Struct
	name    string
	RConfig *RConfig
}

// register a configuration item
// Reg("default","config.yml")
func RegConfig(name, filePath string) {
	_, ok := configs.Load(name)
	if !ok {
		if !strings.HasSuffix(filePath, ".yml") {
			panic("the file is not a yml")
		}

		byteValue, err := ioutil.ReadFile(filePath)
		if err != nil {
			panic(err)
		}

		data := make(map[string]interface{})
		if err = yaml.Unmarshal(byteValue, &data); err != nil {
			panic(err)
		}
		configs.Store(name, &Config{Struct: base.NewStruct(data), name: name})
	}
}

// get a configuration item
// Get("default")
func GetConfig(name string) *Config {
	raw, ok := configs.Load(name)
	if !ok {
		RegConfig(name, name+".yml")
		raw, _ = configs.Load(name)
	}
	return raw.(*Config)
}

/*
将keyPath的设置复制到&value, 可以为空
	pageSize := 0
	c.ValueAssignTo("the.key.path.to.here", &pageSize, 10)

*/
func (c *Config) ValueAssignTo(keyPath string, valuePointer interface{}, default_val interface{}) {
	c.RLock()
	defer c.RUnlock()
	keyps := strings.Split(keyPath, ".")
	var val interface{}
	val = c.Raw
	for _, el := range keyps {
		val = (val.(map[string]interface{}))[el]
		if val == nil {
			break
		}
	}
	if val != nil {
		d, _ := yaml.Marshal(val)
		yaml.Unmarshal(d, valuePointer)
		return
	}
	if default_val != nil {
		d, _ := yaml.Marshal(default_val)
		yaml.Unmarshal(d, valuePointer)
	}
}

func (c *Config) IsDev() bool {
	return c.RConfig.DevMode
}

/*
将keyPath的设置复制到&value, 必须存在

	pageSize := 0
	c.MustValueAssignTo("the.key.path.to.here", &pageSize)

*/
func (c *Config) ValueMustAssignTo(keyPath string, valuePointer interface{}) {
	c.ValueAssignTo(keyPath, valuePointer, nil)
	if valuePointer == nil {
		panic("no value in struct")
	}
}

func DefaultConfig() *Config {
	return GetConfig("default")
}

func IsDev() bool {
	return DefaultConfig().IsDev()
}

func DataAssignTo(val interface{}) {
	DefaultConfig().DataAssignTo(val)
}

func GetInt(key string) (int, bool) {
	return DefaultConfig().GetInt(key)
}

func GetString(key string) (string, bool) {
	return DefaultConfig().GetString(key)
}

func GetFloat(key string) (float64, bool) {
	return DefaultConfig().GetFloat(key)
}

func GetSlice(key string) ([]interface{}, bool) {
	return DefaultConfig().GetSlice(key)
}

func Get(key string) (interface{}, bool) {
	return DefaultConfig().Get(key)
}

func GetStruct(key string) (*base.Struct, bool) {
	return DefaultConfig().GetStruct(key)
}
func ValueAssignTo(keyPath string, valuePointer interface{}, default_val interface{}) {
	DefaultConfig().ValueAssignTo(keyPath, valuePointer, default_val)
}

func ValueMustAssignTo(keyPath string, valuePointer interface{}) {
	DefaultConfig().ValueMustAssignTo(keyPath, valuePointer)
}
