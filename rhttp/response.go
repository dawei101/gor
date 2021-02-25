package rhttp

import (
	"encoding/json"
	"net/http"
)

type Resp struct {
	Status int         `json:"result"`
	Msg    string      `json:"msg"`
	Desc   string      `json:"desc"`
	Data   interface{} `json:"data"`
}

func NewResp(data interface{}) *Resp {
	return &Resp{
		Status: 0,
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
		Data:   nil,
	}
}

func (res *Resp) Flush(w http.ResponseWriter) error {
	body, err := json.Marshal(res)
	if err != nil {
		return err
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(body)
	return nil
}

func (rr *Resp) DataAssignTo(data interface{}) {
	d, _ := json.Marshal(rr.Data)
	json.Unmarshal(d, data)
}

func (res *Resp) FlushHtml(tpl string, w http.ResponseWriter) error {
	// use rrender
}
