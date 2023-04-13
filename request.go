package ngx

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type ContentType string

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
	target  string
	method  string
	header  http.Header
	params  url.Values
	query   url.Values
	body    Body
	Client  *http.Client
	cookies []*http.Cookie
	files   map[string]file
	receive func(total uint64, finished uint64)
	send    func(total uint64, finished uint64)
}

type file struct {
	name     string
	filename string
	filepath string
}

func NewRequest(method, target string, opts ...Option) *Request {
	var nURL, _ = url.Parse(target)
	var req = &Request{}
	req.method = strings.ToUpper(method)
	req.target = target
	req.params = url.Values{}
	if nURL != nil {
		req.query = nURL.Query()
	} else {
		req.query = url.Values{}
	}
	req.header = http.Header{}
	req.Client = http.DefaultClient
	req.SetContentType(ContentTypeURLEncode)

	for _, opt := range opts {
		if opt != nil {
			opt(req)
		}
	}
	return req
}

func NewJSONRequest(method, target string, param interface{}, opts ...Option) *Request {
	var r = NewRequest(method, target, opts...)
	r.WriteJSON(param)
	return r
}

// WriteJSON 将一个对象序列化为 JSON 字符串，并将其作为 http 请求的 body 发送给服务器端。
func (this *Request) WriteJSON(v interface{}) error {
	data, err := json.Marshal(v)
	if err != nil {
		return err
	}
	this.SetBody(bytes.NewReader(data))
	this.SetContentType(ContentTypeJSON)
	return nil
}

func (this *Request) SetContentType(contentType ContentType) {
	this.SetHeader("Content-Type", string(contentType))
}

func (this *Request) AddHeader(key, value string) {
	this.header.Add(key, value)
}

func (this *Request) DelHeader(key string) {
	this.header.Del(key)
}

func (this *Request) SetHeader(key, value string) {
	this.header.Set(key, value)
}

func (this *Request) SetHeaders(header http.Header) {
	this.header = header
}

func (this *Request) SetBody(body Body) {
	this.body = body
}

func (this *Request) AddParam(key, value string) {
	this.params.Add(key, value)
}

func (this *Request) DelParam(key string) {
	this.params.Del(key)
}

func (this *Request) SetParam(key, value string) {
	this.params.Set(key, value)
}

func (this *Request) SetParams(params url.Values) {
	this.params = params
}

func (this *Request) AddQuery(key, value string) {
	this.query.Add(key, value)
}

func (this *Request) DelQuery(key string) {
	this.query.Del(key)
}

func (this *Request) SetQuery(key, value string) {
	this.query.Set(key, value)
}

func (this *Request) AddFile(name, filename, filepath string) {
	if this.files == nil {
		this.files = make(map[string]file)
	}
	if filename == "" {
		filename = name
	}
	this.files[name] = file{name, filename, filepath}
}

func (this *Request) DelFile(name string) {
	if this.files != nil {
		delete(this.files, name)
	}
}

func (this *Request) AddCookie(cookie *http.Cookie) {
	this.cookies = append(this.cookies, cookie)
}

func (this *Request) do(ctx context.Context) (*http.Response, error) {
	var req *http.Request
	var err error
	var body Body
	var toQuery bool

	if this.method == http.MethodGet ||
		this.method == http.MethodTrace ||
		this.method == http.MethodOptions ||
		this.method == http.MethodHead ||
		this.method == http.MethodDelete {
		toQuery = true
	}

	if this.body != nil {
		body = this.body
	} else if len(this.files) > 0 {
		var bodyBuffer = &bytes.Buffer{}
		var bodyWriter = multipart.NewWriter(bodyBuffer)
		var fileContent []byte
		var fileWriter io.Writer

		for _, f := range this.files {
			fileContent, err = os.ReadFile(f.filepath)
			if err != nil {
				return nil, err
			}
			fileWriter, err = bodyWriter.CreateFormFile(f.name, f.filename)
			if err != nil {
				return nil, err
			}
			if _, err = fileWriter.Write(fileContent); err != nil {
				return nil, err
			}
		}
		for key, values := range this.params {
			for _, value := range values {
				bodyWriter.WriteField(key, value)
			}
		}

		if err = bodyWriter.Close(); err != nil {
			return nil, err
		}

		this.SetContentType(ContentType(bodyWriter.FormDataContentType()))
		body = bodyBuffer
	} else if len(this.params) > 0 && !toQuery {
		body = strings.NewReader(this.params.Encode())
	}

	var getBody func() (io.ReadCloser, error)
	var contentLength int64

	if body != nil {
		switch v := body.(type) {
		case *bytes.Buffer:
			contentLength = int64(v.Len())
			buf := v.Bytes()
			getBody = func() (io.ReadCloser, error) {
				r := bytes.NewReader(buf)
				return io.NopCloser(NewReader(r, this.send)), nil
			}
		case *bytes.Reader:
			contentLength = int64(v.Len())
			snapshot := *v
			getBody = func() (io.ReadCloser, error) {
				r := snapshot
				return io.NopCloser(NewReader(&r, this.send)), nil
			}
		case *strings.Reader:
			contentLength = int64(v.Len())
			snapshot := *v
			getBody = func() (io.ReadCloser, error) {
				r := snapshot
				return io.NopCloser(NewReader(&r, this.send)), nil
			}
		default:
		}

		if getBody != nil && contentLength == 0 {
			getBody = func() (io.ReadCloser, error) { return http.NoBody, nil }
		}

		body = NewReader(body, this.send)
	}

	req, err = http.NewRequestWithContext(ctx, this.method, this.target, body)
	if err != nil {
		return nil, err
	}

	req.ContentLength = contentLength
	req.GetBody = getBody

	if toQuery {
		for key, values := range this.params {
			for _, value := range values {
				this.query.Add(key, value)
			}
		}
	}

	req.URL.RawQuery = this.query.Encode()
	req.Header = this.header

	for _, cookie := range this.cookies {
		req.AddCookie(cookie)
	}

	return this.Client.Do(req.Clone(ctx))
}

func (this *Request) exec(rsp *http.Response, w io.Writer) error {
	var nWriter = NewWriter(w, uint64(rsp.ContentLength), this.receive)
	if _, err := io.Copy(nWriter, rsp.Body); err != nil {
		return err
	}
	return nil
}

func (this *Request) Exec(ctx context.Context) *Response {
	rsp, err := this.do(ctx)
	if rsp != nil {
		defer rsp.Body.Close()
	}
	if err != nil {
		return &Response{rsp, nil, err}
	}

	var w = bytes.NewBuffer(nil)

	if err = this.exec(rsp, w); err != nil {
		return &Response{rsp, nil, err}
	}

	return &Response{rsp, w.Bytes(), err}
}

func (this *Request) Download(ctx context.Context, filepath string) *Response {
	rsp, err := this.do(ctx)
	if rsp != nil {
		defer rsp.Body.Close()
	}
	if err != nil {
		return &Response{rsp, nil, err}
	}

	w, err := os.Create(filepath)
	if err != nil {
		return &Response{nil, nil, err}
	}
	defer w.Close()

	if err = this.exec(rsp, w); err != nil {
		return &Response{rsp, nil, err}
	}

	return &Response{rsp, []byte(filepath), err}
}
