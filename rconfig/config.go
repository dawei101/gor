package rconfig

import (
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
	"strings"
	"sync"

	"github.com/dawei101/gor/base"
)

var configs = make(map[string]*Config)
var c_lock sync.RWMutex
var default_config_once sync.Once

// 最顶层设置
type RConfig struct {
	DevMode bool `yaml:"devMode"`
}

type Config struct {
	*base.Struct
	name    string
	RConfig *RConfig
}

//
// 加载并注册配置文件，并按name标识起来
//
func Reg(name, filePath string) *Config {
	if cfg, ok := configs[name]; ok {
		return cfg
	}

	data := make(map[string]interface{})
	rconfig := RConfig{}
	f, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	byteValue, _ := ioutil.ReadAll(f)

	cfname := strings.ToLower(filePath)
	yaml.Unmarshal(byteValue, &data)
	yaml.Unmarshal(byteValue, &rconfig)

	c_lock.Lock()
	cfg := &Config{Struct: base.NewStruct(data), name: name, RConfig: &rconfig}
	configs[name] = cfg
	c_lock.Unlock()
	return cfg
}

// 装在默认的配置文件，必须在使用前加载配置
func RegDefault(filePath string) *Config {
	return Reg("default", filePath)
}

// 获取默认配置
//
//如果未通过`LoadDefaultConfig(filePath)` 或 `LoadConfig(name, filePath)` 加载过配置, 将会panic
//		LoadDefaultConfig("./current/path/config.yml")
//		GetDefaultConfig()
// 一般情况下，我们只需要使用 default 这套config就足够用
func DefaultConfig() *Config {
	return GetConfig("default")
}

// 根据配置名获取配置
//		LoadConfig("myconfig", "./current/path/config.yml")
//		GetConfig("myconfig")
func GetConfig(name string) *Config {
	c_lock.Lock()
	defer c_lock.Unlock()
	cfg, ok := configs[name]
	if !ok {
		panic("no config named:" + name + " loaded, you need load first!!")
	}
	return cfg
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
