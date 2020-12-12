package rhttp

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"

	"github.com/dawei101/gor/rcontext"
	"github.com/dawei101/gor/rlog"
)

func Serve(router *mux.Router, serveAt string) error {
	router.Use(
		rcontext.Middleware_installContext,
		Middleware_httpRequestLog,
		Middleware_panicLog)
	srv := &http.Server{
		Handler:      router,
		Addr:         serveAt,
		ReadTimeout:  10 * time.Minute,
		WriteTimeout: 10 * time.Minute,
	}
	fmt.Printf("Listening and serving HTTP on %s\n", srv.Addr)
	return srv.ListenAndServe()
}

