package rcontext

import (
	"context"
	"net/http"

	"roo.bo/rlib/base"
)

const (
	__context_k = "__context_k"
)

func ReqContext(r *http.Request) *base.Struct {
	return CtxContext(r.Context())
}

func CtxContext(ctx context.Context) *base.Struct {
	if ctx == nil {
		return nil
	}
	c_i := ctx.Value(__context_k)
	if c_i == nil {
		return nil
	}
	return c_i.(*base.Struct)
}

func CtxGet(ctx context.Context, key string) interface{} {
	st := CtxContext(ctx)
	if st != nil {
		if v, ok := st.Get(key); ok {
			return v
		}
	}
	return nil
}

func CtxGetOk(ctx context.Context, key string) (interface{}, bool) {
	st := CtxContext(ctx)
	if st != nil {
		v, ok := st.Get(key)
		return v, ok
	}
	return nil, false
}

func CtxSet(ctx context.Context, key string, val interface{}) {
	st := CtxContext(ctx)
	if st != nil {
		st.Set(key, val)
	}
}

func Get(r *http.Request, key string) interface{} {
	return CtxGet(r.Context(), key)
}

func GetOk(r *http.Request, key string) (interface{}, bool) {
	return CtxGetOk(r.Context(), key)
}

func Set(r *http.Request, key string, val interface{}) {
	CtxSet(r.Context(), key, val)
}

func Middleware_installRContext(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handler.ServeHTTP(w, RequestInstall(r))
	})
}

func RequestInstall(r *http.Request) *http.Request {
	if ReqContext(r) == nil {
		container := base.NewStruct(map[string]interface{}{})
		ctx := context.WithValue(r.Context(), __context_k, container)
		r = r.WithContext(ctx)
	}
	return r
}
