package lb

import "context"

// Resolver 根据服务名称解析可用的服务端点。
//
// 实现必须支持并发调用。
// Resolve 返回的 endpoints 在调用方完成本次读取前不得被修改。
type Resolver interface {
	Resolve(ctx context.Context, service string) ([]Endpoint, error)
}
