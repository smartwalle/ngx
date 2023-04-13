package ngx

import (
	"io"
)

type writer struct {
	io.Writer
	total    uint64
	finished uint64
	handler  func(total uint64, finished uint64)
}

func NewWriter(w io.Writer, total uint64, handler func(total uint64, finished uint64)) *writer {
	return &writer{Writer: w, total: total, handler: handler}
}

func (this *writer) Write(p []byte) (n int, err error) {
	n, err = this.Writer.Write(p)

	if n > 0 {
		this.finished += uint64(n)
		if this.handler != nil {
			this.handler(this.total, this.finished)
		}
	}

	return n, err
}
