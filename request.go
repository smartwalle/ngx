package ngx

import (
	"bytes"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type Request struct {
	url     string
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

func NewRequest(method, urlString string) *Request {
	var r = &Request{}
	r.method = strings.ToUpper(method)
	r.url = urlString
	r.params = url.Values{}
	r.header = http.Header{}
	r.Client = http.DefaultClient
	r.SetContentType("application/x-www-form-urlencoded")
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
	this.params = nil
}

func (this *Request) AddParam(key, value string) {
	this.params.Add(key, value)
	this.body = nil
}

func (this *Request) SetParam(key, value string) {
	this.params.Set(key, value)
	this.body = nil
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

func (this *Request) Exec() *Response {
	var req *http.Request
	var err error
	var body io.Reader
	var rawQuery string

	if this.method == http.MethodGet || this.method == http.MethodHead || this.method == http.MethodDelete {
		if len(this.params) > 0 {
			rawQuery = this.params.Encode()
		}
	} else {
		if this.body != nil {
			body = this.body
		} else if this.file != nil {
			uploadFile, err := os.Open(this.file.path)
			if err != nil {
				return &Response{nil, nil, err}
			}
			defer uploadFile.Close()

			bodyByte := &bytes.Buffer{}
			writer := multipart.NewWriter(bodyByte)
			part, err := writer.CreateFormFile(this.file.name, this.file.filename)
			if err != nil {
				return &Response{nil, nil, err}
			}
			_, err = io.Copy(part, uploadFile)
			if err != nil {
				return &Response{nil, nil, err}
			}

			for key, values := range this.params {
				for _, value := range values {
					writer.WriteField(key, value)
				}
			}

			err = writer.Close()
			if err != nil {
				return &Response{nil, nil, err}
			}

			this.SetContentType(writer.FormDataContentType())
			body = bodyByte
		} else if this.params != nil {
			body = strings.NewReader(this.params.Encode())
		}
	}

	req, err = http.NewRequest(this.method, this.url, body)
	if len(rawQuery) > 0 {
		req.URL.RawQuery = rawQuery
	}

	if err != nil {
		return &Response{nil, nil, err}
	}
	req.Header = this.header

	for _, cookie := range this.cookies {
		req.AddCookie(cookie)
	}

	resp, err := this.Client.Do(req)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return &Response{nil, nil, err}
	}

	data, err := ioutil.ReadAll(resp.Body)
	return &Response{resp, data, err}
}
