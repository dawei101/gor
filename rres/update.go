package rest

import (
	"github.com/gorilla/mux"
	"net/http"
	"roo.bo/rlib"
	"roo.bo/rlib/rsql"
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
	vars := mux.Vars(c.R)
	id, _ := strconv.Atoi(vars["id"])

	req, _ := rlib.JsonBody(c.R)

	for k, v := range c.Force {
		req.Set(k, v)
	}

	err := rsql.Model(c.Model).Where("id = ?", id).Get()
	if err != nil {
		rlib.NewErrResp(-404, "not exists", err.Error()).Flush(c.W)
		return
	}

	req.Set("id", id)
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

	_, err = rsql.Model(c.Model).Update()
	if err != nil {
		rlib.NewErrResp(-422, "update fail", err.Error()).Flush(c.W)
		return
	}

	rlib.NewResp(c.Model).Flush(c.W)
}
