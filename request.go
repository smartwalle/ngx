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
	body        io.Reader
	client      *http.Client
	header      http.Header
	uQuery      url.Values
	query       url.Values
	form        url.Values
	files       FormFiles
	Method      string
	target      string
	contentType ContentType
	cookies     []*http.Cookie
}

func NewRequest(method, target string, opts ...Option) *Request {
	var nURL, _ = url.Parse(target)
	var req = &Request{}
	req.Method = strings.ToUpper(method)
	req.target = target
	req.contentType = ContentTypeURLEncode

	if nURL != nil {
		req.uQuery = nURL.Query()
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

func (r *Request) TrimURLQuery() {
	r.uQuery = nil
}

func (r *Request) SetContentType(contentType ContentType) {
	r.contentType = contentType
}

func (r *Request) SetBody(body io.Reader) {
	r.body = body
}

func (r *Request) SetForm(form url.Values) {
	r.form = form
}

func (r *Request) Form() url.Values {
	if r.form == nil {
		r.form = url.Values{}
	}
	return r.form
}

func (r *Request) SetQuery(query url.Values) {
	r.query = query
}

func (r *Request) Query() url.Values {
	if r.query == nil {
		r.query = url.Values{}
	}
	return r.query
}

func (r *Request) SetHeader(header http.Header) {
	r.header = header
}

func (r *Request) Header() http.Header {
	if r.header == nil {
		r.header = http.Header{}
	}
	return r.header
}

func (r *Request) FileForm() FormFiles {
	if r.files == nil {
		r.files = FormFiles{}
	}
	return r.files
}

func (r *Request) SetFileForm(files FormFiles) {
	r.files = files
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
	var toQuery bool

	if r.Method == http.MethodGet ||
		r.Method == http.MethodTrace ||
		r.Method == http.MethodOptions ||
		r.Method == http.MethodHead ||
		r.Method == http.MethodDelete {
		toQuery = true
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
	} else if len(r.form) > 0 && !toQuery {
		body = strings.NewReader(r.form.Encode())
	}

	req, err = http.NewRequestWithContext(ctx, r.Method, r.target, body)
	if err != nil {
		return nil, err
	}

	var rawQuery = r.uQuery
	if rawQuery == nil {
		rawQuery = url.Values{}
	}

	for key, values := range r.query {
		for _, value := range values {
			rawQuery.Add(key, value)
		}
	}

	if toQuery {
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
