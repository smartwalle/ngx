package ngx

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type Response struct {
	error    error
	response *http.Response
	data     []byte
}

func (this *Response) Status() string {
	if this.response != nil {
		return this.response.Status
	}
	return fmt.Sprintf("%d ServiceUnavailable", http.StatusServiceUnavailable)
}

func (this *Response) StatusCode() int {
	if this.response != nil {
		return this.response.StatusCode
	}
	return http.StatusServiceUnavailable
}

func (this *Response) Proto() string {
	if this.response != nil {
		return this.response.Proto
	}
	return ""
}

func (this *Response) ProtoMajor() int {
	if this.response != nil {
		return this.response.ProtoMajor
	}
	return 1
}

func (this *Response) ProtoMinor() int {
	if this.response != nil {
		return this.response.ProtoMinor
	}
	return 0
}

func (this *Response) Header() http.Header {
	if this.response != nil {
		return this.response.Header
	}
	return http.Header{}
}

func (this *Response) ContentLength() int64 {
	if this.response != nil {
		return this.response.ContentLength
	}
	return 0
}

func (this *Response) TransferEncoding() []string {
	if this.response != nil {
		return this.response.TransferEncoding
	}
	return nil
}

func (this *Response) Close() bool {
	if this.response != nil {
		return this.response.Close
	}
	return true
}

func (this *Response) Uncompressed() bool {
	if this.response != nil {
		return this.response.Uncompressed
	}
	return true
}

func (this *Response) Trailer() http.Header {
	if this.response != nil {
		return this.response.Trailer
	}
	return http.Header{}
}

func (this *Response) Request() *http.Request {
	if this.response != nil {
		return this.response.Request
	}
	return nil
}

func (this *Response) TLS() *tls.ConnectionState {
	if this.response != nil {
		return this.response.TLS
	}
	return nil
}

func (this *Response) Cookies() []*http.Cookie {
	if this.response != nil {
		return this.response.Cookies()
	}
	return nil
}

func (this *Response) Location() (*url.URL, error) {
	if this.response != nil {
		return this.response.Location()
	}
	return nil, nil
}

func (this *Response) ProtoAtLeast(major, minor int) bool {
	if this.response != nil {
		return this.response.ProtoAtLeast(major, minor)
	}
	return false
}

func (this *Response) Error() error {
	return this.error
}

func (this *Response) Reader() io.Reader {
	return bytes.NewReader(this.data)
}

func (this *Response) Bytes() []byte {
	return this.data
}

func (this *Response) String() string {
	return string(this.data)
}

func (this *Response) UnmarshalJSON(v interface{}) error {
	if this.error != nil {
		return this.error
	}
	return json.Unmarshal(this.data, v)
}
