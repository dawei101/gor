package rhttp

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/dawei101/gor/rcontext"
	"github.com/dawei101/gor/rlog"
)

func ListenAndServe(serveAt string, r *rrouter.Router) error {
	router.Use(
		rcontext.Middleware_installContext,
		Middleware_httpRequestLog,
		Middleware_panicLog)
	srv := &http.Server{
		Handler:      r,
		Addr:         serveAt,
		ReadTimeout:  10 * time.Minute,
		WriteTimeout: 10 * time.Minute,
	}
	fmt.Printf("Listening and serving HTTP on %s\n", srv.Addr)
	return srv.ListenAndServe()
}
