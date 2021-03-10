package rhttp

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/dawei101/gor/rlog"
)

const (
	RequestIdKey = "*req*"
)

type _ParamsType int32

const (
	_ParamsTypeNONE         _ParamsType = 0
	_ParamsTypeURLValueLike _ParamsType = 1
	_ParamsTypeStrings      _ParamsType = 2
	_ParamsTypeOther        _ParamsType = 3
)

// DefaultTransport 默认http transport 配置
var defaultHttpTransport = &http.Transport{
	DialContext: (&net.Dialer{
		Timeout:   time.Second,
		KeepAlive: 90 * time.Minute,
		DualStack: false,
	}).DialContext,
	MaxIdleConns:          200,
	MaxIdleConnsPerHost:   10,
	IdleConnTimeout:       90 * time.Second,
	ExpectContinueTimeout: 10 * time.Second,
	ResponseHeaderTimeout: 3 * time.Second,
}

// defaultClient 默认http.Client 链接
var defaultHttpClient = &http.Client{Transport: defaultHttpTransport}

// HTTPClientOption 配置选项
type HTTPClientOption struct {
	Retry         int
	RetryInterval time.Duration
	CloseLog      bool
	Timeout       time.Duration
	Mock          string
}

// defaultHttpClientOption 默认配置
var defaultHttpClientOption = &HTTPClientOption{
	Retry:         0,
	RetryInterval: time.Millisecond * 200,
	CloseLog:      false,
	Timeout:       time.Millisecond * 500,
}

// HTTPClient struct
type HTTPClient struct {
	topic  string
	option *HTTPClientOption
	client *http.Client
	ctx    context.Context
}

func NewHTTPClient(ctx context.Context, topic string, option *HTTPClientOption, httpclient *http.Client) *HTTPClient {
	return &HTTPClient{
		topic:  topic,
		client: httpclient,
		option: option,
		ctx:    ctx,
	}
}

func (c HTTPClient) requestLog(url string, requeststr string, responsestr string, usedtm int, errormsg string, subTopic string, code int, method string) {
	if c.option.CloseLog {
		return
	}

	logstr := fmt.Sprintf("[REQUEST-TAG][%s] url(%s), method(%s), request(%s), response(%s), usedtm(%d), code(%d), error(%s)", c.topic, url, method, requeststr, responsestr, usedtm, code, errormsg)

	if code == 200 {
		getApiLog().Info(c.ctx, logstr)
	} else {
		getApiLog().Warning(c.ctx, logstr)
	}
}

//Get method
// body 内容可以为一下数据结构 url.Values,map[string]string, map[string][string]
//
func (c HTTPClient) Get(urlstr string, body interface{}, header map[string]string) (*http.Response, error) {
	urlstr = addParams(urlstr, toURLValues(body))
	return c.doCurlExec("GET", urlstr, "", header)
}

//Post method
// body 如果传入内容 url.Values,map[string]string, map[string][string] 那么将按照Content-Type=application/x-www-form-urlencoded 进行请求
// 如果body实体中含有文件，那么Content-Type=multipart/form-data 传输
// 如果body 实体为 string,[]byte,io.Reader 那么直接进行post请求,额外的需要自定义content-type
func (c HTTPClient) Post(urlstr string, body interface{}, header map[string]string) (*http.Response, error) {
	t := checkParamsType(body)
	if t == _ParamsTypeStrings {
		return c.doCurlExec("POST", urlstr, toString(body), header)
	}

	paramsvalues := toURLValues(body)
	if checkParamsFile(paramsvalues) {
		c.PostMultipart(urlstr, paramsvalues, header)
	}

	if header == nil {
		header = make(map[string]string)
	}
	header["Content-Type"] = "application/x-www-form-urlencoded"
	return c.doCurlExec("POST", urlstr, paramsvalues.Encode(), header)
}

// PostMultipart method
func (c HTTPClient) PostMultipart(urlstr string, body interface{}, header map[string]string) (*http.Response, error) {
	postbody := &bytes.Buffer{}
	writer := multipart.NewWriter(postbody)
	paramsvalues := toURLValues(body)
	for k, v := range paramsvalues {
		for _, vv := range v {
			if k[0] == '@' {
				err := addFormFile(writer, k[1:], vv)
				if err != nil {
					return nil, err
				}
			} else {
				writer.WriteField(k, vv)
			}
		}
	}
	if header == nil {
		header = make(map[string]string)
	}
	header["Content-Type"] = writer.FormDataContentType()
	err := writer.Close()
	if err != nil {
		return nil, err
	}
	copypostbody, err := ioutil.ReadAll(postbody)
	if err != nil {
		return nil, err
	}
	return c.doCurlExec("POST", urlstr, string(copypostbody), header)
}

