package rrest

import (
	"net/http"
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
	if pageSize < 5 {
		pageSize = 5
	}
	return &Pagination{
		Page:     page,
		PageSize: pageSize,
		Total:    0,
	}
}

func (l *List) Handle(w http.ResponseWriter, r *http.Request) {
	sql := rsql.Model(l.Model)
	qs := r.URL.Query()

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
		err := l.BeforeQuery(r, sql)
		if err != nil {
			rhttp.FlushErr(w, r, err)
			return
		}
	}
	p := reqPagination(r)
	p.Total, _ = sql.Count()
	sql.Offset(p.PageSize * p.Page).Limit(p.PageSize)
	rhttp.NewResp(&Data{
		Items:      sql.All,
		Pagination: p,
	})
}
