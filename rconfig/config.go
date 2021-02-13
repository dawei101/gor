package rconfig

import (
	"errors"
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

type Config struct {
	*base.Struct
	name string
}

//
// 加载并注册配置文件，并按name标识起来
//
func Reg(name, filePath string) *Config {
	if cfg, ok := configs[name]; ok {
		return cfg
	}

	data := make(map[string]interface{})
	f, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	byteValue, _ := ioutil.ReadAll(f)

	cfname := strings.ToLower(filePath)
	if cfname != strings.TrimRight(cfname, ".yml") {
		yaml.Unmarshal(byteValue, &data)
	} else {
		panic(errors.New("config file type not support:" + filePath))
	}

	c_lock.Lock()
	cfg := &Config{Struct: base.NewStruct(data), name: name}
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
