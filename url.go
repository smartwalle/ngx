package ngx

import (
	"net/url"
)

type URL struct {
	url    *url.URL
	params url.Values
}

func NewURL(rawURL string) (u *URL, err error) {
	unescape, err := url.QueryUnescape(rawURL)
	if err != nil {
		return nil, err
	}
	newURL, err := url.Parse(unescape)
	if err != nil {
		return nil, err
	}

	u = &URL{}
	u.url = newURL
	u.params = u.url.Query()
	return u, nil
}

func MustURL(rawURL string) (u *URL) {
	u, err := NewURL(rawURL)
	if err != nil {
		panic(err)
	}
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
	return this.params.Get(key)
}

func (this *URL) Query() url.Values {
	return this.params
}

func (this *URL) RawURL() *url.URL {
	return this.url
}

func URLEncode(s string) string {
	s = url.QueryEscape(s)
	return s
}

func URLDecode(s string) string {
	s, _ = url.QueryUnescape(s)
	return s
}
