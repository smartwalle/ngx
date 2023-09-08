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
	body    io.Reader
	client  *http.Client
	header  http.Header
	query   url.Values
	form    url.Values
	files   FormFiles
	Method  string
	target  string
	cookies []*http.Cookie
}

func NewRequest(method, target string, opts ...Option) *Request {
	var nURL, _ = url.Parse(target)
	var req = &Request{}
	req.Method = strings.ToUpper(method)
	req.target = target

	if nURL != nil {
		req.query = nURL.Query()
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

func (this *Request) SetContentType(contentType ContentType) {
	this.Header().Set(kContentType, string(contentType))
}

func (this *Request) SetBody(body io.Reader) {
	this.body = body
}

func (this *Request) SetForm(form url.Values) {
	this.form = form
}

func (this *Request) Form() url.Values {
	if this.form == nil {
		this.form = url.Values{}
	}
	return this.form
}

func (this *Request) SetQuery(query url.Values) {
	this.query = query
}

func (this *Request) Query() url.Values {
	if this.query == nil {
		this.query = url.Values{}
	}
	return this.query
}

func (this *Request) SetHeader(header http.Header) {
	this.header = header
}

func (this *Request) Header() http.Header {
	if this.header == nil {
		this.header = http.Header{}
	}
	return this.header
}

func (this *Request) FileForm() FormFiles {
	if this.files == nil {
		this.files = FormFiles{}
	}
	return this.files
}

func (this *Request) SetFileForm(files FormFiles) {
	this.files = files
}

func (this *Request) AddCookie(cookie *http.Cookie) {
	this.cookies = append(this.cookies, cookie)
}

func (this *Request) SetCookies(cookies []*http.Cookie) {
	this.cookies = cookies
}

func (this *Request) Do(ctx context.Context) (*http.Response, error) {
	var req *http.Request
	var err error
	var body io.Reader
	var toQuery bool

	if this.Method == http.MethodGet ||
		this.Method == http.MethodTrace ||
		this.Method == http.MethodOptions ||
		this.Method == http.MethodHead ||
		this.Method == http.MethodDelete {
		toQuery = true
	}

	if this.body != nil {
		body = this.body
	} else if len(this.files) > 0 {
		var bodyBuffer = &bytes.Buffer{}
		var bodyWriter = multipart.NewWriter(bodyBuffer)

		for _, f := range this.files {
			if err = f.WriteTo(bodyWriter); err != nil {
				return nil, err
			}
		}
		for key, values := range this.form {
			for _, value := range values {
				bodyWriter.WriteField(key, value)
			}
		}

		if err = bodyWriter.Close(); err != nil {
			return nil, err
		}

		this.SetContentType(ContentType(bodyWriter.FormDataContentType()))
		body = bodyBuffer
	} else if len(this.form) > 0 && !toQuery {
		body = strings.NewReader(this.form.Encode())
	}

	req, err = http.NewRequestWithContext(ctx, this.Method, this.target, body)
	if err != nil {
		return nil, err
	}

	if toQuery {
		for key, values := range this.form {
			for _, value := range values {
				this.Query().Add(key, value)
			}
		}
	}

	req.URL.RawQuery = this.query.Encode()

	var header = this.Header()
	if _, ok := header[kContentType]; !ok {
		header.Set(kContentType, string(ContentTypeURLEncode))
	}
	req.Header = header

	for _, cookie := range this.cookies {
		req.AddCookie(cookie)
	}

	return this.client.Do(req)
}
