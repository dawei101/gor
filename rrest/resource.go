package rrest

import (
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"github.com/gorilla/mux"

	"github.com/dawei101/gor/rlog"
	"github.com/dawei101/gor/rsql"
)

type resource struct {
	IResource
	name     string
	model    rsql.IModel
	keyField string
	db       string
	flag     PageFlag
}

func NewResource(name string, model rsql.IModel, keyField string, dbconfig string) *resource {
	_, err := getDBFieldValue(model, keyField)
	if err != nil {
		panic(fmt.Sprintf("no field=%s found in resource=%s", keyField, name))
	}
	if dbconfig == "" {
		dbconfig = "default"
	}
	rr := &resource{
		model:    model,
		name:     name,
		keyField: keyField,
		db:       dbconfig,
		flag:     PageAll,
	}
	rr.flag = PageAll
	return rr
}

func (rr *resource) getDB() *rsql.DB {
	return rsql.Use(rr.db)
}

func (rr *resource) generateModel() rsql.IModel {
	m := reflect.New(reflect.TypeOf(rr.model).Elem()).Elem().Addr().Interface()
	return m.(rsql.IModel)
}

func (rr *resource) generateModels() []rsql.IModel {
	ms := reflect.New(reflect.SliceOf(reflect.TypeOf(rr.model))).Elem().Addr().Interface()
	return ms.([]rsql.IModel)
}

func (rr *resource) Route(router *mux.Router) {

	path_prefix := fmt.Sprintf("/%s", rr.name)
	r := router.PathPrefix(path_prefix).Subrouter()

	r.HandleFunc("/docs", rr.doc).Methods("GET")

	if PageDetail&rr.flag > 0 {
		r.HandleFunc("/{key}", handleHttp(rr, detail)).Methods("GET")
	}

	if PageUpdate&rr.flag > 0 {
		r.HandleFunc("/{key}", handleHttp(rr, update)).Methods("POST")
	}

	if PageCreate&rr.flag > 0 {
		r.HandleFunc("/", handleHttp(rr, create)).Methods("PUT")
	}

	if PageDelete&rr.flag > 0 {
		r.HandleFunc("/{key}", handleHttp(rr, del)).Methods("DELETE")
	}

	if PageSearch&rr.flag > 0 {
		r.HandleFunc("/", handleHttp(rr, search)).Methods("POST")
	}

	if PageRelation&rr.flag > 0 {
		subr := r.PathPrefix("/{key}").Subrouter()
		for _, res := range rr.getRelations() {
			res.Route(subr)
		}
	}
}

func (rr *resource) doc(w http.ResponseWriter, r *http.Request) {
}

func (rr *resource) lockWhere(r *http.Request, sql *rsql.Builder) (*rsql.Builder, error) {
	if rlog.Level == rlog.LEVEL_DEBUG {
		sql = sql.ShowSQL()
	}
	return sql.Where(fmt.Sprintf(" %s = ? ", rr.keyField), mux.Vars(r)["key"]), nil
}

func (rr *resource) lockFields(r *http.Request, row rsql.IModel) error {
	return nil
}

func (rr *resource) getRelations() []IResource {
	rt := reflect.TypeOf(rr.model).Elem()
	if rt.Kind() != reflect.Struct {
		panic("bad type")
	}
	all := []IResource{}
	for i := 0; i < rt.NumField(); i++ {
		f := rt.Field(i)
		rel := f.Tag.Get("relation")
		if rel == "" {
			continue
		}
		if rt.Kind() != reflect.Struct && f.Type.Kind() != reflect.Slice {
			panic("model relation field is not struct:" + f.Name)
		}
		fields := strings.Split(rel, ",")
		if len(fields) != 2 {
			panic("relation field set is not correct:" + f.Name)
		}
		conn := f.Tag.Get("connection")
		model := reflect.New(f.Type.Elem()).Elem().Addr().Interface().(rsql.IModel)
		name := f.Tag.Get("json")
		if rt.Kind() == reflect.Struct {
			all = append(all, newToOneResource(NewResource(name, model, model.PK(), conn), fields[1], rr, fields[0]))
		} else {
			all = append(all, newToManyResource(NewResource(name, model, model.PK(), conn), fields[1], rr, fields[0]))
		}
	}
	return all
}
