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
		if err := one.ValidateModel(r, old); err != nil {
			rhttp.FlushErr(w, r, err)
			return
		}
	}

	_, err = rsql.Model(old).Where(old.PK()+"=?", id).Delete()
	if err != nil {
		rhttp.NewErrResp(500, "update fail", err.Error()).Json(w)
		return
	}
	rhttp.NewResp(nil).Json(w)
}
