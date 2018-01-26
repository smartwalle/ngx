package ngx

import (
	"net/url"
)

type URL struct {
	url    *url.URL
	params url.Values
}

func NewURL(rawURL string) *URL {
	unescape, _ := url.QueryUnescape(rawURL)
	var u = &URL{}
	u.url, _ = url.Parse(unescape)
	u.params = u.url.Query()
	return u
}

func (this *URL) String() string {
	this.url.RawQuery = this.params.Encode()
	return this.url.String()
}

func (this *URL) Add(key, value string) {
	this.params.Add(key, value)
}

func (this *URL) Del(key string) {
	this.params.Del(key)
}

func (this *URL) Set(key, value string) {
	this.params.Set(key, value)
}

func (this *URL) Get(key string) string {
	return this.Get(key)
}

func (this *URL) Query() url.Values {
	return this.params
}
