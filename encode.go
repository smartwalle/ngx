package ngx

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"strings"
)

type BodyEncoder func(req *Request) (io.Reader, error)

func Body(r io.Reader) BodyEncoder {
	return func(req *Request) (io.Reader, error) {
		return r, nil
	}
}

func JSONEncoder(v interface{}) BodyEncoder {
	return func(req *Request) (io.Reader, error) {
		var buffer = &bytes.Buffer{}
		if err := json.NewEncoder(buffer).Encode(v); err != nil {
			return nil, err
		}
		req.ContentType = ContentTypeJSON
		return buffer, nil
	}
}

func multiEncoder() BodyEncoder {
	return func(req *Request) (io.Reader, error) {
		var multiBuffer = &bytes.Buffer{}
		var multiWriter = multipart.NewWriter(multiBuffer)
		for key, file := range req.File {
			if err := file.Write(key, multiWriter); err != nil {
				return nil, err
			}
		}
		for key, values := range req.Form {
			for _, value := range values {
				if err := multiWriter.WriteField(key, value); err != nil {
					return nil, err
				}
			}
		}
		if err := multiWriter.Close(); err != nil {
			return nil, err
		}
		req.ContentType = multiWriter.FormDataContentType()
		return multiBuffer, nil
	}
}

func formEncoder() BodyEncoder {
	return func(req *Request) (io.Reader, error) {
		return strings.NewReader(req.Form.Encode()), nil
	}
}
