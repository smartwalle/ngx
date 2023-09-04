package ngx

import (
	"io"
	"mime/multipart"
	"os"
)

type FormFile interface {
	WriteTo(writer *multipart.Writer) error
}

type file struct {
	name     string
	filename string
	filepath string
}

func (f file) WriteTo(writer *multipart.Writer) error {
	nFile, err := os.Open(f.filepath)
	if err != nil {
		return err
	}
	defer nFile.Close()
	nWriter, err := writer.CreateFormFile(f.name, f.filename)
	if err != nil {
		return err
	}
	if _, err = io.Copy(nWriter, nFile); err != nil {
		return err
	}
	return nil
}

type fileObject struct {
	name     string
	filename string
	reader   io.Reader
}

func (f fileObject) WriteTo(writer *multipart.Writer) error {
	if f.reader == nil {
		return nil
	}
	nWriter, err := writer.CreateFormFile(f.name, f.filename)
	if err != nil {
		return err
	}
	if _, err = io.Copy(nWriter, f.reader); err != nil {
		return err
	}
	return nil
}
