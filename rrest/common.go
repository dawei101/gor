package rrestapi

import (
	"github.com/dawei101/gor/rsql"
)

var (
	mapper = rsql.NewReflectMapper("db")
)

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
