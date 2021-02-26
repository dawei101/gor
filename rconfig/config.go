package rconfig

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"strings"
	"sync"
)

var configs = sync.Map{}

// 最顶层设置
type RConfig struct {
	DevMode bool `yaml:"devMode"`
}

type Config struct {
	name    string
	data    map[string]interface{}
	RConfig *RConfig
}

const DefaultName = "default"

// register a configuration item
// Reg("default","config.yml")
func Reg(name, file string) {
	_, ok := configs.Load(name)
	if !ok {
		if !strings.HasSuffix(file, ".yml") {
			panic("the file is not a yml")
		}

		byteVal, err := ioutil.ReadFile(file)
		if err != nil {
			panic(err)
		}

		data := make(map[string]interface{})
		if err = yaml.Unmarshal(byteVal, &data); err != nil {
			panic(err)
		}
		configs.Store(name, &Config{name: name, data: data})
	}
}

// get a configuration item
// Get("default")
func Get(name string) *Config {
	raw, ok := configs.Load(name)
	if !ok {
		panic(fmt.Sprintf("the config `%s` is not exist", name))
	}
	return raw.(*Config)
}

// set val to the struct
// ValTo("mysql",&m)
func (c *Config) ValTo(keyPath string, ptr interface{}) bool {
	var val interface{}
	val = c.data
	for _, el := range strings.Split(keyPath, ".") {
		val = (val.(map[string]interface{}))[el]
		if val == nil {
			break
		}
	}
	if val == nil {
		return false
	}
	d, err := yaml.Marshal(val)
	if err != nil {
		return false
	}
	err = yaml.Unmarshal(d, ptr)
	return err == nil
}

func (c *Config) IsDev() bool {
	return c.RConfig.DevMode
}

func DefConf() *Config {
	return Get(DefaultName)
}
