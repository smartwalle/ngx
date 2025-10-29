package ngx

import (
	"bytes"
	"io"
	"net/http"
	"strings"
)

func CURL(req *http.Request) (string, error) {
	var buffer = &bytes.Buffer{}

	buffer.WriteString("curl")
	// Method
	buffer.WriteString(" --request ")
	buffer.WriteString(req.Method)

	// URL
	buffer.WriteString(" ")
	escape(buffer, req.URL.String())

	// Header
	for key, values := range req.Header {
		for _, value := range values {
			buffer.WriteString(" --header ")
			escape(buffer, key, ":", value)
		}
	}

	// Body
	var err error
	var body io.Reader

	body, req.Body, err = DrainBody(req.Body)
	if err != nil {
		return "", err
	}
	var bodyBuffer bytes.Buffer
	if _, err = bodyBuffer.ReadFrom(body); err != nil {
		return "", err
	}
	if bodyBuffer.Len() > 0 {
		buffer.WriteString(" --data ")
		escape(buffer, bodyBuffer.String())
	}

	return buffer.String(), nil
}

func escape(buffer *bytes.Buffer, values ...string) {
	buffer.WriteString("'")
	for _, value := range values {
		buffer.WriteString(strings.Replace(value, "'", "'\\''", -1))
	}
	buffer.WriteString("'")
}
