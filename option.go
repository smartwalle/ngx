package ngx

import (
	"net/http"
	"net/url"
)

type Option func(req *Request)

func WithClient(client *http.Client) Option {
	return func(req *Request) {
		req.client = client
	}
}

func WithHeader(header http.Header) Option {
	return func(req *Request) {
		req.header = header
	}
}

func WithForm(form url.Values) Option {
	return func(req *Request) {
		req.form = form
	}
}

func WithQuery(query url.Values) Option {
	return func(req *Request) {
		req.query = query
	}
}

func WithBody(body Body) Option {
	return func(req *Request) {
		req.body = body
	}
}

func WithCookies(cookies []*http.Cookie) Option {
	return func(req *Request) {
		req.cookies = cookies
	}
}

// WithReceive 获取从服务端已接收数据大小
func WithReceive(fn func(total, chunk, finished uint64)) Option {
	return func(req *Request) {
		req.receive = fn
	}
}

// WithSend 获取向服务端已发送数据大小
func WithSend(fn func(total, chunk, finished uint64)) Option {
	return func(req *Request) {
		req.send = fn
	}
}
