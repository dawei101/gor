package rlog

import (
	"context"
	"fmt"
	"net/http"
	"sync/atomic"

	"github.com/rs/xid"
)

var reqIdGenerator = func(r *http.Request) string {
	return xid.New().String()
}

func SetReqIdGenerator(fn func(*http.Request) string) {
	reqIdGenerator = fn
}

func Middleware_installRLog(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handler.ServeHTTP(w, prepareRequest(r))
	})
}

const (
	rlogCounterKey = "__rlog_c"
	rlogReqidKey   = "_rlog_rid"
)

type counter struct {
	count int32
}

func (c *counter) Rise() int32 {
	return atomic.AddInt32(&c.count, 1)
}

func prepareRequest(r *http.Request) *http.Request {
	reqid := reqIdGenerator(r)
	ctx := context.WithValue(r.Context(), rlogCounterKey, &counter{})
	ctx = context.WithValue(ctx, rlogReqidKey, reqid)
	return r.WithContext(ctx)
}

func getReqId(ctx context.Context) string {
	s := ctx.Value(rlogReqidKey)
	if s == nil {
		return "nil"
	}
	return fmt.Sprintf("%s", s)
}

func logSerialNum(ctx context.Context) int32 {
	c := ctx.Value(rlogCounterKey)
	if c == nil {
		return 0
	}
	ct := c.(*counter)
	return ct.Rise()
}
