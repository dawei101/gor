package rhttp

import (
	"net/http"

	"github.com/unrolled/render"
)

var renderer *render.Render = render.New()

type Resp struct {
	Status int         `json:"result"`
	Msg    string      `json:"msg"`
	Desc   string      `json:"desc"`
	Data   interface{} `json:"data"`
}

func NewResp(data interface{}) *Resp {
	return &Resp{
		Status: 200,
		Msg:    "ok",
		Desc:   "",
		Data:   data,
	}
}

func NewErrResp(status int, msg string, desc string) *Resp {
	return &Resp{
		Status: status,
		Msg:    msg,
		Desc:   desc,
		Data:   map[string]interface{}{},
	}
}

func (res *Resp) Json(w http.ResponseWriter) error {
	return renderer.JSON(w, 200, res)
}

func (res *Resp) Html(w http.ResponseWriter, tpl string) error {
	// use rrender
	//return renderer.HTML(w, res.Status, tpl, res)
	return res.Json(w)
}
