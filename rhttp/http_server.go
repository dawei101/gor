package rhttp

import (
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"runtime/debug"
	"strings"
	"time"

	"github.com/gorilla/mux"

	"roo.bo/rlib/rcontext"
	"roo.bo/rlib/rlog"
)

func HttpServe(router *mux.Router, serveAt string) error {
	router.Use(
		rcontext.Middleware_installContext,
		middleware_httpRequestLog,
		middleware_panicLog)
	srv := &http.Server{
		Handler:      router,
		Addr:         serveAt,
		ReadTimeout:  10 * time.Minute,
		WriteTimeout: 10 * time.Minute,
	}
	fmt.Printf("Listening and serving HTTP on %s\n", srv.Addr)
	return srv.ListenAndServe()
}

type gzipResponseWriter struct {
	io.Writer
	http.ResponseWriter
}

func (w gzipResponseWriter) Write(b []byte) (int, error) {
	return w.ResponseWriter.Write(b)
}

func Middleware_Gzip(handler http.Handler) http.Handler {
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

func middleware_panicLog(handle http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				rlog.Error(r.Context(), err, strings.ReplaceAll(string(debug.Stack()), "\n", "\t"))
				NewErrResp(-2, "server went wrong", "").Flush(w)
			}
		}()
		handle.ServeHTTP(w, r)

	})
}
