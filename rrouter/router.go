package rrouter

import (
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"
	"github.com/rs/xid"
)

type Middleware func(httprouter.Handle) httprouter.Handle

type Router struct {
	*httprouter.Router
	middlewares []Middleware
}

func New() *Router {
	return &Router{
		httprouter.New(),
		[]Middleware{},
	}
}

func (r *Router) Use(mw Middleware) {
	r.middlewares = append(r.middlewares, mw)
}

func (r *Router) Wrap(fn httprouter.Handle) httprouter.Handle {
	l := len(r.middlewares)
	if l == 0 {
		return fn
	}

	// There is at least one item in the list. Starting
	// with the last item, create the handler to be
	// returned:
	var result httprouter.Handle
	result = r.middlewares[l-1](fn)

	// Reverse through the stack for the remaining elements,
	// and wrap the result with each layer:
	for i := 0; i < (l - 1); i++ {
		result = r.middlewares[l-(2+i)](result)
	}

	return result
}

func Vars(r *http.Request) map[string]string {
	vars := map[string]string{}
	ps := httprouter.ParamsFromContext(r.Context())
	for _, p := range ps {
		vars[p.Key] = vars[p.Value]
	}
	return vars
}

func Var(r *http.Request, key string) string {
	ps := httprouter.ParamsFromContext(r.Context())
	return ps.ByName(key)
}

func (r *Router) Subroute(prefix string, sub *Router) {
	if strings.HasSuffix(prefix, "/") {
		panic("subroute's prefix could not end with `/`")
	}
	subpath_name := "subroute" + xid.New().String()
	path := prefix + "/:" + subpath_name
	handler := func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		path := "/" + ps.ByName(subpath_name)
		r.URL.Path = path
		sub.ServeHTTP(w, r)
	}
	for _, method := range []string{"GET", "POST", "PUT", "DELETE"} {
		r.Handle(method, path, handler)
	}
}
