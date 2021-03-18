package rrest

import (
	"net/http"

	"github.com/dawei101/gor/rhttp"
	"github.com/dawei101/gor/rlog"
	"github.com/dawei101/gor/rsql"
	"github.com/dawei101/gor/rvalid"
)

type Create struct {
	One
}

func (one *Create) Handle(w http.ResponseWriter, r *http.Request) {

	m := newModel(one.Model)
	if err := rhttp.JsonBodyTo(r, m); err != nil {
		rhttp.FlushErr(w, r, err)
		return
	}
	if err := rvalid.ValidField(m); err != nil {
		rhttp.FlushErr(w, r, err)
		return
	}

	if err := one.ValidateModel(r, m); err != nil {
		rhttp.FlushErr(w, r, err)
		return
	}

	id, err := rsql.Model(m).Create()
	if err != nil {
		rhttp.FlushErr(w, r, err)
	}
	rlog.Debug(r.Context(), "resouce created, id=", id)
	rhttp.NewResp(m).Json(w)
}
