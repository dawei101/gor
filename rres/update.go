package rest

import (
	"github.com/dawei101/gor/rhttp"
	"github.com/dawei101/gor/rrouter"
	"github.com/dawei101/gor/rsql"
	"github.com/dawei101/gor/rvalid"
	"net/http"
	"strconv"
)

type Update struct {
	Validate map[string]Validator
	Force    map[string]string
	Model    interface{}
	R        *http.Request
	W        http.ResponseWriter
}

func (c Update) Parse() {
	vars := rrouter.Vars(c.R)
	id, _ := strconv.Atoi(vars["id"])

	req, _ := rhttp.JsonBody(c.R)

	for k, v := range c.Force {
		req.Set(k, v)
	}

	err := rsql.Model(c.Model).Where("id = ?", id).Get()
	if err != nil {
		rhttp.NewErrResp(-404, "not exists", err.Error()).Flush(c.W)
		return
	}

	req.Set("id", id)
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

	_, err = rsql.Model(c.Model).Update()
	if err != nil {
		rhttp.NewErrResp(-422, "update fail", err.Error()).Flush(c.W)
		return
	}

	rhttp.NewResp(c.Model).Flush(c.W)
}
