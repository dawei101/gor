package rcontext

import (
	"context"
	"net/http"

	"roo.bo/rlib/base"
)

const (
	__context_k = "__context_k"
)

func Context(r *http.Request) *base.Struct {
	c_i := r.Context().Value(__context_k)
	if c_i == nil {
		return nil
	}
	return c_i.(*base.Struct)
}

func Middleware_installContext(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		container := base.NewStruct(map[string]interface{}{})
		ctx := context.WithValue(r.Context(), __context_k, container)
		r = r.WithContext(ctx)
		handler.ServeHTTP(w, r)
	})

}
