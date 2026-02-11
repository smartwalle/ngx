package ngx

import (
	"io"
	"mime/multipart"
	"os"
)

type FileForm map[string]File

func (f FileForm) AddFile(key, filename, filepath string) {
	f.Add(key, fileInfo{filename: filename, filepath: filepath})
}

func (f FileForm) AddObject(key, filename string, file io.Reader) {
	f.Add(key, fileObject{filename: filename, reader: file})
}

func (f FileForm) Add(key string, file File) {
	f[key] = file
}

func (f FileForm) Del(key string) {
	delete(f, key)
}

func (f FileForm) Has(key string) bool {
	_, ok := f[key]
	return ok
}

type File interface {
	Write(key string, writer *multipart.Writer) error
}

type fileInfo struct {
	filename string
	filepath string
}

func (f fileInfo) Write(key string, writer *multipart.Writer) error {
	file, err := os.Open(f.filepath)
	if err != nil {
		return err
	}
	defer file.Close()
	fileWriter, err := writer.CreateFormFile(key, f.filename)
	if err != nil {
		return err
	}
	if _, err = io.Copy(fileWriter, file); err != nil {
		return err
	}
	return nil
}

type fileObject struct {
	filename string
	reader   io.Reader
}

func (f fileObject) Write(key string, writer *multipart.Writer) error {
	if f.reader == nil {
		return nil
	}
	fileWriter, err := writer.CreateFormFile(key, f.filename)
	if err != nil {
		return err
	}
	if _, err = io.Copy(fileWriter, f.reader); err != nil {
		return err
	}
	return nil
}
