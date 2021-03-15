package rest

import (
	"fmt"
	"github.com/dawei101/gor/rhttp"
	"github.com/dawei101/gor/rrouter"
	"github.com/dawei101/gor/rsql"
	"net/http"
	"strconv"
)

type Detail struct {
	Model interface{}
	Force map[string]string
	R     *http.Request
	W     http.ResponseWriter
}

func (d Detail) Parse() {
	vars := rrouter.Vars(d.R)
	id, _ := strconv.Atoi(vars["id"])
	q := rsql.Model(d.Model)
	for k, v := range d.Force {
		q.Where(fmt.Sprintf("`%s` = ?", k), v)
	}
	err := rsql.Model(d.Model).Where("id = ?", id).Get()
	if err != nil {
		rhttp.NewErrResp(-404, "", err.Error()).Flush(d.W)
		return
	}
	rhttp.NewResp(d.Model).Flush(d.W)
}
