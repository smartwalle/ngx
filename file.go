package ngx

import (
	"io"
	"mime/multipart"
	"os"
)

type FileForm map[string]File

func (f FileForm) AddFilePath(name, filename, filepath string) {
	if filename == "" {
		filename = name
	}
	f.Add(name, file{name: name, filename: filename, filepath: filepath})
}

func (f FileForm) AddFileObject(name, filename string, file io.Reader) {
	if filename == "" {
		filename = name
	}
	f.Add(name, fileObject{name: name, filename: filename, reader: file})
}

func (f FileForm) Add(name string, file File) {
	f[name] = file
}

func (f FileForm) Del(name string) {
	delete(f, name)
}

func (f FileForm) Has(name string) bool {
	_, ok := f[name]
	return ok
}

type File interface {
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
