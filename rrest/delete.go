package rrest

import (
	"database/sql"
	"net/http"

	"github.com/dawei101/gor/rhttp"
	"github.com/dawei101/gor/rrouter"
	"github.com/dawei101/gor/rsql"
)

type Delete struct {
	One
}

func (one Delete) Handle(w http.ResponseWriter, r *http.Request) {
	id := rrouter.Var(r, "id")

	newm := newModel(one.Model)
	if err := rhttp.JsonBodyTo(r, newm); err != nil {
		rhttp.FlushErr(w, r, err)
		return
	}

	old := newModel(one.Model)
	err := rsql.Model(old).Where(old.PK()+" = ?", id).Get()
	if err == sql.ErrNoRows {
		rhttp.NewErrResp(404, "no resource found", err.Error()).Json(w)
		return
	}
	if err != nil {
		rhttp.NewErrResp(500, "server error", err.Error()).Json(w)
		return
	}

	if err := one.ValidateModel(r, newm); err != nil {
		rhttp.FlushErr(w, r, err)
		return
	}

	fields := forceZeroFields(old, newm)
	_, err = rsql.Model(newm).Where(old.PK()+"=?", id).Update(fields...)
	if err != nil {
		rhttp.NewErrResp(500, "update fail", err.Error()).Json(w)
		return
	}

	rhttp.NewResp(newm).Json(w)
}
