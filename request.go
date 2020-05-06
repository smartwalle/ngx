package ngx

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"strings"
)

// --------------------------------------------------------------------------------
const (
	ContentTypeJSON      = "application/json"
	ContentTypeXML       = "application/xml"
	ContentTypeForm      = "application/x-www-form-urlencoded"
	ContentTypeFormData  = "application/x-www-form-urlencoded"
	ContentTypeURLEncode = "application/x-www-form-urlencoded"
	ContentTypeHTML      = "text/html"
	ContentTypeText      = "text/plain"
	ContentTypeMultipart = "multipart/form-data"
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

// --------------------------------------------------------------------------------
type Request struct {
	target  string
	method  string
	header  http.Header
	params  url.Values
	body    io.Reader
	Client  *http.Client
	cookies []*http.Cookie
	file    *file
}

type file struct {
	name     string
	filename string
	path     string
}

// --------------------------------------------------------------------------------
func NewRequest(method, target string) *Request {
	var r = &Request{}
	r.method = strings.ToUpper(method)
	r.target = target
	r.params = url.Values{}
	r.header = http.Header{}
	r.Client = http.DefaultClient
	r.SetContentType(ContentTypeURLEncode)
	return r
}

func NewRequestWithJSON(method, target string, param interface{}) *Request {
	var r = &Request{}
	r.method = strings.ToUpper(method)
	r.target = target
	r.params = url.Values{}
	r.header = http.Header{}
	r.Client = http.DefaultClient
	r.WriteJSON(param)
	return r
}

func (this *Request) SetContentType(contentType string) {
	this.SetHeader("Content-Type", contentType)
}

func (this *Request) AddHeader(key, value string) {
	this.header.Add(key, value)
}

func (this *Request) SetHeader(key, value string) {
	this.header.Set(key, value)
}

func (this *Request) SetHeaders(header http.Header) {
	this.header = header
}

func (this *Request) SetBody(body io.Reader) {
	this.body = body
}

func (this *Request) AddParam(key, value string) {
	this.params.Add(key, value)
}

func (this *Request) SetParam(key, value string) {
	this.params.Set(key, value)
}

func (this *Request) SetParams(params url.Values) {
	this.params = params
}

func (this *Request) AddFile(name, filename, path string) {
	this.file = &file{name, filename, path}
}

func (this *Request) RemoveFile() {
	this.file = nil
}

func (this *Request) AddCookie(cookie *http.Cookie) {
	this.cookies = append(this.cookies, cookie)
}

func (this *Request) Do() (*http.Response, error) {
	return this.DoWithContext(nil)
}

func (this *Request) DoWithContext(ctx context.Context) (*http.Response, error) {
	var req *http.Request
	var err error
	var body io.Reader
	var rawQuery string

	if this.method == http.MethodGet || this.method == http.MethodHead || this.method == http.MethodDelete {
		if len(this.params) > 0 {
			rawQuery = this.params.Encode()
		}
		if this.body != nil {
			body = this.body
		}
	} else {
		if this.body != nil {
			body = this.body
		} else if this.file != nil {
			uploadFile, err := os.Open(this.file.path)
			if err != nil {
				return nil, err
			}
			defer uploadFile.Close()

			bodyByte := &bytes.Buffer{}
			writer := multipart.NewWriter(bodyByte)
			part, err := writer.CreateFormFile(this.file.name, this.file.filename)
			if err != nil {
				return nil, err
			}
			_, err = io.Copy(part, uploadFile)
			if err != nil {
				return nil, err
			}

			for key, values := range this.params {
				for _, value := range values {
					writer.WriteField(key, value)
				}
			}

			if err = writer.Close(); err != nil {
				return nil, err
			}

			this.SetContentType(writer.FormDataContentType())
			body = bodyByte
		} else if len(this.params) > 0 {
			body = strings.NewReader(this.params.Encode())
		}
	}

	req, err = http.NewRequest(this.method, this.target, body)
	if ctx != nil && req != nil {
		req = req.WithContext(ctx)
	}
	if len(rawQuery) > 0 {
		req.URL.RawQuery = rawQuery
	}

	if err != nil {
		return nil, err
	}
	req.Header = this.header

	for _, cookie := range this.cookies {
		req.AddCookie(cookie)
	}

	return this.Client.Do(req)
}

func (this *Request) Exec() *Response {
	return this.ExecWithContext(nil)
}

func (this *Request) ExecWithContext(ctx context.Context) *Response {
	rsp, err := this.DoWithContext(ctx)
	if rsp != nil {
		defer rsp.Body.Close()
	}
	if err != nil {
		return &Response{rsp, nil, err}
	}
	data, err := ioutil.ReadAll(rsp.Body)
	return &Response{rsp, data, err}
}

func (this *Request) Download(savePath string) *Response {
	return this.DownloadWithContext(nil, savePath)
}

func (this *Request) DownloadWithContext(ctx context.Context, savePath string) *Response {
	rsp, err := this.DoWithContext(ctx)
	if rsp != nil {
		defer rsp.Body.Close()
	}
	if err != nil {
		return &Response{rsp, nil, err}
	}

	nFile, err := os.Create(savePath)
	if err != nil {
		return &Response{nil, nil, err}
	}
	defer nFile.Close()

	buf := make([]byte, 1024)
	for {
		size, _ := rsp.Body.Read(buf)
		if size == 0 {
			break
		}
		nFile.Write(buf[:size])
	}
	data := []byte(savePath)
	return &Response{rsp, data, err}
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
