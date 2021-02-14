package rhttp

import (
	"fmt"
	"net/http"
)

type RespErr struct {
	Code int
	Msg  string
	Desc string
}

func NewRespErr(code int, msg, desc string) RespErr {
	return RespErr{
		Code: code,
		Msg:  msg,
		Desc: desc,
	}
}

func (e RespErr) Error() string {
	return fmt.Sprintf("err:(%s). desc:(%s)", e.Msg, e.Desc)
}

func (e RespErr) Flush(w http.ResponseWriter) {
	NewErrResp(e.Code, e.Msg, e.Desc).Flush(w)
}