package rhttp

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/dawei101/gor/base"
	"github.com/dawei101/gor/rcontext"
)

func JsonBodyTo(r *http.Request, v interface{}) error {
	structs, err := JsonBody(r)
	if err != nil {
		return err
	}
	structs.DataAssignTo(v)
	return nil
}

func JsonBody(r *http.Request) (st *base.Struct, err error) {
	ctype := r.Header.Get("Content-Type")

	if r.Method == "GET" {
		return nil, NewRespErr(422, "http GET method is not correct", "")
	}

	if !strings.Contains(ctype, "application/json") {
		return nil, NewRespErr(422, "body is not json", "")
	}

	ctx := rcontext.Context(r)
	if e_i, ok := ctx.Get("rbody_err"); ok {
		err = e_i.(error)
	}
	rbody, ok := ctx.Get("rbody")
	if !ok {
		body, err := ioutil.ReadAll(r.Body)
		r.Body.Close()
		r.Body = ioutil.NopCloser(bytes.NewBuffer([]byte("--data used already --")))
		val := map[string]interface{}{"_": "no data in body"}
		if err != nil || len(body) == 0 {
			ctx.Set("rbody", base.NewStruct(val))
		}
		if err != nil {
			ctx.Set("rbody_err", err)
			return nil, err
		}
		err = json.Unmarshal(body, &val)
		if err != nil {
			ctx.Set("rbody_err", err)
			return nil, err
		}
		st = base.NewStruct(val)
		ctx.Set("rbody", st)
	} else {
		st = rbody.(*base.Struct)
	}
	return st, err
}
