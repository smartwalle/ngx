package ngx

import (
	"context"
	"io"
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
	method      string
	url         *url.URL
	Client      *http.Client
	Header      http.Header
	Body        BodyEncoder // 如果同时设置了 Body 和 Form(FileForm)，则 Body 的优先级高于 Form(FileForm)，且 Form(FileForm) 中的信息将被舍弃。
	Query       url.Values  // 该参数将拼接在 URL 的查询参数中。
	Form        url.Values  // 对于 POST 请求，该参数将通过 Body 传递；对于 GET 一类的请求，该参数将和 Query 合并之后，拼接在 URL 的查询参数中。
	FileForm    FileForm    // 上传文件。
	ContentType ContentType // 如果设置了 ContentType，则会覆盖 Header 中的 Content-Type 值。
	Cookies     []*http.Cookie
}

func NewRequest(method, rawURL string, opts ...Option) *Request {
	var nURL, _ = url.Parse(rawURL)
	var req = &Request{}
	req.method = strings.ToUpper(method)
	req.url = nURL
	req.ContentType = ContentTypeURLEncode
	req.Header = http.Header{}
	req.Query = nURL.Query()
	req.Form = url.Values{}
	req.FileForm = FileForm{}

	req.url.RawQuery = ""
	for _, opt := range opts {
		if opt != nil {
			opt(req)
		}
	}
	return req
}

func (r *Request) JoinPath(elems ...string) {
	r.url = r.url.JoinPath(elems...)
}

func (r *Request) AddCookie(cookie *http.Cookie) {
	r.Cookies = append(r.Cookies, cookie)
}

func (r *Request) SetCookies(cookies []*http.Cookie) {
	r.Cookies = cookies
}

func (r *Request) Request(ctx context.Context) (req *http.Request, err error) {
	var body io.Reader
	var forceQuery bool

	if r.method == http.MethodGet ||
		r.method == http.MethodTrace ||
		r.method == http.MethodOptions ||
		r.method == http.MethodHead ||
		r.method == http.MethodDelete {
		forceQuery = true
	}

	var bodyEncoder BodyEncoder
	if r.Body != nil {
		bodyEncoder = r.Body
	} else if len(r.FileForm) > 0 {
		bodyEncoder = multiEncoder()
	} else if len(r.Form) > 0 && !forceQuery {
		bodyEncoder = formEncoder()
	}
	if bodyEncoder != nil {
		body, err = bodyEncoder(r)
		if err != nil {
			return nil, err
		}
	}

	req, err = http.NewRequestWithContext(ctx, r.method, r.url.String(), body)
	if err != nil {
		return nil, err
	}

	var rawQuery = r.Query
	if forceQuery {
		if rawQuery == nil {
			rawQuery = url.Values{}
		}
		for key, values := range r.Form {
			for _, value := range values {
				rawQuery.Add(key, value)
			}
		}
	}
	if len(rawQuery) > 0 {
		req.URL.RawQuery = rawQuery.Encode()
	}

	var header = r.Header
	if header == nil {
		header = http.Header{}
	}
	if r.ContentType != "" {
		header.Set(kContentType, string(r.ContentType))
	}
	req.Header = header

	for _, cookie := range r.Cookies {
		req.AddCookie(cookie)
	}
	return req, nil
}

func (r *Request) Do(ctx context.Context) (*http.Response, error) {
	var req, err = r.Request(ctx)
	if err != nil {
		return nil, err
	}
	var client = r.Client
	if client == nil {
		client = http.DefaultClient
	}
	return client.Do(req)
}
