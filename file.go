package ngx

import (
	"io"
	"mime/multipart"
	"os"
)

type FileForm map[string]File

func (f FileForm) AddFilePath(name, filename, filepath string) {
	f.Add(name, fileInfo{filename: filename, filepath: filepath})
}

func (f FileForm) AddFileObject(name, filename string, file io.Reader) {
	f.Add(name, fileObject{filename: filename, reader: file})
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
	Write(name string, writer *multipart.Writer) error
}

type fileInfo struct {
	filename string
	filepath string
}

func (f fileInfo) Write(name string, writer *multipart.Writer) error {
	file, err := os.Open(f.filepath)
	if err != nil {
		return err
	}
	defer file.Close()
	fileWriter, err := writer.CreateFormFile(name, f.filename)
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

func (f fileObject) Write(name string, writer *multipart.Writer) error {
	if f.reader == nil {
		return nil
	}
	fileWriter, err := writer.CreateFormFile(name, f.filename)
	if err != nil {
		return err
	}
	if _, err = io.Copy(fileWriter, f.reader); err != nil {
		return err
	}
	return nil
}
