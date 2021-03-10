package rhttp

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/rs/xid"

	"github.com/dawei101/gor/base"
	"github.com/dawei101/gor/rlog"
)

var (
	reqlog *rlog.Log
	apiLog *rlog.Log
)

const (
	_log_extra   = "_log_extra"
	RequestIdKey = "*req*"
)

func getReqLog() *rlog.Log {
	if reqlog != nil {
		return reqlog
	}
	return rlog.DefLog()
}

func GenerateReqId(r *http.Request) string {
	reqids_, ok := r.URL.Query()[RequestIdKey]
	if !ok || len(reqids_) == 0 {
		return xid.New().String()
	} else {
		return reqids_[0]
	}
}

type reqLogExtra struct {
	*base.Struct
}

func requestLogExtra(r *http.Request) *reqLogExtra {
	extra_i := r.Context().Value(_log_extra)
	if extra_i == nil {
		return &reqLogExtra{base.NewStruct(map[string]interface{}{})}
	}
	return extra_i.(*reqLogExtra)
}

// 在request的log中添加extra信息
//
func AddRequestLogExtra(r *http.Request, name string, val interface{}) {
	requestLogExtra(r).Set(name, val)
}

func prepareContext(r *http.Request) *http.Request {
	ctx := context.WithValue(r.Context(), _log_extra, requestLogExtra(r))
	return r.WithContext(ctx)
}

type multiResponseWriter struct {
	rw http.ResponseWriter
	mw io.Writer
}

func newMultiResponseWriter(w http.ResponseWriter, buf io.Writer) *multiResponseWriter {
	return &multiResponseWriter{
		rw: w,
		mw: io.MultiWriter(w, buf),
	}
}

func (r *multiResponseWriter) Header() http.Header {
	return r.rw.Header()
}

func (r *multiResponseWriter) Write(i []byte) (int, error) {
	return r.mw.Write(i)
}

func (r *multiResponseWriter) WriteHeader(statusCode int) {
	r.rw.WriteHeader(statusCode)
}

func Middleware_httpRequestLog(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		r = prepareContext(r)

		// make multipe writer
		resp := &bytes.Buffer{}
		m := newMultiResponseWriter(w, resp)
		bodystruct, err := JsonBody(r)
		body := []byte{}
		if err != nil {
			body = []byte(err.Error())
		} else {
			body = bodystruct.JsonMarshal()
		}

		st := time.Now().Local()

		handler.ServeHTTP(m, r)

		et := time.Now().Local()
		dt := et.Sub(st)
		url := r.URL
		method := r.Method

		loginfos := []interface{}{
			fmt.Sprintf("method(%s)", method),
			fmt.Sprintf("url(%s)", url.String()),
			fmt.Sprintf("startms(%d)", st.UnixNano()/1e6),
			fmt.Sprintf("usedms(%d)", dt.Milliseconds()),
			fmt.Sprintf("request(%s)", body),
			fmt.Sprintf("response(%s)", resp),
			fmt.Sprintf("extra(%s)", requestLogExtra(r).JsonMarshal()),
		}

		getReqLog().Info(r.Context(), loginfos)
	})
}
