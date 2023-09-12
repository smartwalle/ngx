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

var tests = []struct {
	method     string
	path       string
	rawQuery   string
	header     http.Header
	query      url.Values
	form       url.Values
	statusCode int
	response   []byte
}{
	{
		method:     http.MethodGet,
		path:       "/get200",
		statusCode: http.StatusOK,
		response:   []byte("Good"),
	},
	{
		method:     http.MethodGet,
		path:       "/get200_1",
		query:      url.Values{"q1": []string{"qv1"}},
		form:       url.Values{"f1": []string{"fv1"}},
		statusCode: http.StatusBadRequest,
		response:   []byte("Good"),
	},
	{
		method:     http.MethodGet,
		path:       "/get200_2",
		rawQuery:   "rq1=rqv1&rq2=rqv2",
		query:      url.Values{"q1": []string{"qv1"}},
		form:       url.Values{"f1": []string{"fv1"}},
		statusCode: http.StatusBadRequest,
		response:   []byte("Good"),
	},
	{
		method:     http.MethodGet,
		path:       "/get400",
		statusCode: http.StatusBadRequest,
		response:   []byte("Good"),
	},
	{
		method:     http.MethodPost,
		path:       "/post200_1",
		query:      url.Values{"q1": []string{"qv1"}},
		form:       url.Values{"f1": []string{"fv1"}},
		statusCode: http.StatusBadRequest,
		response:   []byte("Good"),
	},
	{
		method:     http.MethodPost,
		path:       "/post200_2",
		rawQuery:   "rq1=rqv1&rq2=rqv2",
		query:      url.Values{"q1": []string{"qv1"}},
		form:       url.Values{"f1": []string{"fv1"}},
		statusCode: http.StatusBadRequest,
		response:   []byte("Good"),
	},
}

func NewServer(t *testing.T) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()

		for _, test := range tests {
			if test.method == r.Method && test.path == r.URL.Path {
				var f1 = url.Values{}
				copyValues(f1, r.Form)

				var f2 = url.Values{}
				var rQuery, _ = url.ParseQuery(test.rawQuery)
				copyValues(f2, rQuery)
				copyValues(f2, test.query)
				copyValues(f2, test.form)

				if f1.Encode() != f2.Encode() {
					t.Fatalf("请求：%s-%s 参数不匹配, 期望: %s, 实际: %s \n", test.method, test.path, f2.Encode(), f1.Encode())
				}

				w.WriteHeader(test.statusCode)
				w.Write(test.response)
			}
		}
	}))
}

func copyValues(dst, src url.Values) {
	for k, vs := range src {
		dst[k] = append(dst[k], vs...)
	}
}

func TestNewRequest(t *testing.T) {
	var server = NewServer(t)

	for _, test := range tests {
		var req = ngx.NewRequest(test.method, server.URL+test.path+"?"+test.rawQuery, ngx.WithQuery(test.query))
		req.SetHeader(test.header)
		req.SetForm(ngx.CloneValues(test.form))

		var rsp, err = req.Do(context.Background())
		if err != nil {
			t.Fatalf("访问：%s-%s 发生错误: %v", test.method, test.path, err)
			continue
		}

		if rsp.StatusCode != test.statusCode {
			t.Fatalf("访问：%s-%s 期望: %d，实际：%d \n", test.method, test.path, test.statusCode, rsp.StatusCode)
		}

		var body, _ = io.ReadAll(rsp.Body)
		rsp.Body.Close()

		if bytes.Compare(body, test.response) != 0 {
			t.Fatalf("访问：%s-%s 期望: %s，实际：%s \n", test.method, test.path, string(test.response), string(body))
		}
	}
}
