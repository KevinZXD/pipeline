package pipeline

import (
	"context"
	"sync"
)

// parallel 并行管理器
//
// 每个任务实现Node接口，装载到Parallel中
// 然后使用WaitGroup来管理任务并发，Wait等待结束后，可通过Value函数获取结果
type parallel struct {
	nodes []Node
	wg    *sync.WaitGroup
	value map[string]interface{}
}

// NewParallel 创建并行实例，可装载多个Node，使用Run并发执行
func NewParallel() *parallel {
	return &parallel{
		nodes: make([]Node, 0),
		wg:    &sync.WaitGroup{},
		value: make(map[string]interface{}),
	}
}

// Append 装载多个Node
func (cg *parallel) Append(ws ...Node) {
	cg.nodes = append(cg.nodes, ws...)
}

// Run 并行执行每个Node
//
// 如果没有装载Node，将什么都不执行
func (cg *parallel) Run(ctx context.Context) {
	if len(cg.nodes) == 0 {
		return
	}
	for _, w := range cg.nodes {
		cg.wg.Add(1)
		go func(w Node) {
			w.Run(ctx)
			cg.wg.Done()
		}(w)
	}
	// 判断context是否已经done
	done := make(chan struct{})
	go func() {
		cg.wg.Wait()
		close(done)
	}()
	// 如果context已经done，及时退出，以免一直卡在前端
	select {
	case <-done:
	case <-ctx.Done():
		return
	}
	for _, w := range cg.nodes {
		values := w.Values()
		if values == nil {
			continue
		}
		for name, val := range values {
			cg.value[name] = val
		}
	}
}

// Values 获取并行Node的所有结果
func (cg *parallel) Values() map[string]interface{} {
	return cg.value
}
