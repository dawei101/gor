package rrestapi

import (
	"github.com/dawei101/gor/rsql"
	"github.com/dawei101/gor/rvalid"
	"github.com/julienschmidt/httprouter"
)

type Resource struct {
	model    rsql.IModel
	keyField string
	db       string
	parent   *Resource
	isMany   bool
}

func ModelToApi(model rsql.IModel, keyField string) *Resource {
	return &Resource{
		model:    model,
		keyField: model.PK(),
		db:       "default",
	}
}

func ModelToApiWithKeyField(model rsql.IModel, keyField string) *Resource {
	return &Resource{
		model:    model,
		keyField: keyField,
		db:       "default",
	}
}

func (rr *Resource) WithDB(db string) *Resource {
	r.db = db
	return r
}

func (rr *Resource) Router() *httprouter.Router {
	r.HandleFunc("/{key}", handleHttp(rr, detail)).Methods("GET")
	r.HandleFunc("/{key}", handleHttp(rr, update)).Methods("POST")
	r.HandleFunc("/", handleHttp(rr, create)).Methods("PUT")
	r.HandleFunc("/{key}", handleHttp(rr, del)).Methods("DELETE")
	r.HandleFunc("/", handleHttp(rr, search)).Methods("POST")
}

func (rr *Resource) Model() rsql.IModel {
	m := reflect.New(reflect.TypeOf(r.model).Elem()).Elem().Addr().Interface()
	return m.(rsql.IModel)
}

func (rr *Resource) Models() []rsql.IModel {
	ms := reflect.New(reflect.SliceOf(reflect.TypeOf(rr.model))).Elem().Addr().Interface()
	return ms.([]rsql.IModel)
}

func (rr *Resource) UrlKeyField() string {
	return r.keyField
}

func (rr *Resource) Detail(w http.ResponseWriter, r *http.Request) {
	m := rr.Model()
	sql, err := rr.beforeSaveSql(r, rsql.Use(rr.db).Model(m))
	if err != nil {
		return err, nil
	}
	err = sql.Get()
	return err, m
}

func (rr *Resource) Update(w http.ResponseWriter, r *http.Request) {
	old := rr.Model()
	changeto := rr.Model()
	rhttp.JsonBodyTo(r, changeto)
	if err := rvalid.ValidField(changeto); err != nil {
		return rhttp.NewRespErr(422, err.Error(), ""), nil
	}

	sql, err := rr.beforeSaveSql(r, rsql.Use(rr.db).Model(old))
	if err != nil {
		return err, nil
	}
	fields := forceZeroFields(old, changeto)

	if err := rr.beforeSaveSql(r, changeto); err != nil {
		return err, nil
	}
	sql, _ = rr.beforeSaveSql(r, rsql.Use(rr.db).Model(changeto))

	effected, err := sql.Update(fields...)
	if err == nil {
		if effected == 0 {
			return rhttp.NewRespErr(404, "make sure resource exists", ""), nil
		}
		err = rsql.Use(rr.db).Model(changeto).Get()
	}
	return err, changeto

}
func (rr *Resource) Create(w http.ResponseWriter, r *http.Request) {
	m := rr.Model()
	rhttp.JsonBodyTo(r, m)
	if err := rvalid.ValidField(m); err != nil {
		return rhttp.NewRespErr(422, err.Error(), ""), nil
	}
	if err := rr.beforeSaveSql(r, m); err != nil {
		return err, nil
	}

	id, err := rsql.Use(rr.db).Model(m).Create()
	rlog.Debug(r.Context(), "resouce created, id=", id)
	return err, m

}
func (rr *Resource) Delete(w http.ResponseWriter, r *http.Request) {
	m := rr.Model()
	sql, err := rr.beforeSaveSql(r, rsql.Use(rr.db).Model(m))
	if err != nil {
		return err, nil
	}
	effected, err := sql.Delete()
	if effected == 0 {
		err = rhttp.NewRespErr(404, "not found", "0 rows updated")
	}
	return err, nil
}

