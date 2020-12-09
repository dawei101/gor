package rrest

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"roo.bo/rlib/rcontext"
	"roo.bo/rlib/rhttp"
	"roo.bo/rlib/rlog"
	"roo.bo/rlib/rsql"
)

type toOneResource struct {
	*resource
	parent       *resource
	parentField  string
	relatedField string
}

func newToOneResource(resource *resource, relatedField string, parent *resource, parentField string) *toOneResource {
	if _, err := getDBFieldValue(parent.model, parentField); err != nil {
		panic(fmt.Sprintf("no field=%s found in resource=%s", parentField, parent.name))
	}
	if _, err := getDBFieldValue(resource.model, relatedField); err != nil {
		panic(fmt.Sprintf("no field=%s found in resource=%s", relatedField, resource.name))
	}
	return &toOneResource{
		resource:     resource,
		parent:       parent,
		parentField:  parentField,
		relatedField: relatedField,
	}
}

func (rr *toOneResource) Route(router *mux.Router) {
	path_prefix := fmt.Sprintf("/%s", rr.name)
	r := router.PathPrefix(path_prefix).Subrouter()

	// todo related resource
	//r.HandleFunc("/docs", rr.Doc).Methods("GET")

	if PageDetail&rr.flag > 0 {
		r.HandleFunc("", handleHttp(rr, detail)).Methods("GET")
	}

	if PageUpdate&rr.flag > 0 {
		r.HandleFunc("", handleHttp(rr, update)).Methods("POST")
	}

	if PageCreate&rr.flag > 0 {
		r.HandleFunc("", handleHttp(rr, create)).Methods("PUT")
	}

	if PageDelete&rr.flag > 0 {
		r.HandleFunc("", handleHttp(rr, del)).Methods("DELETE")
	}
}
func (rr *toOneResource) lockWhere(r *http.Request, sql *rsql.Builder) (*rsql.Builder, error) {
	if rlog.Level == rlog.LEVEL_DEBUG {
		sql = sql.ShowSQL()
	}
	parent, err := rr.getParentModel(r)
	if err != nil {
		return sql, err
	}
	if parent == nil {
		return sql, rhttp.NewRespErr(404, "no resource found", "")
	}
	val, _ := getDBFieldValue(parent, rr.parentField)
	return sql.Where(fmt.Sprintf(" %s =? ", rr.relatedField), val), nil
}

func (rr *toOneResource) lockFields(r *http.Request, row rsql.IModel) error {
	parent, err := rr.getParentModel(r)
	if err != nil {
		return err
	}
	if parent == nil {
		return rhttp.NewRespErr(404, "no resource found", "")
	}
	val, _ := getDBFieldValue(parent, rr.parentField)
	rlog.Debug(r.Context(), "try to set:", rr.relatedField, " to val:", val)
	setDBFieldValue(row, rr.relatedField, val)
	return nil
}

func (rr *toOneResource) getParentModel(r *http.Request) (rsql.IModel, error) {
	ctx := rcontext.Context(r)
	p, ok := ctx.Get("parent")
	var err error
	if !ok {
		err, p = detail(rr.parent, r)
		rlog.Debug(r.Context(), "get parent:", p)
		if err != nil {
			return nil, err
		}
		ctx.Set("parent", p)
	}
	if p == nil {
		return nil, nil
	}
	return p.(rsql.IModel), nil
}
