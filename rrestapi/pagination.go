package rrest

import (
	"encoding/json"
	"net/http"

	"roo.bo/rlib/rhttp"
)

type listPage struct {
	Items      []interface{} `json:"items"`
	Pagination pagination    `json:"pagination"`
}

type pagination struct {
	Page     int `json:"page"`
	PageSize int `json:"pageSize"`
	Total    int `json:"total"`
}

func paginationFromRequest(r *http.Request) *pagination {
	p := &pagination{
		Page:     0,
		PageSize: 20,
		Total:    0,
	}
	rhttp.JsonBodyTo(r, p)
	if p.PageSize < 5 {
		p.PageSize = 5
	}
	if p.Page < 0 {
		p.Page = 0
	}
	return p
}

func pageIt(items interface{}, page, pageSize int, total int) *listPage {
	p := &listPage{
		Pagination: pagination{
			Page:     page,
			PageSize: pageSize,
			Total:    total,
		},
	}

	bts, _ := json.Marshal(items)
	json.Unmarshal(bts, &p.Items)
	return p
}
