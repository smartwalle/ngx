package ngx

import (
	"io"
	"net/http"
	"net/url"
)

type Option func(req *Request)

func WithClient(client *http.Client) Option {
	return func(req *Request) {
		req.Client = client
	}
}

func WithHeader(header http.Header) Option {
	return func(req *Request) {
		for key, values := range header {
			for _, value := range values {
				req.Header.Add(key, value)
			}
		}
	}
}

func WithForm(form url.Values) Option {
	return func(req *Request) {
		for key, values := range form {
			for _, value := range values {
				req.Form.Add(key, value)
			}
		}
	}
}

func WithQuery(query url.Values) Option {
	return func(req *Request) {
		for key, values := range query {
			for _, value := range values {
				req.Query.Add(key, value)
			}
		}
	}
}

func WithBody(body io.Reader) Option {
	return func(req *Request) {
		req.Body = Body(body)
	}
}

func WithCookies(cookies []*http.Cookie) Option {
	return func(req *Request) {
		req.Cookies = cookies
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