func (c HTTPClient) sendJSON(method, urlstr string, data interface{}, header map[string]string) (*http.Response, error) {
	if header == nil {
		header = make(map[string]string)
	}
	var body string
	header["Content-Type"] = "application/json"
	switch t := data.(type) {
	case []byte:
		body = string(t)
	case string:
		body = t
	default:
		bodybyte, err := json.Marshal(data)
		if err != nil {
			return nil, err
		}
		body = string(bodybyte)
	}

	return c.doCurlExec(method, urlstr, body, header)
}

//PostJSON method
func (c HTTPClient) PostJSON(urlstr string, data interface{}, header map[string]string) (*http.Response, error) {
	return c.sendJSON("POST", urlstr, data, header)
}

func (c HTTPClient) doCurlExec(method, url string, body string, header map[string]string) (resp *http.Response, err error) {
	resp, err = c.doCurl(method, url, body, header)
	retry := c.option.Retry
	retryInterval := c.option.RetryInterval
	idleDelay := time.NewTimer(retryInterval)
	defer idleDelay.Stop()
	for resp == nil && err != nil && retry > 0 {
		idleDelay.Reset(retryInterval)
		select {
		case <-c.ctx.Done():
			return
		// case <-time.After(retryInterval):
		case <-idleDelay.C:
			// fmt.Println("retry")
			resp, err = c.doCurl(method, url, body, header)
			retry--
		}
	}
	return
}

func (c HTTPClient) doCurl(method, urlstr string, body string, header map[string]string) (*http.Response, error) {

	params := map[string]string{
		RequestIdKey: rlog.CtxId(c.ctx),
	}
	urlstr = addParams(urlstr, toURLValues(params))

	request, err := http.NewRequest(method, urlstr, strings.NewReader(body))
	if err != nil {
		return nil, err
	}
	nctx := c.ctx
	if c.option.Timeout > 0 {
		var cancel context.CancelFunc
		nctx, cancel = context.WithTimeout(nctx, c.option.Timeout)
		defer cancel()
	}

	request = request.WithContext(nctx)

	reqStartTime := int(time.Now().UnixNano() / 1e6)
	for key, value := range header {
		request.Header.Add(key, value)
	}
	errno := 200
	errmsg := ""
	resp, err := c.client.Do(request)
	respstr := ""

	defer func() {
		usetm := int(time.Now().UnixNano()/1e6) - reqStartTime
		c.requestLog(request.URL.String(), body, respstr, usetm, errmsg, "", errno, request.Method)
	}()

	if err != nil {
		errno = 0
		errmsg = err.Error()
		return resp, err
	} else {
		errno = resp.StatusCode
		respbody := []byte{}
		respbody, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			resp.Body.Close()
			return nil, err
		}
		resp.Body.Close()
		if errno != 200 {
			errmsg = respstr
		} else {
			respstr = string(respbody)
		}
	}

	if err == nil {
		resp.Body = ioutil.NopCloser(strings.NewReader(respstr))
	}
	return resp, err
}

//

func toURLValues(v interface{}) url.Values {
	switch t := v.(type) {
	case url.Values:
		return t
	case map[string][]string:
		return url.Values(t)
	case map[string]string:
		rst := make(url.Values)
		for k, v := range t {
			rst.Add(k, v)
		}
		return rst
	case nil:
		return make(url.Values)
	default:
		panic("Invalid value")
	}
}

func addParams(url string, params url.Values) string {
	if len(params) == 0 {
		return url
	}
	if !strings.Contains(url, "?") {
		url += "?"
	}
	if strings.HasSuffix(url, "?") || strings.HasSuffix(url, "&") {
		url += params.Encode()
	} else {
		url += "&" + params.Encode()
	}

	return url
}

func toString(params interface{}) string {
	switch params.(type) {
	case string:
		return params.(string)
	case []byte:
		return string(params.([]byte))
	case *bytes.Reader:
		rst, err := ioutil.ReadAll(params.(*bytes.Reader))
		if err != nil {
			return ""
		}
		return string(rst)
	default:
		panic("Invalid value")
	}
}

func addFormFile(writer *multipart.Writer, name, path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()
	part, err := writer.CreateFormFile(name, filepath.Base(path))
	if err != nil {
		return err
	}
	_, err = io.Copy(part, file)
	return err
}

func checkParamsType(v interface{}) _ParamsType {
	switch v.(type) {
	case url.Values, map[string][]string, map[string]string:
		return _ParamsTypeURLValueLike
	case []byte, string, *bytes.Reader:
		return _ParamsTypeStrings
	case nil:
		return _ParamsTypeNONE
	default:
		return _ParamsTypeOther

	}
}

func checkParamsFile(params url.Values) bool {
	for k := range params {
		if k[0] == '@' {
			return true
		}
	}

	return false
}
func SetApiLog(log *rlog.Log) {
	apiLog = log
}

func getApiLog() *rlog.Log {
	if apiLog != nil {
		return apiLog
	}
	return rlog.DefaultLog()
}