func (rr *Resource) List(w http.ResponseWriter, r *http.Request) {
	objs := rr.Models()
	pag := paginationFromRequest(r)

	sql, err := rr.beforeSaveSql(r, rsql.Use(rr.db).Model(objs))
	if err != nil {
		return err, nil
	}
	c, err := sql.Count()

	if err != nil {
		return err, nil
	}
	if c == 0 {
		return nil, pageIt([]int{}, pag.Page, pag.PageSize, int(c))
	}

	sql, _ = rr.beforeSaveSql(r, rsql.Use(rr.db).Model(objs))
	err = sql.Offset(pag.Page * pag.PageSize).Limit(pag.PageSize).All()
	if err != nil {
		return err, nil
	}
	data := pageIt(objs, pag.Page, pag.PageSize, int(c))
	return nil, data
}

func (rr *Resource) beforeSearchSql(r *http.Request, sql *rsql.Builder) {
	if rlog.Level == rlog.LEVEL_DEBUG {
		sql = sql.ShowSQL()
	}
	if rr.parent == nil {
		// one to many resource
		return sql.Where(fmt.Sprintf(" %s = ? ", rr.keyField), rrouter.Var(r, "key")), nil
	}
	if rr.isMany {
		// one to many resource
		pm, err := rr.parent.Model(r)
		if err != nil {
			return sql, err
		}
		if pm == nil {
			return sql, rhttp.NewRespErr(404, "no resource found", "")
		}

		pval, _ := getDBFieldValue(pm, rr.parentField)
		if err != nil {
			return sql, err
		}
		rlog.Debug(r.Context(), "get parent keyField=", rr.parentField, " value=", pval)
		wsql := sql.Where(fmt.Sprintf(" %s =? ", rr.relatedField), pval)

		rel_key := rrouter.Var(r, "rel_key")
		if rel_key != "" {
			wsql = wsql.Where(fmt.Sprintf(" %s =? ", rr.keyField), rel_key)
		}
		return wsql, nil
	} else {
		// one to one resource
		pm, err := rr.parent.Model(r)
		if err != nil {
			return sql, err
		}
		if pm == nil {
			return sql, rhttp.NewRespErr(404, "no resource found", "")
		}
		val, _ := getDBFieldValue(pm, rr.parentField)
		return sql.Where(fmt.Sprintf(" %s =? ", rr.relatedField), val), nil
	}
}

func (rr *Resource) beforeSaveSql(r *http.Request, row rsql.IModel) error {
	if rr.parent == nil {
		return nil
	}
	if rr.isMany {
		// one to many resource
		pm, err := rr.parent.Model(r)
		if err != nil {
			return err
		}
		if pm == nil {
			return rhttp.NewRespErr(404, "no resource found", "")
		}

		val, _ := getDBFieldValue(pm, rr.parentField)
		rlog.Debug(r.Context(), "get parent parentField=", rr.parentField, " value=", val)
		err = setDBFieldValue(row, rr.relatedField, val)
		if err != nil {
			msg := fmt.Sprintf("parent keyfield type may not corrent(require int(etc int)/string), get error:", err.Error())
			rlog.Debug(r.Context(), msg)
		}

		rel_key := rrouter.Var("rel_key")
		if rel_key != "" {
			err := setDBFieldValue(row, rr.keyField, val)
			if err != nil {
				msg := fmt.Sprintf("keyfield type may not corrent(require int(etc int)/string), get error:", err.Error())
				rlog.Debug(r.Context(), msg)
			}
		}
	} else {
		// one to one resource
		pm, err := rr.parent.Model(r)
		if err != nil {
			return err
		}
		if pm == nil {
			return rhttp.NewRespErr(404, "no resource found", "")
		}
		val, _ := getDBFieldValue(pm, rr.parentField)
		rlog.Debug(r.Context(), "try to set:", rr.relatedField, " to val:", val)
		setDBFieldValue(row, rr.relatedField, val)
	}
	return nil
}

func (rr *Resource) getRelations() []Resource {
	rt := reflect.TypeOf(rr.model).Elem()
	if rt.Kind() != reflect.Struct {
		panic("bad type")
	}
	all := []Resource{}
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
			res := ModelToApiWithKeyField(model, model.PK())
			all = append(all, newToOneResource(NewResource(name, model, model.PK(), conn), fields[1], rr, fields[0]))
		} else {
			all = append(all, newToManyResource(NewResource(name, model, model.PK(), conn), fields[1], rr, fields[0]))
		}
	}
	return all
}
