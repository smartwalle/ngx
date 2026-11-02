package ngx

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

// Picker 根据服务名称选择服务端点。
type Picker interface {
	Pick(ctx context.Context, service string) (string, error)
}

// Transport 实现基于服务发现的 HTTP 请求级负载均衡。
type Transport struct {
	base   *http.Transport
	picker Picker
}

// NewTransport 创建服务发现 Transport。
func NewTransport(base *http.Transport) *Transport {
	if base == nil {
		defaultTransport, ok := http.DefaultTransport.(*http.Transport)
		if !ok {
			panic("ngx: http.DefaultTransport is not *http.Transport")
		}
		base = defaultTransport
	}
	return &Transport{
		base: base.Clone(),
	}
}

// UsePicker 设置负载均衡器。
func (t *Transport) UsePicker(picker Picker) {
	t.picker = picker
}

// RoundTrip 实现请求级负载均衡。
func (t *Transport) RoundTrip(req *http.Request) (*http.Response, error) {
	if t.picker == nil {
		return t.base.RoundTrip(req)
	}

	var reqScheme = strings.ToLower(req.URL.Scheme)
	if reqScheme != "http" && reqScheme != "https" {
		return t.base.RoundTrip(req)
	}

	var service = normalizeHost(req.URL.Hostname())

	selected, err := t.picker.Pick(req.Context(), service)
	if err != nil {
		return nil, err
	}

	endpointURL, err := parseEndpointURL(selected)
	if err != nil {
		return nil, err
	}

	// 克隆请求，避免修改调用方持有的 Request。
	var outReq = req.Clone(req.Context())

	// 实际连接目标使用 Endpoint。
	outReq.URL.Scheme = endpointURL.Scheme
	outReq.URL.Host = endpointURL.Host

	// HTTP Host 默认使用 Endpoint。
	outReq.Host = ""

	resp, err := t.base.RoundTrip(outReq)
	if resp != nil {
		// 对外保留原始请求，避免暴露内部改写后的 URL。
		resp.Request = req
	}

	return resp, err
}

// CloseIdleConnections 关闭底层 Transport 的空闲连接。
func (t *Transport) CloseIdleConnections() {
	t.base.CloseIdleConnections()
}

// parseEndpointURL 解析并验证 Endpoint 地址。
func parseEndpointURL(addr string) (*url.URL, error) {
	if addr == "" {
		return nil, fmt.Errorf("ngx: invalid endpoint address")
	}

	u, err := url.Parse(addr)
	if err != nil {
		return nil, fmt.Errorf("ngx: invalid endpoint address %q: %w", addr, err)
	}

	u.Scheme = strings.ToLower(u.Scheme)

	if u.Scheme != "http" && u.Scheme != "https" {
		return nil, fmt.Errorf("ngx: unsupported endpoint scheme %q", u.Scheme)
	}

	if u.Host == "" || u.Hostname() == "" {
		return nil, fmt.Errorf("ngx: endpoint address %q has no host", addr)
	}

	if u.User != nil || u.Opaque != "" {
		return nil, fmt.Errorf("ngx: invalid endpoint address %q", addr)
	}

	if (u.Path != "" && u.Path != "/") ||
		u.RawQuery != "" ||
		u.ForceQuery ||
		u.Fragment != "" {
		return nil, fmt.Errorf("ngx: endpoint address %q must contain only scheme and host", addr)
	}

	return u, nil
}

// normalizeHost 规范化服务名称。
func normalizeHost(host string) string {
	host = strings.ToLower(host)
	return strings.TrimSuffix(host, ".")
}
