package rrest

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	"roo.bo/rlib/rcontext"
	"roo.bo/rlib/rhttp"
	"roo.bo/rlib/rlog"
	"roo.bo/rlib/rsql"
)

type toManyResource struct {
	*resource
	parent       *resource
	parentField  string
	relatedField string
}

func newToManyResource(resource *resource, relatedField string, parent *resource, parentField string) *toManyResource {
	if _, err := getDBFieldValue(parent.model, parentField); err != nil {
		panic(fmt.Sprintf("no field=%s found in resource=%s", parentField, parent.name))
	}
	if _, err := getDBFieldValue(resource.model, relatedField); err != nil {
		panic(fmt.Sprintf("no field=%s found in resource=%s", relatedField, resource.name))
	}
	return &toManyResource{
		resource:     resource,
		parent:       parent,
		parentField:  parentField,
		relatedField: relatedField,
	}
}

func (rr *toManyResource) Route(router *mux.Router) {
	rlog.Debug(context.Background(), "install resouce:", rr.name)

	path_prefix := fmt.Sprintf("/%s", rr.name)
	r := router.PathPrefix(path_prefix).Subrouter()

	// todo related resource
	//r.HandleFunc("/docs", rr.Doc).Methods("GET")

	if PageDetail&rr.flag > 0 {
		r.HandleFunc("/{rel_key}", handleHttp(rr, detail)).Methods("GET")
	}

	if PageUpdate&rr.flag > 0 {
		r.HandleFunc("/{rel_key}", handleHttp(rr, update)).Methods("POST")
	}

	if PageCreate&rr.flag > 0 {
		r.HandleFunc("/", handleHttp(rr, create)).Methods("PUT")
	}

	if PageDelete&rr.flag > 0 {
		r.HandleFunc("/{rel_key}", handleHttp(rr, del)).Methods("DELETE")
	}
}

func (rr *toManyResource) lockWhere(r *http.Request, sql *rsql.Builder) (*rsql.Builder, error) {
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

	pval, _ := getDBFieldValue(parent, rr.parentField)
	if err != nil {
		return sql, err
	}
	rlog.Debug(r.Context(), "get parent keyField=", rr.parentField, " value=", pval)
	wsql := sql.Where(fmt.Sprintf(" %s =? ", rr.relatedField), pval)

	rel_key := mux.Vars(r)["rel_key"]
	if rel_key != "" {
		wsql = wsql.Where(fmt.Sprintf(" %s =? ", rr.keyField), rel_key)
	}
	return wsql, nil
}

func (rr *toManyResource) lockFields(r *http.Request, row rsql.IModel) error {
	parent, err := rr.getParentModel(r)
	if err != nil {
		return err
	}
	if parent == nil {
		return rhttp.NewRespErr(404, "no resource found", "")
	}

	val, _ := getDBFieldValue(parent, rr.parentField)
	rlog.Debug(r.Context(), "get parent parentField=", rr.parentField, " value=", val)
	err = setDBFieldValue(row, rr.relatedField, val)
	if err != nil {
		msg := fmt.Sprintf("parent keyfield type may not corrent(require int(etc int)/string), get error:", err.Error())
		rlog.Debug(r.Context(), msg)
	}

	rel_key := mux.Vars(r)["rel_key"]
	if rel_key != "" {
		err := setDBFieldValue(row, rr.keyField, val)
		if err != nil {
			msg := fmt.Sprintf("keyfield type may not corrent(require int(etc int)/string), get error:", err.Error())
			rlog.Debug(r.Context(), msg)
		}
	}
	return nil
}

func (rr *toManyResource) getParentModel(r *http.Request) (rsql.IModel, error) {
	ctx := rcontext.Context(r)
	p, ok := ctx.Get("parent")
	if !ok {
		err, p := detail(rr.parent, r)
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
