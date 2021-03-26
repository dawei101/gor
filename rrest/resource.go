package rrest

import (
	"net/http"
	"reflect"

	"github.com/dawei101/gor/rsql"
)

type ValidateModelFunc func(r *http.Request, model rsql.IModel) error

type Action interface {
	Handle(w http.ResponseWriter, r *http.Request)
}

type One struct {
	Model         interface{}
	ValidateModel ValidateModelFunc
}

func (one *One) newModel() rsql.IModel {
	m := reflect.New(reflect.TypeOf(one.Model).Elem()).Elem().Addr().Interface()
	return m.(rsql.IModel)
}

var (
	mapper = rsql.NewReflectMapper("db")
)

func forceZeroFields(old, changeto rsql.IModel) (zeroFields []string) {
	pk := old.PK()
	nonzeros := map[string]int{}
	fields := mapper.FieldMap(reflect.ValueOf(old))
	for field, fval := range fields {
		if !rsql.IsZero(fval) {
			nonzeros[field] = 0
		}
	}

	fields = mapper.FieldMap(reflect.ValueOf(changeto))
	for field, fval := range fields {
		if _, ok := nonzeros[field]; ok && rsql.IsZero(fval) && field != pk {
			zeroFields = append(zeroFields, field)
		}
	}
	return zeroFields
}
