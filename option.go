package ngx

type Option func(opts *Request)

func WithReceived(f func(total uint64, received uint64)) Option {
	return func(req *Request) {
		req.received = f
	}
}
