package ngx

import (
	"io"
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

func WithBody(body io.Reader) Option {
	return func(req *Request) {
		req.body = body
	}
}

func WithCookies(cookies []*http.Cookie) Option {
	return func(req *Request) {
		req.cookies = cookies
	}
}

func CloneValues(src url.Values) url.Values {
	if src == nil {
		return nil
	}

	nv := 0
	for _, vv := range src {
		nv += len(vv)
	}
	sv := make([]string, nv)
	dst := make(url.Values, len(src))
	for k, vv := range src {
		if vv == nil {
			dst[k] = nil
			continue
		}
		n := copy(sv, vv)
		dst[k] = sv[:n:n]
		sv = sv[n:]
	}
	return dst
}

func CloneHeader(src http.Header) http.Header {
	return src.Clone()
}
