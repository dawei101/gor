package rhttp

import ()

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

var config struct {
	Server *ServerConfig           `yml:"server"`
	Client map[string]ClientConfig `yml:"client"`
}
