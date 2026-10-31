package lb

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

// Endpoint 表示一个服务端点。
type Endpoint struct {
	Name string
	Addr string
}

// Transport 实现基于服务发现的 HTTP 请求级负载均衡。
type Transport struct {
	base *http.Transport

	resolver Resolver
	balancer Balancer
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

// UseBalancer 设置负载均衡器。
func (t *Transport) UseBalancer(balancer Balancer) {
	t.balancer = balancer
}

// UseResolver 设置服务发现。
func (t *Transport) UseResolver(resolver Resolver) {
	t.resolver = resolver
}

// RoundTrip 实现请求级负载均衡。
func (t *Transport) RoundTrip(req *http.Request) (*http.Response, error) {
	if t.balancer == nil || t.resolver == nil {
		return t.base.RoundTrip(req)
	}

	var reqScheme = strings.ToLower(req.URL.Scheme)
	if reqScheme != "http" && reqScheme != "https" {
		return t.base.RoundTrip(req)
	}

	var service = normalizeHost(req.URL.Hostname())

	endpoints, err := t.resolver.Resolve(req.Context(), service)
	if err != nil {
		return nil, err
	}
	if len(endpoints) == 0 {
		return t.base.RoundTrip(req)
	}

	selected, err := t.balancer.Pick(req.Context(), service, endpoints)
	if err != nil {
		return nil, err
	}
	if selected.Addr == "" {
		return nil, errors.New("ngx: balancer selected endpoint with empty address")
	}

	endpointURL, err := parseEndpointAddr(selected.Addr)
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

// parseEndpointAddr 解析并验证 Endpoint 地址。
func parseEndpointAddr(addr string) (*url.URL, error) {
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
