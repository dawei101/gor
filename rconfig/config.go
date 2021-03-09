package rconfig

import (
	"errors"
	"fmt"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"strings"
	"sync"
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
	rConfig RConfig
}

func init() {
	err := Reg(DefName, DefFile)
	if err != nil {
		panic(err)
	}
}

// register a configuration item
// Reg("default","config.yml")
func Reg(name, file string) error {
	_, ok := configs.Load(name)
	if !ok {
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
		c.ValTo("", &c.rConfig)
		configs.Store(name, c)
	}

	return nil
}

// get a configuration item
// Get("default")
func Get(name string) *Config {
	raw, ok := configs.Load(name)
	if ok {
		return raw.(*Config)
	}
	return &Config{}
}

// set val to the struct
// ValTo("mysql",&m)
func (c *Config) ValTo(keyPath string, ptr interface{}) bool {
	var val interface{} = c.data
	if keyPath != "" {
		for _, key := range strings.Split(keyPath, ".") {
			val = (val.(map[string]interface{}))[key]
			if val == nil {
				return false
			}
		}
	}

	d, err := yaml.Marshal(val)
	if err != nil {
		return false
	}
	return yaml.Unmarshal(d, ptr) == nil
}

func (c *Config) IsDev() bool {
	return c.rConfig.DevMode
}

func DefConf() *Config {
	return Get(DefName)
}
