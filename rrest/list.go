package rrest

import (
	"net/http"
	"reflect"
	"strconv"

	"github.com/dawei101/gor/rhttp"
	"github.com/dawei101/gor/rsql"
)

type List struct {
	Action
	Model       rsql.IModel
	Filters     []string
	BeforeQuery QueryBuilder
}

type QueryBuilder func(r *http.Request, sql *rsql.Builder) error
type Data struct {
	Items      interface{} `json:"items"`
	Pagination *Pagination `json:"pagination"`
}

type Pagination struct {
	Page     int   `json:"page"`
	PageSize int   `json:"pageSize"`
	Total    int64 `json:"total"`
}

func reqPagination(r *http.Request) *Pagination {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 0 {
		page = 0
	}
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("pageSize"))
	if pageSize < 20 {
		pageSize = 20
	}

	return &Pagination{
		Page:     page,
		PageSize: pageSize,
		Total:    0,
	}
}

func (l *List) newModels() (models_ptr interface{}) {
	return reflect.New(reflect.SliceOf(reflect.TypeOf(l.Model))).Elem().Addr().Interface()
}

func (l *List) Handle(w http.ResponseWriter, r *http.Request) {
	models := l.newModels()
	sql := rsql.Model(models)
	qs := r.URL.Query()
	var err error

	if l.Filters != nil {
		for _, f := range l.Filters {
			op := NewFilterOp(f)
			if val, ok := qs[op.field]; ok {
				op.Query(sql, val)
			}
			if val, ok := qs[op.query_field]; ok {
				op.Query(sql, val)
			}
		}
	}
	if l.BeforeQuery != nil {
		err = l.BeforeQuery(r, sql)
		if err != nil {
			rhttp.FlushErr(w, r, err)
			return
		}
	}
	p := reqPagination(r)
	p.Total, err = sql.Count()
	if err != nil {
		rhttp.FlushErr(w, r, err)
		return
	}
	if p.Total == 0 {
		rhttp.NewResp(&Data{
			Items:      []int{},
			Pagination: p,
		}).Json(w)
		return
	}

	sql.Offset(p.PageSize * p.Page).Limit(p.PageSize)
	err = sql.All()
	if err != nil {
		rhttp.FlushErr(w, r, err)
		return
	}
	rhttp.NewResp(&Data{
		Items:      models,
		Pagination: p,
	}).Json(w)
}
