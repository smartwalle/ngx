package ngx

import (
	"encoding/json"
	"net/http"
	"crypto/tls"
	"net/url"
	"io"
	"fmt"
)

type Response struct {
	rsp   *http.Response
	data  []byte
	error error
}

func (this *Response) Status() string {
	if this.rsp != nil {
		return this.rsp.Status
	}
	return fmt.Sprintf("%d ServiceUnavailable", http.StatusServiceUnavailable)
}

func (this *Response) StatusCode() int {
	if this.rsp != nil {
		return this.rsp.StatusCode
	}
	return http.StatusServiceUnavailable
}

func (this *Response) Proto() string {
	if this.rsp != nil {
		return this.rsp.Proto
	}
	return ""
}

func (this *Response) ProtoMajor() int {
	if this.rsp != nil {
		return this.rsp.ProtoMajor
	}
	return 1
}

func (this *Response) ProtoMinor() int {
	if this.rsp != nil {
		return this.rsp.ProtoMinor
	}
	return 0
}

func (this *Response) Header() http.Header {
	if this.rsp != nil {
		return this.rsp.Header
	}
	return http.Header{}
}

func (this *Response) ContentLength() int64 {
	if this.rsp != nil {
		return this.rsp.ContentLength
	}
	return 0
}

func (this *Response) TransferEncoding() []string {
	if this.rsp != nil {
		return this.rsp.TransferEncoding
	}
	return nil
}

func (this *Response) Close() bool {
	if this.rsp != nil {
		return this.rsp.Close
	}
	return true
}

func (this *Response) Uncompressed() bool {
	if this.rsp != nil {
		return this.rsp.Uncompressed
	}
	return true
}

func (this *Response) Trailer() http.Header {
	if this.rsp != nil {
		return this.rsp.Trailer
	}
	return http.Header{}
}

func (this *Response) Request() *http.Request {
	if this.rsp != nil {
		return this.rsp.Request
	}
	return nil
}

func (this *Response) TLS() *tls.ConnectionState {
	if this.rsp != nil {
		return this.rsp.TLS
	}
	return nil
}

func (this *Response) Cookies() []*http.Cookie {
	if this.rsp != nil {
		return this.rsp.Cookies()
	}
	return nil
}

func (this *Response) Location() (*url.URL, error) {
	if this.rsp != nil {
		return this.rsp.Location()
	}
	return nil, nil
}

func (this *Response) ProtoAtLeast(major, minor int) bool {
	if this.rsp != nil {
		return this.rsp.ProtoAtLeast(major, minor)
	}
	return false
}

func (this *Response) Write(w io.Writer) error {
	if this.rsp != nil {
		return this.rsp.Write(w)
	}
	return nil
}

func (this *Response) Error() error {
	return this.error
}

func (this *Response) Bytes() ([]byte, error) {
	return this.data, this.error
}

func (this *Response) MustBytes() []byte {
	return this.data
}

func (this *Response) String() (string, error) {
	return string(this.data), this.error
}

func (this *Response) MustString() string {
	return string(this.data)
}

func (this *Response) UnmarshalJSON(v interface{}) error {
	if this.error != nil {
		return this.error
	}
	return json.Unmarshal(this.data, v)
}
