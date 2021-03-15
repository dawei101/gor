package rest

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"roo.bo/rlib"
	"roo.bo/rlib/rsql"
	"strconv"
)

type Detail struct {
	Model interface{}
	Force map[string]string
	R     *http.Request
	W     http.ResponseWriter
}

func (d Detail) Parse() {
	vars := mux.Vars(d.R)
	id, _ := strconv.Atoi(vars["id"])
	q := rsql.Model(d.Model)
	for k, v := range d.Force {
		q.Where(fmt.Sprintf("`%s` = ?", k), v)
	}
	err := rsql.Model(d.Model).Where("id = ?", id).Get()
	if err != nil {
		rlib.NewErrResp(-404, "", err.Error()).Flush(d.W)
		return
	}
	rlib.NewResp(d.Model).Flush(d.W)
}
