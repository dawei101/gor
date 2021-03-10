package rhttp

import (
	"fmt"
	"net/http"
	"time"

	"github.com/dawei101/gor/rrouter"
)

func ListenAndServe(serveAt string, r *rrouter.Router) error {
	srv := &http.Server{
		Handler:      r,
		Addr:         serveAt,
		ReadTimeout:  10 * time.Minute,
		WriteTimeout: 10 * time.Minute,
	}
	fmt.Printf("Listening and serving HTTP on %s\n", srv.Addr)
	return srv.ListenAndServe()
}
