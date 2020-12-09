package rrest

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strconv"

	"database/sql"
	"github.com/gorilla/mux"

	"roo.bo/rlib/rhttp"
	"roo.bo/rlib/rlog"
	"roo.bo/rlib/rsql"
)

var (
	mapper      = rsql.NewReflectMapper("db")
	display_log = false
)

type PageFlag int

type handler func(IResource, *http.Request) (error, interface{})

const (
	PageDetail PageFlag = 1 << iota
	PageUpdate
	PageCreate
	PageDelete
	PageSearch
	PageRelation

	PageAll PageFlag = -1 //位运算 ~0 = 111111111111
)

type IResource interface {
	Route(*mux.Router)
	lockWhere(*http.Request, *rsql.Builder) (*rsql.Builder, error)
	lockFields(*http.Request, rsql.IModel) error
	getDB() *rsql.DB
	generateModel() rsql.IModel
	generateModels() []rsql.IModel
}

func forceZeroFields(old, changeto rsql.IModel) (zeroFields []string) {
	nonzeros := map[string]int{}
	fields := mapper.FieldMap(reflect.ValueOf(old))
	for field, fval := range fields {
		if !rsql.IsZero(fval) {
			nonzeros[field] = 0
		}
	}

	fields = mapper.FieldMap(reflect.ValueOf(changeto))
	for field, fval := range fields {
		if _, ok := nonzeros[field]; ok && rsql.IsZero(fval) {
			zeroFields = append(zeroFields, field)
		}
	}
	return zeroFields
}

func handleHttp(rr IResource, h handler) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		err, data := h(rr, r)
		if err != nil {
			if rerr, ok := err.(rhttp.RespErr); ok {
				rerr.Flush(w)
				return
			}
			if err == sql.ErrNoRows {
				rhttp.NewErrResp(404, "not found", err.Error()).Flush(w)
				return
			}
			rhttp.NewErrResp(500, "server err", err.Error()).Flush(w)
		} else {
			rhttp.NewResp(data).Flush(w)
		}
	}
}

func getDBFieldValue(s interface{}, dbfield string) (interface{}, error) {
	fields := mapper.FieldMap(reflect.ValueOf(s))

	for field, fval := range fields {
		if field == dbfield {
			return fval.Interface(), nil
		}
	}
	rlog.Debug(context.Background(), fmt.Sprintf("get %s failed from instance:%v", dbfield, s))
	return nil, errors.New("get dbfield failed!")
}

func setDBFieldValue(s interface{}, dbfield string, toval interface{}) error {
	tostr := fmt.Sprintf("%v", toval)
	fields := mapper.FieldMap(reflect.ValueOf(s))

	for field, fval := range fields {
		if field == dbfield {
			switch fval.Kind() {
			case reflect.String:
				fval.SetString(tostr)
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
				i, err := strconv.ParseInt(tostr, 10, 64)
				if err != nil {
					return err
				}
				fval.SetInt(i)
			}
			return nil

		}
	}
	rlog.Debug(context.Background(), fmt.Sprintf("get %s failed from obj:%v", dbfield, s))
	return errors.New("set field failed")

}
