package pipeline

import "context"

// Node 任务接口
type Node interface {
	// Run 执行任务主逻辑
	Run(context.Context)
	// NodeValue 返回任务执行的结果
	Values() map[string]interface{}
}

const (
	upstreamContextKey = "_upstreamContextKey" // 串行context传递上游结果的key
)

// UpstreamValue 返回上游的值，如果不存在返回nil
//
// 该函数用于Series串行内的Node
func UpstreamValue(ctx context.Context) map[string]interface{} {
	val := ctx.Value(upstreamContextKey)
	if val != nil {
		return val.(map[string]interface{})
	}
	return nil
}

// NodeValue 构造原子Node的值 省的每次都要实现一坨代码
func NodeValue(name string, value interface{}) map[string]interface{} {
	m := make(map[string]interface{})
	m[name] = value
	return m
}
