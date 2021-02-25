package rhttp

import (
	"github.com/unrolled/render"
)

type RenderConfig struct {
	Directory string `yml:"directory"`
	Layout    string `yml:"layout"`
}

type ServerConfig struct {
	Render         *RenderConfig `yml:"render"`
	ReadTimeoutMs  int64         `yml:"readtimeoutMs"`
	WriteTimeoutMs int64         `yml:"writetimeoutMs"`
	Middlewares    []string      `yml:"middlewares"`
}

type ClientConfig struct {
	BaseUrl string `yml:"baseUrl"`
}

func NewRender(opts ...render.Options) *Render {
	return &Render{render.New(render.Options{
		Charset:    "ISO-8859-1",
		Directory:  "",
		FileSystem: &LocalFileSystem{},
	})}
}

var Config struct {
	Server *ServerConfig           `yml:"server"`
	Client map[string]ClientConfig `yml:"client"`
}

func init() {
	rconfig.ValueAssignTo("rhttp", &Config, nil)
}
