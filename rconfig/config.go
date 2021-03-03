package rconfig

import (
	"errors"
	"fmt"
	"io/ioutil"
	"strings"
	"sync"

	"gopkg.in/yaml.v3"
)

const DefName = "default"
const DefFile = "config.yml"

var configs = sync.Map{}

// 最顶层设置
type RConfig struct {
	DevMode bool `yaml:"devMode"`
}

type Config struct {
	name    string
	data    map[string]interface{}
	RConfig RConfig
}

// register a configuration item
// Reg("default","config.yml")
func Reg(name, file string) error {
	_, ok := configs.Load(name)
	if ok {
		return nil
	}

	if !strings.HasSuffix(file, ".yml") {
		return errors.New(fmt.Sprintf("the file `%s` is not a yml", file))
	}

	byteVal, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}

	data := make(map[string]interface{})
	if err = yaml.Unmarshal(byteVal, &data); err != nil {
		return err
	}
	c := &Config{name: name, data: data}
	c.ValTo("", &c.RConfig)
	configs.Store(name, c)
	return nil
}

// get a configuration item
// Get("default")
func Get(name string) *Config {
	raw, ok := configs.Load(name)
	if !ok {
		fmt.Printf("the config `%s` is not exist\n", name)
		return nil
	}
	return raw.(*Config)
}

// set val to the struct
// ValTo("mysql",&m)
func (c *Config) ValTo(keyPath string, ptr interface{}) bool {
	var val interface{}
	val = c.data
	for _, key := range strings.Split(keyPath, ".") {
		if key == "" {
			continue
		}
		val = (val.(map[string]interface{}))[key]
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
	return yaml.Unmarshal(d, ptr) == nil
}

func (c *Config) IsDev() bool {
	return c.RConfig.DevMode
}

func DefConf() *Config {
	return Get(DefName)
}
