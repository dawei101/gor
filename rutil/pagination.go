package rutil

import (
	"github.com/dawei101/gor/rsql"
)

type Pagination struct {
	Page     int64 `json:"page"`
	PageSize int64 `json:"pageSize"`
	Total    int64 `json:"total"`
}

func NewPagination(page, pageSize, total int64) *Pagination {
	return &Pagination{page, pageSize, total}
}

func ReqPagination(r *http.Request) *Pagination {
	p := &Pagination{
		Page:     0,
		PageSize: 20,
		Total:    0,
	}
	JsonBodyTo(r, p)
	if p.PageSize < 5 {
		p.PageSize = 5
	}
	if p.Page < 0 {
		p.Page = 0
	}
	return p
}

func Paginate(r *http.Request, sql *rsql.Sql) (error, *rhttp.Pagination) {
	c, err := sql.Count()
	if err != nil {
		return err, nil
	}
	pag := rhttp.ReqPagination(r)
	sql.Offset(pag.Page * pag.PageSize).Limit(pag.PageSize)
	return NewPagination(pag.Page, pag.PageSize, c)
}

func ListPage(items interface{}, p *Pagination) map[string]interface{} {
	return map[string]interface{}{
		"item":       items,
		"pagination": p,
	}
}
