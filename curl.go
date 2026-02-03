package ngx

import (
	"bytes"
	"io"
	"net/http"

	"github.com/smartwalle/ngx/curl"
)

func CURL(req *http.Request) (string, error) {
	var cmd = curl.New(req.Method, req.URL.String())

	for key, values := range req.Header {
		for _, value := range values {
			cmd.Header(key, value)
		}
	}

	var userAgent = req.UserAgent()
	if userAgent != "" {
		cmd.UserAgent(userAgent)
	}

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
		cmd.Data(bodyBuffer.String())
	}
	return cmd.Encode(), nil
}
