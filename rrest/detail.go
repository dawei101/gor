package rrest

import (
	"database/sql"
	"net/http"

	"github.com/dawei101/gor/rhttp"
	"github.com/dawei101/gor/rrouter"
	"github.com/dawei101/gor/rsql"
)

type Detail struct {
	One
}

func (one *Detail) Handle(w http.ResponseWriter, r *http.Request) {
	id := rrouter.Var(r, "id")
	model := one.newModel()
	err := rsql.Model(model).Where(model.PK()+" = ?", id).Get()
	if err == sql.ErrNoRows {
		rhttp.NewErrResp(404, "no resource found", err.Error()).Json(w)
		return
	}
	if err != nil {
		rhttp.NewErrResp(500, "server error", err.Error()).Json(w)
		return
	}

	if one.ValidateModel != nil {
		err := one.ValidateModel(r, model)
		if err != nil {
			rhttp.FlushErr(w, r, err)
			return
		}
	}
	rhttp.NewResp(model).Json(w)
}
