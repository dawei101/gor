package rrest

import (
	"net/http"

	"roo.bo/rlib/rhttp"
	"roo.bo/rlib/rlog"
	"roo.bo/rlib/rvalid"
)

func detail(rr IResource, r *http.Request) (error, interface{}) {
	m := rr.generateModel()
	sql, err := rr.lockWhere(r, rr.getDB().Model(m))
	if err != nil {
		return err, nil
	}
	err = sql.Get()
	return err, m
}

func update(rr IResource, r *http.Request) (error, interface{}) {
	old := rr.generateModel()
	changeto := rr.generateModel()
	rhttp.JsonBodyTo(r, changeto)
	if err := rvalid.ValidField(changeto); err != nil {
		return rhttp.NewRespErr(422, err.Error(), ""), nil
	}

	sql, err := rr.lockWhere(r, rr.getDB().Model(old))
	if err != nil {
		return err, nil
	}
	fields := forceZeroFields(old, changeto)

	if err := rr.lockFields(r, changeto); err != nil {
		return err, nil
	}
	sql, _ = rr.lockWhere(r, rr.getDB().Model(changeto))

	effected, err := sql.Update(fields...)
	if err == nil {
		if effected == 0 {
			return rhttp.NewRespErr(404, "make sure resource exists", ""), nil
		}
		err = rr.getDB().Model(changeto).Get()
	}
	return err, changeto
}

func create(rr IResource, r *http.Request) (error, interface{}) {
	m := rr.generateModel()
	rhttp.JsonBodyTo(r, m)
	if err := rvalid.ValidField(m); err != nil {
		return rhttp.NewRespErr(422, err.Error(), ""), nil
	}
	if err := rr.lockFields(r, m); err != nil {
		return err, nil
	}

	id, err := rr.getDB().Model(m).Create()
	rlog.Debug(r.Context(), "resouce created, id=", id)
	return err, m
}

func del(rr IResource, r *http.Request) (error, interface{}) {
	m := rr.generateModel()
	sql, err := rr.lockWhere(r, rr.getDB().Model(m))
	if err != nil {
		return err, nil
	}
	effected, err := sql.Delete()
	if effected == 0 {
		err = rhttp.NewRespErr(404, "not found", "0 rows updated")
	}
	return err, nil
}

func search(rr IResource, r *http.Request) (error, interface{}) {
	objs := rr.generateModels()
	pag := paginationFromRequest(r)

	sql, err := rr.lockWhere(r, rr.getDB().Model(objs))
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

	sql, _ = rr.lockWhere(r, rr.getDB().Model(objs))
	err = sql.Offset(pag.Page * pag.PageSize).Limit(pag.PageSize).All()
	if err != nil {
		return err, nil
	}
	data := pageIt(objs, pag.Page, pag.PageSize, int(c))
	return nil, data
}
