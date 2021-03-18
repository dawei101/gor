package rrest

import (
	"fmt"
	"strings"

	"github.com/dawei101/gor/rsql"
)

type FilterOp struct {
	field       string
	op          string
	query_field string
}

func NewFilterOp(filter string) *FilterOp {
	z := strings.SplitN(filter, "__", 2)
	field := z[0]
	op := "eq"
	if len(z) == 2 {
		op = z[1]
	}
	return &FilterOp{
		field:       field,
		op:          op,
		query_field: field + "__" + op,
	}
}

func (f *FilterOp) Query(sql *rsql.Builder, value interface{}) {
	switch f.op {
	case "eq":
		sql.Where(fmt.Sprintf(" %s = ? ", f.field), value)
	case "like":
		sql.Where(fmt.Sprintf(" %s like ? ", f.field), fmt.Sprintf("%%%s%%", value))
	case "startwith":
		sql.Where(fmt.Sprintf(" %s like ? ", f.field), fmt.Sprintf("%s%%", value))
	case "gt":
		sql.Where(fmt.Sprintf(" %s > ? ", f.field), value)
	case "gte":
		sql.Where(fmt.Sprintf(" %s >= ? ", f.field), value)
	case "lt":
		sql.Where(fmt.Sprintf(" %s < ? ", f.field), value)
	case "lte":
		sql.Where(fmt.Sprintf(" %s <= ? ", f.field), value)
	case "in":
		sql.Where(fmt.Sprintf(" %s in (?) ", f.field), value)
	}
}
