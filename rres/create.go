package rest

import (
	"github.com/dawei101/gor/base"
	"github.com/dawei101/gor/rhttp"
	"github.com/dawei101/gor/rsql"
	"github.com/dawei101/gor/rvalid"

	"net/http"
)

type Validator func(k string, form *base.Struct) (string, bool)

type Create struct {
	Validate map[string]Validator
	Force    map[string]string
	Model    interface{}
	R        *http.Request
	W        http.ResponseWriter
}

func (c Create) Parse() {
	req, _ := rhttp.JsonBody(c.R)

	for k, v := range c.Force {
		req.Set(k, v)
	}

	req.DataAssignTo(c.Model)

	if err := rvalid.FieldValid(c.Model); err != nil {
		rhttp.NewErrResp(-422, "", err.Error()).Flush(c.W)
		return
	}

	for k, f := range c.Validate {
		msg, ok := f(k, req)
		if !ok {
			rhttp.NewErrResp(-422, msg, "").Flush(c.W)
			return
		}
	}

	_, err := rsql.Model(c.Model).ShowSQL().Create()
	if err != nil {
		rhttp.NewErrResp(-422, "create fail", err.Error()).Flush(c.W)
		return
	}

	rhttp.NewResp(c.Model).Flush(c.W)
}
