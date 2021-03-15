package rhttp

import (
	"net/http"
	"strings"
)

type RespErr struct {
	Status int
	Msg    string
	Desc   string
}

func NewRespErr(status int, msg, desc string) RespErr {
	return RespErr{
		Status: status,
		Msg:    msg,
		Desc:   desc,
	}
}

func (e RespErr) Error() string {
	return e.Msg
}

func (e RespErr) Flush(w http.ResponseWriter, r *http.Request) {
	res := NewErrResp(e.Status, e.Msg, e.Desc)

	accept := strings.ToLower(r.Header.Get("Accept"))
	for _, ctype := range strings.Split(strings.Split(accept, ";")[0], ",") {
		if strings.HasSuffix(strings.TrimRight(ctype, " "), "/json") {
			res.Json(w)
			return
		}
	}
	res.Html(w, string(res.Status))
}

func FlushErr(w http.ResponseWriter, r *http.Request, err error) {
	if rerr, ok := err.(RespErr); ok {
		rerr.Flush(w, r)
		return
	}
	NewRespErr(500, "server error", err.Error()).Flush(w, r)
}
