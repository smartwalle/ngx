package ngx_test

import (
	"bytes"
	"context"
	"github.com/smartwalle/ngx"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

type TestValue struct {
	method      string
	path        string
	contentType ngx.ContentType
	// 请求信息
	rawQuery string
	header   http.Header
	query    url.Values
	form     url.Values
	body     []byte
	// 响应信息
	rCode int
	rBody []byte

	builder func(t *testing.T, test TestValue, server *httptest.Server) *ngx.Request
	handler func(t *testing.T, test TestValue, w http.ResponseWriter, r *http.Request)
}

var tests = []TestValue{
	{
		method:      http.MethodGet,
		path:        "/get200",
		contentType: ngx.ContentTypeURLEncode,
		rCode:       http.StatusOK,
		rBody:       []byte("Good"),
	},
	{
		method:      http.MethodGet,
		path:        "/get200_1",
		contentType: ngx.ContentTypeURLEncode,
		query:       url.Values{"q1": []string{"qv1"}},
		form:        url.Values{"f1": []string{"fv1"}},
		rCode:       http.StatusOK,
		rBody:       []byte("Good"),
	},
	{
		method:      http.MethodGet,
		path:        "/get200_2",
		contentType: ngx.ContentTypeURLEncode,
		rawQuery:    "rq1=rqv1&rq2=rqv2",
		query:       url.Values{"q1": []string{"qv1"}},
		form:        url.Values{"f1": []string{"fv1"}},
		rCode:       http.StatusOK,
		rBody:       []byte("Good"),
	},
	{
		method:      http.MethodGet,
		path:        "/get200_3",
		contentType: ngx.ContentTypeURLEncode,
		rawQuery:    "rq1=rqv1&rq2=rqv2",
		query:       url.Values{"q1": []string{"qv1"}},
		form:        url.Values{"f1": []string{"fv1"}},
		body:        []byte("Get Body"),
		rCode:       http.StatusOK,
		rBody:       []byte("Good"),
	},
	{
		method:      http.MethodGet, // TrimURLQuery
		path:        "/get200_4",
		contentType: ngx.ContentTypeURLEncode,
		rawQuery:    "rq1=rqv1&rq2=rqv2",
		query:       url.Values{"q1": []string{"qv1"}},
		form:        url.Values{"f1": []string{"fv1"}},
		body:        []byte("Get Body"),
		rCode:       http.StatusOK,
		rBody:       []byte("Good"),
		builder: func(t *testing.T, test TestValue, server *httptest.Server) *ngx.Request {
			var req = ngx.NewRequest(test.method, server.URL+test.path+"?"+test.rawQuery)
			req.Header = test.header
			req.ContentType = test.contentType

			req.Query = ngx.CloneValues(test.query)
			req.Form = ngx.CloneValues(test.form)
			if len(test.body) > 0 {
				req.Body = bytes.NewReader(test.body)
			}
			return req
		},
		handler: func(t *testing.T, test TestValue, w http.ResponseWriter, r *http.Request) {
			r.ParseForm()

			var f1 = url.Values{}
			copyValues(f1, r.Form)

			var f2 = url.Values{}
			copyValues(f2, test.query)
			copyValues(f2, test.form)

			if f1.Encode() != f2.Encode() {
				t.Fatalf("请求：%s-%s 参数不匹配, 期望: %s, 实际: %s \n", test.method, test.path, f2.Encode(), f1.Encode())
			}

			var body, _ = io.ReadAll(r.Body)
			if bytes.Compare(body, test.body) != 0 {
				t.Fatalf("请求：%s-%s Body 不匹配, 期望: %s，实际：%s \n", test.method, test.path, string(test.body), string(body))
			}

			w.WriteHeader(test.rCode)
			w.Write(test.rBody)
		},
	},
	{
		method:      http.MethodGet,
		path:        "/get400",
		contentType: ngx.ContentTypeURLEncode,
		rCode:       http.StatusBadRequest,
		rBody:       []byte("Good"),
	},
	{
		method:      http.MethodGet,
		path:        "/content_type",
		contentType: ngx.ContentTypeText,
		rCode:       http.StatusBadRequest,
		rBody:       []byte("Good"),
		builder: func(t *testing.T, test TestValue, server *httptest.Server) *ngx.Request {
			var req = ngx.NewRequest(test.method, server.URL+test.path+"?"+test.rawQuery)
			// Header 中有设置 Content-Type 时，单独调用 SetContentType 方法设置的 Content-Type 将被忽略
			req.Header.Set("Content-Type", string(ngx.ContentTypeURLEncode))
			req.ContentType = test.contentType

			req.Query = ngx.CloneValues(test.query)
			req.Form = ngx.CloneValues(test.form)
			if len(test.body) > 0 {
				req.Body = bytes.NewReader(test.body)
			}
			return req
		},
		handler: func(t *testing.T, test TestValue, w http.ResponseWriter, r *http.Request) {
			r.ParseForm()

			var f1 = url.Values{}
			copyValues(f1, r.Form)

			var f2 = url.Values{}
			copyValues(f2, test.query)
			copyValues(f2, test.form)

			if r.Header.Get("Content-Type") != string(ngx.ContentTypeText) {
				t.Fatalf("请求：%s-%s ContentType 不匹配, 期望: %s, 实际: %s \n", test.method, test.path, ngx.ContentTypeText, r.Header.Get("Content-Type"))
			}

			if f1.Encode() != f2.Encode() {
				t.Fatalf("请求：%s-%s 参数不匹配, 期望: %s, 实际: %s \n", test.method, test.path, f2.Encode(), f1.Encode())
			}

			var body, _ = io.ReadAll(r.Body)
			if bytes.Compare(body, test.body) != 0 {
				t.Fatalf("请求：%s-%s Body 不匹配, 期望: %s，实际：%s \n", test.method, test.path, string(test.body), string(body))
			}

			w.WriteHeader(test.rCode)
			w.Write(test.rBody)
		},
	},
	{
		method:      http.MethodPost,
		path:        "/post200_1",
		contentType: ngx.ContentTypeURLEncode,
		query:       url.Values{"q1": []string{"qv1"}},
		form:        url.Values{"f1": []string{"fv1"}},
		rCode:       http.StatusOK,
		rBody:       []byte("Good"),
	},
	{
		method:      http.MethodPost,
		path:        "/post200_2",
		contentType: ngx.ContentTypeURLEncode,
		rawQuery:    "rq1=rqv1&rq2=rqv2",
		query:       url.Values{"q1": []string{"qv1"}},
		form:        url.Values{"f1": []string{"fv1"}},
		rCode:       http.StatusOK,
		rBody:       []byte("Good"),
	},
	{
		method:      http.MethodPost,
		path:        "/post200_3",
		contentType: ngx.ContentTypeText,
		rawQuery:    "rq1=rqv1&rq2=rqv2",
		query:       url.Values{"q1": []string{"qv1"}},
		form:        url.Values{"f1": []string{"fv1"}},
		body:        []byte("Get Body"),
		rCode:       http.StatusOK,
		rBody:       []byte("Good"),
		handler: func(t *testing.T, test TestValue, w http.ResponseWriter, r *http.Request) {
			r.ParseForm()

			var f1 = url.Values{}
			copyValues(f1, r.Form)

			var f2 = url.Values{}
			var rQuery, _ = url.ParseQuery(test.rawQuery)
			copyValues(f2, rQuery)
			copyValues(f2, test.query)

			if f1.Encode() != f2.Encode() {
				t.Fatalf("请求：%s-%s 参数不匹配, 期望: %s, 实际: %s \n", test.method, test.path, f2.Encode(), f1.Encode())
			}

			var body, _ = io.ReadAll(r.Body)
			if bytes.Compare(body, test.body) != 0 {
				t.Fatalf("请求：%s-%s Body 不匹配, 期望: %s，实际：%s \n", test.method, test.path, string(test.body), string(body))
			}

			w.WriteHeader(test.rCode)
			w.Write(test.rBody)
		},
	},
}

func defaultRequest(t *testing.T, test TestValue, server *httptest.Server) *ngx.Request {
	var req = ngx.NewRequest(test.method, server.URL+test.path+"?"+test.rawQuery)
	req.Header = test.header
	req.ContentType = test.contentType

	for key, values := range test.query {
		for _, value := range values {
			req.Query.Add(key, value)
		}
	}
	for key, values := range test.form {
		for _, value := range values {
			req.Form.Add(key, value)
		}
	}
	if len(test.body) > 0 {
		req.Body = bytes.NewReader(test.body)
	}
	return req
}

func defaultHandler(t *testing.T, test TestValue, w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	var f1 = url.Values{}
	copyValues(f1, r.Form)

	var f2 = url.Values{}
	var rQuery, _ = url.ParseQuery(test.rawQuery)
	copyValues(f2, rQuery)
	copyValues(f2, test.query)
	copyValues(f2, test.form)

	if r.Header.Get("Content-Type") != string(test.contentType) {
		t.Fatalf("请求：%s-%s ContentType 不匹配, 期望: %s, 实际: %s \n", test.method, test.path, test.contentType, r.Header.Get("Content-Type"))
	}

	if f1.Encode() != f2.Encode() {
		t.Fatalf("请求：%s-%s 参数不匹配, 期望: %s, 实际: %s \n", test.method, test.path, f2.Encode(), f1.Encode())
	}

	var body, _ = io.ReadAll(r.Body)
	if !bytes.Equal(body, test.body) {
		t.Fatalf("请求：%s-%s Body 不匹配, 期望: %s，实际：%s \n", test.method, test.path, string(test.body), string(body))
	}

	w.WriteHeader(test.rCode)
	w.Write(test.rBody)
}

func NewServer(t *testing.T) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for _, test := range tests {
			if test.method == r.Method && test.path == r.URL.Path {
				var handler = test.handler
				if handler == nil {
					handler = defaultHandler
				}
				handler(t, test, w, r)
			}
		}
	}))
}

func copyValues(dst, src url.Values) {
	for k, vs := range src {
		dst[k] = append(dst[k], vs...)
	}
}

func TestRequest_Do(t *testing.T) {
	var server = NewServer(t)

	for _, test := range tests {
		var requestBuilder = test.builder
		if requestBuilder == nil {
			requestBuilder = defaultRequest
		}
		var req = requestBuilder(t, test, server)

		var rsp, err = req.Do(context.Background())
		if err != nil {
			t.Fatalf("访问：%s-%s 发生错误: %v", test.method, test.path, err)
			continue
		}

		if rsp.StatusCode != test.rCode {
			t.Fatalf("访问：%s-%s 期望: %d，实际：%d \n", test.method, test.path, test.rCode, rsp.StatusCode)
		}

		var body, _ = io.ReadAll(rsp.Body)
		rsp.Body.Close()

		if bytes.Compare(body, test.rBody) != 0 {
			t.Fatalf("访问：%s-%s 期望: %s，实际：%s \n", test.method, test.path, string(test.rBody), string(body))
		}
	}
}
