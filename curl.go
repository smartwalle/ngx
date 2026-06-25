package ngx

import (
	"bytes"
	"errors"
	"net/http"
	"net/url"
	"strings"

	"github.com/smartwalle/ngx/curl"
)

func (r *Request) CURL() (string, error) {
	if r.err != nil {
		return "", r.err
	}
	if r.url == nil {
		return "", errors.New("ngx: request url is nil; use NewRequest with a valid URL")
	}

	var body *bytes.Buffer
	var contentType = r.ContentType
	var forceQuery bool

	if r.method == http.MethodGet ||
		r.method == http.MethodTrace ||
		r.method == http.MethodOptions ||
		r.method == http.MethodHead ||
		r.method == http.MethodDelete {
		forceQuery = true
	}

	var bodyEncoder BodyEncoder
	if r.Body != nil {
		bodyEncoder = r.Body
	} else if len(r.File) > 0 {
		// File is translated to curl --form later so paths remain visible.
	} else if len(r.Form) > 0 && !forceQuery {
		bodyEncoder = formEncoder()
	}
	if bodyEncoder != nil {
		var reader, nContentType, err = bodyEncoder(r)
		if err != nil {
			return "", err
		}
		if nContentType != "" {
			contentType = nContentType
		}
		if reader != nil {
			body = &bytes.Buffer{}
			if _, err = body.ReadFrom(reader); err != nil {
				return "", err
			}
		}
	}

	var rawURL = *r.url
	var rawQuery = CloneValues(r.Query)
	if forceQuery {
		if rawQuery == nil {
			rawQuery = url.Values{}
		}
		for key, values := range r.Form {
			for _, value := range values {
				rawQuery.Add(key, value)
			}
		}
	}
	if len(rawQuery) > 0 {
		rawURL.RawQuery = rawQuery.Encode()
	}

	var cmd = curl.New(r.method, rawURL.String())
	var useMultipartForm bool
	if len(r.File) > 0 && r.Body == nil {
		if err := addCurlForm(cmd, r); err != nil {
			return "", err
		}
		useMultipartForm = true
	} else if body != nil && body.Len() > 0 {
		cmd.Data(body.String())
	}

	var header = CloneHeader(r.Header)
	if header == nil {
		header = http.Header{}
	}
	if !useMultipartForm && contentType != "" {
		header.Set(kContentType, contentType)
	}
	for key, values := range header {
		for _, value := range values {
			cmd.Header(key, value)
		}
	}
	addCurlCookies(cmd, r.Cookies)

	return cmd.Encode(), nil
}

func addCurlForm(cmd *curl.Command, req *Request) error {
	for _, file := range req.File {
		if _, ok := file.(fileInfo); !ok {
			return errors.New("ngx: cannot convert file object to curl command; file path is required")
		}
	}
	for key, file := range req.File {
		var info = file.(fileInfo)
		cmd.File(key, info.filepath, info.filename)
	}
	for key, values := range req.Form {
		for _, value := range values {
			cmd.Form(key, value)
		}
	}
	return nil
}

func addCurlCookies(cmd *curl.Command, cookies []*http.Cookie) {
	if len(cookies) == 0 {
		return
	}
	var values = make([]string, 0, len(cookies))
	for _, cookie := range cookies {
		if cookie == nil {
			continue
		}
		if value := cookie.String(); value != "" {
			values = append(values, value)
		}
	}
	if len(values) > 0 {
		cmd.CookieRaw(strings.Join(values, "; "))
	}
}
