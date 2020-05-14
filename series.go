package pipeline

import (
	"context"
)

// series 串行管理器
//
// 每个任务实现Worker接口，装载到SeriesGroup中
// 调用Work函数时，将会顺序执行每个worker
type series struct {
	nodes []Node
	value map[string]interface{}
}

// NewSeries 创建串行实例。可使用Append装载多个Node，使用Run串行执行
func NewSeries() *series {
	return &series{
		nodes: make([]Node, 0),
		value: make(map[string]interface{}),
	}
}

// Append 添加Worker
func (sg *series) Append(ws ...Node) {
	sg.nodes = append(sg.nodes, ws...)
}

// Run 串行执行每个Node
//
// 如果未装载node，将什么都不执行
func (sg *series) Run(ctx context.Context) {
	if len(sg.nodes) == 0 {
		return
	}
	for _, w := range sg.nodes {
		done := make(chan struct{})
		go func() {
			w.Run(ctx)
			close(done)
		}()
		// 使用select来监听context是否已经done
		select {
		case <-ctx.Done():
			return
		case <-done:
		}
		// 将串行上游的值添加到context并传递给下游
		values := w.Values()
		ctx = context.WithValue(ctx, upstreamContextKey, values)
		if values == nil {
			continue
		}
		for name, val := range values {
			sg.value[name] = val
		}
	}
}

// Values 获取串行node的所有结果
func (sg *series) Values() map[string]interface{} {
	return sg.value
}
