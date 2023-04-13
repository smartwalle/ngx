package ngx

import "net/http"

type Option func(req *Request)

// WithReceive 获取从服务端已接收数据大小
func WithReceive(f func(total uint64, finished uint64)) Option {
	return func(req *Request) {
		req.receive = f
	}
}

// WithSend 获取向服务端已发送数据大小
func WithSend(f func(total uint64, finished uint64)) Option {
	return func(req *Request) {
		req.send = f
	}
}

func WithClient(c *http.Client) Option {
	return func(req *Request) {
		if c != nil {
			req.Client = c
		}
	}
}
