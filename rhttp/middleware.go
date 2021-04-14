package rhttp

import (
	"compress/gzip"
	"io"
	"net/http"
	"runtime/debug"
	"strings"

	"github.com/dawei101/gor/rlog"
	"github.com/dawei101/gor/rrouter"
)

func init() {
	rrouter.RegGlobalMiddleware(Middleware_panicLog)
}

type gzipResponseWriter struct {
	io.Writer
	http.ResponseWriter
}

func (w gzipResponseWriter) Write(b []byte) (int, error) {
	return w.ResponseWriter.Write(b)
}

func Middleware_gzip(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			handler.ServeHTTP(w, r)
			return
		}
		w.Header().Set("Content-Encoding", "gzip")
		gz := gzip.NewWriter(w)
		defer gz.Close()
		gzw := gzipResponseWriter{Writer: gz, ResponseWriter: w}
		handler.ServeHTTP(gzw, r)
	})
}

func Middleware_panicLog(handle http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				rlog.Error(r.Context(), err, strings.ReplaceAll(string(debug.Stack()), "\n", "\t"))
				FlushErr(w, r, NewRespErr(500, "server went wrong", ""))
			}
		}()
		handle(w, r)
	}
}
