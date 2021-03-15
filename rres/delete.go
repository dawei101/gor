package rest

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"roo.bo/rlib"
	"roo.bo/rlib/rsql"
	"strconv"
)

type Delete struct {
	Model interface{}
	Force map[string]string
	R     *http.Request
	W     http.ResponseWriter
}

func (d Delete) Parse() {
	vars := mux.Vars(d.R)
	id, _ := strconv.Atoi(vars["id"])
	q := rsql.Model(d.Model)

	for k, v := range d.Force {
		q.Where(fmt.Sprintf("`%s` = ?", k), v)
	}

	q.Where("id = ?", id)

	err := q.Get()
	if err != nil {
		rlib.NewErrResp(-404, "", err.Error()).Flush(d.W)
		return
	}

	_, err = q.Delete()

	if err != nil {
		rlib.NewErrResp(-422, "", err.Error()).Flush(d.W)
		return
	}

	rlib.NewResp("").Flush(d.W)
}
