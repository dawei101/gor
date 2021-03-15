package rest

import (
	"net/http"
	"roo.bo/rlib"
	"roo.bo/rlib/base"
	"roo.bo/rlib/rsql"
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
	req, _ := rlib.JsonBody(c.R)

	for k, v := range c.Force {
		req.Set(k, v)
	}

	req.DataAssignTo(c.Model)

	if err := rlib.FieldValid(c.Model); err != nil {
		rlib.NewErrResp(-422, "", err.Error()).Flush(c.W)
		return
	}

	for k, f := range c.Validate {
		msg, ok := f(k, req)
		if !ok {
			rlib.NewErrResp(-422, msg, "").Flush(c.W)
			return
		}
	}

	_, err := rsql.Model(c.Model).ShowSQL().Create()
	if err != nil {
		rlib.NewErrResp(-422,"create fail",err.Error()).Flush(c.W)
		return
	}


	rlib.NewResp(c.Model).Flush(c.W)
}
