package ngx

import (
	"io"
	"mime/multipart"
	"os"
)

type FileForm map[string]File

func (f FileForm) AddFilePath(key, filename, filepath string) {
	f.Add(key, FilePath{Filename: filename, Filepath: filepath})
}

func (f FileForm) AddFileReader(key, filename string, file io.Reader) {
	f.Add(key, FileReader{Filename: filename, Reader: file})
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

type FilePath struct {
	Filename string
	Filepath string
}

func (f FilePath) Write(key string, writer *multipart.Writer) error {
	file, err := os.Open(f.Filepath)
	if err != nil {
		return err
	}
	defer file.Close()
	fileWriter, err := writer.CreateFormFile(key, f.Filename)
	if err != nil {
		return err
	}
	if _, err = io.Copy(fileWriter, file); err != nil {
		return err
	}
	return nil
}

type FileReader struct {
	Filename string
	Reader   io.Reader
}

func (f FileReader) Write(key string, writer *multipart.Writer) error {
	if f.Reader == nil {
		return nil
	}
	fileWriter, err := writer.CreateFormFile(key, f.Filename)
	if err != nil {
		return err
	}
	if _, err = io.Copy(fileWriter, f.Reader); err != nil {
		return err
	}
	return nil
}
