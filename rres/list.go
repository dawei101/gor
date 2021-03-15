package rest

import (
	"fmt"
	"net/http"
	"roo.bo/rlib"
	"roo.bo/rlib/rsql"
	"strconv"
	"strings"
)

const prePage = 1
const prePageSize = 5

type List struct {
	Filters []string
	Force   map[string]string
	Model   interface{}
	R       *http.Request
	W       http.ResponseWriter
}

type Data struct {
	Items    interface{} `json:"items"`
	Paginate Paginate    `json:"pagination"`
}

type Paginate struct {
	Page     int   `json:"page"`
	PageSize int   `json:"pageSize"`
	Total    int64 `json:"total"`
}

func (l List) Parse() {
	builder := rsql.Model(l.Model).ShowSQL()
	q := l.R.URL.Query()
	var orders []string
	//filter
	var filterMap = map[string]struct{}{}
	for _, v := range l.Filters {
		filterMap[v] = struct{}{}
	}

	for k, v := range q {
		_, ok := filterMap[k]
		if !ok {
			continue
		}
		keys := strings.SplitN(k, "__", 2)
		field := keys[0]

		op := "eq"
		if len(keys) == 2 {
			op = keys[1]
		}
		var value string
		value = v[0]
		if value == "" {
			continue
		}

		switch op {
		case "eq":
			builder.Where(fmt.Sprintf(" `%s` = ? ", field), value)
		case "like":
			builder.Where(fmt.Sprintf(" `%s` like ? ", field), fmt.Sprintf("%%%s%%", value))
		case "start":
			builder.Where(fmt.Sprintf(" `%s` like ? ", field), fmt.Sprintf("%s%%", value))
		case "gt":
			builder.Where(fmt.Sprintf(" `%s` > ? ", field), value)
		case "gte":
			builder.Where(fmt.Sprintf(" `%s` >= ? ", field), value)
		case "lt":
			builder.Where(fmt.Sprintf(" `%s` < ? ", field), value)
		case "lte":
			builder.Where(fmt.Sprintf(" `%s` <= ? ", field), value)
		case "in":
			builder.Where(fmt.Sprintf(" `%s` in (?) ", field), strings.Split(value, ","))
		case "between":
			b := strings.Split(value, ",")
			if len(b) != 2 {
				continue
			}
			builder.Where(fmt.Sprintf(" `%s` between ? and ? ", field), b[0], b[1])
		case "notBetween":
			b := strings.Split(value, ",")
			if len(b) != 2 {
				continue
			}
			builder.Where(fmt.Sprintf(" `%s` not between ? and ? ", field), b[0], b[1])
		case "sort":
			if !in(value, []string{"desc", "asc"}) {
				value = "asc"
			}
			orders = append(orders, fmt.Sprintf("%s %s", field, value))
		}
	}

	if 0 < len(orders) {
		builder.OrderBy(strings.Join(orders, ","))
	}

	//paginate
	page, _ := strconv.Atoi(q.Get("page"))
	if page < 1 {
		page = prePage
	}

	pageSize, _ := strconv.Atoi(q.Get("pageSize"))
	if pageSize < 5 {
		pageSize = prePageSize
	}

	for k, v := range l.Force {
		builder.Where(fmt.Sprintf("`%s` = ?", k), v)
	}

	var d = Data{
		Paginate: Paginate{
			Page:     page,
			PageSize: pageSize,
		},
	}
	d.Paginate.Total, _ = builder.Count()

	builder.Offset(pageSize * (page - 1)).Limit(pageSize)

	builder.All()

	d.Items = l.Model
	rlib.NewResp(d).Flush(l.W)
}

func in(target string, arr []string) bool {
	for _, ele := range arr {
		if target == ele {
			return true
		}
	}
	return false
}
