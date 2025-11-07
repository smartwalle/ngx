package ngx

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"strings"
)

type BodyEncoder interface {
	Encode(req *Request) (io.Reader, error)
}

func Body(r io.Reader) BodyEncoder {
	return rawEncoder{r: r}
}

type rawEncoder struct {
	r io.Reader
}

func (e rawEncoder) Encode(req *Request) (io.Reader, error) {
	return e.r, nil
}

func JSONEncoder(v interface{}) BodyEncoder {
	return jsonEncoder{v: v}
}

type jsonEncoder struct {
	v interface{}
}

func (e jsonEncoder) Encode(req *Request) (io.Reader, error) {
	var buffer = &bytes.Buffer{}
	if err := json.NewEncoder(buffer).Encode(e.v); err != nil {
		return nil, err
	}
	req.ContentType = ContentTypeJSON
	return buffer, nil
}

type multiEncoder struct {
}

func (e multiEncoder) Encode(req *Request) (io.Reader, error) {
	var multiBuffer = &bytes.Buffer{}
	var multiWriter = multipart.NewWriter(multiBuffer)
	for key, file := range req.FileForm {
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
	req.ContentType = ContentType(multiWriter.FormDataContentType())
	return multiBuffer, nil
}

type formEncoder struct {
}

func (e formEncoder) Encode(req *Request) (io.Reader, error) {
	return strings.NewReader(req.Form.Encode()), nil
}
