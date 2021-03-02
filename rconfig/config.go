package rconfig

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
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
	RConfig RConfig
}

func init() {
	_, err := os.Stat(DefFile)
	if err != nil {
		panic(err)
	}

	Reg(DefName, DefFile)
}

// register a configuration item
// Reg("default","config.yml")
func Reg(name, file string) {
	_, ok := configs.Load(name)
	if ok {
		return
	}

	if !strings.HasSuffix(file, ".yml") {
		panic(fmt.Sprintf("the file `%s` is not a yml", file))
	}

	byteVal, err := ioutil.ReadFile(file)
	if err != nil {
		panic(err)
	}

	data := make(map[string]interface{})
	if err = yaml.Unmarshal(byteVal, &data); err != nil {
		panic(err)
	}
	c := &Config{name: name, data: data}
	c.ValTo("", &c.RConfig)
	configs.Store(name, c)
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
