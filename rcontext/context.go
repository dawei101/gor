package rcontext

import (
	"context"
	"net/http"

	"github.com/dawei101/gor/base"
	"github.com/dawei101/gor/rrouter"
)

func init() {
	rrouter.RegGlobalMiddleware(Middleware_installRContext)
}

const (
	__context_k = "__context_k"
)

func ReqCtx(r *http.Request) *base.Struct {
	return Ctx(r.Context())
}

func Ctx(ctx context.Context) *base.Struct {
	c_i := ctx.Value(__context_k)
	if c_i == nil {
		panic("need to install rcontext.Middleware_installRContext middleware before use it")
	}
	return c_i.(*base.Struct)
}

func Middleware_installRContext(handle http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handle(w, RequestInstall(r))
	}
}

func RequestInstall(r *http.Request) *http.Request {
	if ReqCtx(r) == nil {
		container := base.NewStruct(map[string]interface{}{})
		ctx := context.WithValue(r.Context(), __context_k, container)
		r = r.WithContext(ctx)
	}
	return r
}
