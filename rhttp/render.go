package rhttp

import (
	"github.com/unrolled/render"
)

type Render struct {
	*render.Render
}

// 创建一个渲染器
//
// 具体使用参见：https://github.com/unrolled/render
//
// TODO config中设置好template 目录，直接使用
func NewRender(opts ...render.Options) *Render {
	return &Render{render.New(render.Options{
		Charset: "ISO-8859-1",
	})}
}
