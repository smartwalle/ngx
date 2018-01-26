package ngx

import (
	"net/http"
	"encoding/json"
)

type Response struct {
	*http.Response
	data  []byte
	error error
}

func (this *Response) Error() error {
	return this.error
}

func (this *Response) Bytes() ([]byte, error) {
	return this.data, this.error
}

func (this *Response) MustBytes() ([]byte) {
	return this.data
}

func (this *Response) String() (string, error) {
	return string(this.data), this.error
}

func (this *Response) MustString() (string) {
	return string(this.data)
}

func (this *Response) UnmarshalJSON(v interface{}) (error) {
	if this.error != nil {
		return this.error
	}
	return json.Unmarshal(this.data, v)
}
