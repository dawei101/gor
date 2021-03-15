// +build !testing

package rlog

import (
	"github.com/dawei101/gor/rrouter"
)

func init() {
	if err := loadConfig(); err != nil {
		panic(err)
	}
	rrouter.RegGlobalMiddleware(Middleware_installRLog)
	go midnightRotate()
}
