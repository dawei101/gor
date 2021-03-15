package rrouter

import (
	"net/http"
	"strings"

	"github.com/dawei101/gor/base"
	"github.com/julienschmidt/httprouter"
	"github.com/rs/xid"
)

var globalMiddlewares *base.OrderedMap

func RegGlobalMiddleware(mw Middleware) {
	if _, ok := globalMiddlewares.Get(mw); ok {
		return
	}
	globalMiddlewares.Set(mw, true)
}

type Middleware func(http.HandlerFunc) http.HandlerFunc

type Router struct {
	*httprouter.Router
	middlewares []Middleware
}

func New() *Router {
	mws := []Middleware{}
	for _, mw := range globalMiddlewares.Keys() {
		mws = append(mws, mw.(Middleware))
	}
	return &Router{
		httprouter.New(),
		mws,
	}
}

func (r *Router) Use(mw Middleware) {
	r.middlewares = append(r.middlewares, mw)
}

func (r *Router) Wrap(fn http.HandlerFunc) http.HandlerFunc {
	l := len(r.middlewares)
	if l == 0 {
		return fn
	}

	// There is at least one item in the list. Starting
	// with the last item, create the handler to be
	// returned:
	var result http.HandlerFunc
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
	handler := func(w http.ResponseWriter, req *http.Request) {
		path := "/" + Var(req, subpath_name)
		req.URL.Path = path
		sub.ServeHTTP(w, req)
	}
	for _, method := range []string{"GET", "POST", "PUT", "DELETE"} {
		r.Handle(method, path, handler)
	}
}

// GET is a shortcut for router.Handle(http.MethodGet, path, handle)
func (r *Router) GET(path string, handle http.HandlerFunc) {
	r.Handle(http.MethodGet, path, handle)
}

// HEAD is a shortcut for router.Handle(http.MethodHead, path, handle)
func (r *Router) HEAD(path string, handle http.HandlerFunc) {
	r.Handle(http.MethodHead, path, handle)
}

// OPTIONS is a shortcut for router.Handle(http.MethodOptions, path, handle)
func (r *Router) OPTIONS(path string, handle http.HandlerFunc) {
	r.Handle(http.MethodOptions, path, handle)
}

// POST is a shortcut for router.Handle(http.MethodPost, path, handle)
func (r *Router) POST(path string, handle http.HandlerFunc) {
	r.Handle(http.MethodPost, path, handle)
}

// PUT is a shortcut for router.Handle(http.MethodPut, path, handle)
func (r *Router) PUT(path string, handle http.HandlerFunc) {
	r.Handle(http.MethodPut, path, handle)
}

// PATCH is a shortcut for router.Handle(http.MethodPatch, path, handle)
func (r *Router) PATCH(path string, handle http.HandlerFunc) {
	r.Handle(http.MethodPatch, path, handle)
}

// DELETE is a shortcut for router.Handle(http.MethodDelete, path, handle)
func (r *Router) DELETE(path string, handle http.HandlerFunc) {
	r.Handle(http.MethodDelete, path, handle)
}

// Handle registers a new request handle with the given path and method.
//
// For GET, POST, PUT, PATCH and DELETE requests the respective shortcut
// functions can be used.
//
// This function is intended for bulk loading and to allow the usage of less
// frequently used, non-standardized or custom methods (e.g. for internal
// communication with a proxy).
func (r *Router) Handle(method, path string, handle http.HandlerFunc) {
	r.HandlerFunc(method, path, handle)
}
