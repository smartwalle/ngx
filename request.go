package ngx

import (
	"bytes"
	"context"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"
)

type ContentType string

const kContentType = "Content-Type"

const (
	ContentTypeJSON      ContentType = "application/json"
	ContentTypeXML       ContentType = "application/xml"
	ContentTypeURLEncode ContentType = "application/x-www-form-urlencoded"
	ContentTypeHTML      ContentType = "text/html"
	ContentTypeText      ContentType = "text/plain"
	ContentTypeMultipart ContentType = "multipart/form-data"
)

const (
	Post    = http.MethodPost
	Get     = http.MethodGet
	Head    = http.MethodHead
	Put     = http.MethodPut
	Delete  = http.MethodDelete
	Patch   = http.MethodPatch
	Options = http.MethodOptions
)

type Request struct {
	client      *http.Client
	header      http.Header
	body        io.Reader
	method      string
	rawURL      string
	rawQuery    url.Values
	query       url.Values
	form        url.Values
	files       FormFiles
	contentType ContentType
	cookies     []*http.Cookie
}

func NewRequest(method, rawURL string, opts ...Option) *Request {
	var nURL, _ = url.Parse(rawURL)
	var req = &Request{}
	req.method = strings.ToUpper(method)
	req.rawURL = rawURL
	req.contentType = ContentTypeURLEncode

	if nURL != nil {
		req.rawQuery = nURL.Query()
	}

	for _, opt := range opts {
		if opt != nil {
			opt(req)
		}
	}

	if req.client == nil {
		req.client = http.DefaultClient
	}

	return req
}

// TrimURLQuery 清除 URL 中原有的查询参数信息
func (r *Request) TrimURLQuery() {
	r.rawQuery = nil
}

// SetContentType 设置 Content-Type
func (r *Request) SetContentType(contentType ContentType) {
	r.contentType = contentType
}

// SetBody 设置请求体
//
// 如果同时设置了 Body 和 Form，Body 的优先级则高于 Form，Form 中的信息将被舍弃。
func (r *Request) SetBody(body io.Reader) {
	r.body = body
}

// SetForm 设置请求参数（表单）
func (r *Request) SetForm(form url.Values) {
	r.form = form
}

// Form 获取请求参数（表单）
//
// 对于 POST 请求，该参数将通过 Body 传递；
// 对于 GET 一类的请求，该参数将拼接在 URL 的查询参数中。
func (r *Request) Form() url.Values {
	if r.form == nil {
		r.form = url.Values{}
	}
	return r.form
}

// SetQuery 设置查询参数信息
//
// 该参数将拼接在 URL 的查询参数中。
func (r *Request) SetQuery(query url.Values) {
	r.query = query
}

// Query 获取查询参数信息
func (r *Request) Query() url.Values {
	if r.query == nil {
		r.query = url.Values{}
	}
	return r.query
}

// SetHeader 设置请求头
func (r *Request) SetHeader(header http.Header) {
	r.header = header
}

// Header 获取请求头
func (r *Request) Header() http.Header {
	if r.header == nil {
		r.header = http.Header{}
	}
	return r.header
}

// SetFileForm 设置上传文件信息
func (r *Request) SetFileForm(files FormFiles) {
	r.files = files
}

// FileForm 获取上传文件信息
func (r *Request) FileForm() FormFiles {
	if r.files == nil {
		r.files = FormFiles{}
	}
	return r.files
}

func (r *Request) AddCookie(cookie *http.Cookie) {
	r.cookies = append(r.cookies, cookie)
}

func (r *Request) SetCookies(cookies []*http.Cookie) {
	r.cookies = cookies
}

func (r *Request) Do(ctx context.Context) (*http.Response, error) {
	var req *http.Request
	var err error
	var body io.Reader
	var mergeToRawQuery bool

	if r.method == http.MethodGet ||
		r.method == http.MethodTrace ||
		r.method == http.MethodOptions ||
		r.method == http.MethodHead ||
		r.method == http.MethodDelete {
		mergeToRawQuery = true
	}

	if r.body != nil {
		body = r.body
	} else if len(r.files) > 0 {
		var bodyBuffer = &bytes.Buffer{}
		var bodyWriter = multipart.NewWriter(bodyBuffer)

		for _, f := range r.files {
			if err = f.WriteTo(bodyWriter); err != nil {
				return nil, err
			}
		}
		for key, values := range r.form {
			for _, value := range values {
				bodyWriter.WriteField(key, value)
			}
		}

		if err = bodyWriter.Close(); err != nil {
			return nil, err
		}

		r.SetContentType(ContentType(bodyWriter.FormDataContentType()))
		body = bodyBuffer
	} else if len(r.form) > 0 && !mergeToRawQuery {
		body = strings.NewReader(r.form.Encode())
	}

	req, err = http.NewRequestWithContext(ctx, r.method, r.rawURL, body)
	if err != nil {
		return nil, err
	}

	var rawQuery = CloneValues(r.rawQuery)
	if rawQuery == nil {
		rawQuery = url.Values{}
	}

	for key, values := range r.query {
		for _, value := range values {
			rawQuery.Add(key, value)
		}
	}

	if mergeToRawQuery {
		for key, values := range r.form {
			for _, value := range values {
				rawQuery.Add(key, value)
			}
		}
	}
	req.URL.RawQuery = rawQuery.Encode()

	var header = r.Header()
	if _, ok := header[kContentType]; !ok {
		header.Set(kContentType, string(r.contentType))
	}
	req.Header = header

	for _, cookie := range r.cookies {
		req.AddCookie(cookie)
	}

	return r.client.Do(req)
}
