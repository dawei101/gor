package rsession

import (
	"github.com/dawei101/gor/rcontext"

	"github.com/go-session/session"
)

func Session(r *http.Request) *session.Store {
	return rcontext.Get(r, "__session").(*session.Store)
}

func Destory(r *http.Request) {
}

func Middleware_installRSession(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		store, err := session.Start(nil, w, r)
		if err != nil {
			// 400 err
			return
		}
		rcontext.Set(r, "__session", &store)
		handler.ServeHTTP(w, r)
	})
}
