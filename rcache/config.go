package rcache

type FilecacheConfig struct {
}

type MemcacheConfig struct {
}

var config struct {
	Redis string          `yaml:"redis"`
	File  FilecacheConfig `yaml:"file"`
	Mem   MemcacheConfig  `yaml:"file"`
}
