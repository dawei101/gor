package rrest

import (
	"database/sql"
	"net/http"

	"github.com/dawei101/gor/rhttp"
	"github.com/dawei101/gor/rrouter"
	"github.com/dawei101/gor/rsql"
)

type Update struct {
	One
}

func (one Update) Handle(w http.ResponseWriter, r *http.Request) {
	id := rrouter.Var(r, "id")

	newm := one.newModel()
	if err := rhttp.JsonBodyTo(r, newm); err != nil {
		rhttp.FlushErr(w, r, err)
		return
	}

	old := one.newModel()
	err := rsql.Model(old).Where(old.PK()+" = ?", id).Get()
	if err == sql.ErrNoRows {
		rhttp.NewErrResp(404, "no resource found", err.Error()).Json(w)
		return
	}
	if err != nil {
		rhttp.NewErrResp(500, "server error", err.Error()).Json(w)
		return
	}

	if one.ValidateModel != nil {
		if err := one.ValidateModel(r, newm); err != nil {
			rhttp.FlushErr(w, r, err)
			return
		}
	}

	fields := forceZeroFields(old, newm)
	_, err = rsql.Model(newm).Where(old.PK()+"=?", id).Update(fields...)
	if err != nil {
		rhttp.NewErrResp(500, "update fail", err.Error()).Json(w)
		return
	}

	rhttp.NewResp(map[string]string{}).Json(w)
}
