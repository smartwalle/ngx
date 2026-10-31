package lb

import (
	"context"
	"errors"
	"sync"
)

var ErrNoEndpoints = errors.New("ngx: no endpoints available")

// Balancer 从候选端点中选择一个服务端点。
//
// 实现必须支持并发调用。
// Pick 不得修改传入的 endpoints 切片。
// 返回值必须是 endpoints 中的一个端点。
type Balancer interface {
	Pick(ctx context.Context, service string, endpoints []Endpoint) (Endpoint, error)
}

// RoundRobinBalancer 使用轮询策略选择服务端点。
//
// RoundRobinBalancer 的零值可以直接使用。
type RoundRobinBalancer struct {
	mu   sync.Mutex
	next map[string]uint64
}

// NewRoundRobinBalancer 创建一个轮询负载均衡器。
func NewRoundRobinBalancer() *RoundRobinBalancer {
	return &RoundRobinBalancer{
		next: make(map[string]uint64),
	}
}

// Pick 按轮询顺序选择服务端点。
func (b *RoundRobinBalancer) Pick(ctx context.Context, service string, endpoints []Endpoint) (Endpoint, error) {
	if len(endpoints) == 0 {
		return Endpoint{}, ErrNoEndpoints
	}

	b.mu.Lock()
	if b.next == nil {
		b.next = make(map[string]uint64)
	}
	var index = b.next[service]
	b.next[service] = index + 1

	b.mu.Unlock()

	return endpoints[index%uint64(len(endpoints))], nil
}
