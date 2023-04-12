package ngx

type receive struct {
	total    uint64
	received uint64
	handler  func(total uint64, received uint64)
}

func (this *receive) Write(p []byte) (int, error) {
	var n = len(p)
	this.received += uint64(n)
	if this.handler != nil {
		this.handler(this.total, this.received)
	}
	return n, nil
}
