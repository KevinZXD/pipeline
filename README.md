# Pipeline 流程管理器

Pipeline 用于管理任务，将多个任务的异步操作抽象成同步操作，结合异步的性能和同步的优雅，简单并灵活地运行多个任务，并获取任务结果。

该库借鉴了[errgroup](https://godoc.org/golang.org/x/sync/errgroup)库，降低难度，开箱即用。

每个任务需实现`pipeline.Node`接口，然后创建并行组`Parallel`或串行组`Series`示例，将Node加入到组中。最后调用Run开始执行。

实现如图所示的流程图（其中NodeX表示请求需要花费X秒，如Node5表示请求耗时5s）

<div  align="center">
<img src="http://git.intra.weibo.com/adx/pipeline/uploads/1be9ae6243015f33a177221972fec7d1/%E6%B5%81%E7%A8%8B%E5%9B%BE.png" width="600" alt="流程图"/>
</div>

示例代码：

```go
seriesURLs := []string{
    "http://httpbin.org/delay/1",
    "http://httpbin.org/delay/2",
    "http://httpbin.org/delay/3",
}

parallelURLs := []string{
    "http://httpbin.org/delay/4",
    "http://httpbin.org/delay/5",
    "http://httpbin.org/delay/7",
}

// 创建串行组
series := pipeline.NewSeries()
for i, url := range seriesURLs {
    // 加入到串行组中
    series.Append(newURL("series_"+strconv.Itoa(i), url))
}
// 创建并行组
parallel := pipeline.NewParallel()
for i, url := range parallelURLs {
    // 加入到并行组中
    parallel.Append(newURL("parallel_"+strconv.Itoa(i), url))
}
// 将串行组加入到并行组中
parallel.Append(series)

// 执行pipeline
parallel.Run(context.Background())

// 返回结果
fmt.Println("result:", parallel.Values())
```

运行结果：

```
upstream value: map[series_0:http://httpbin.org/delay/1]
upstream value: map[series_1:http://httpbin.org/delay/2]
result: map[parallel_0:http://httpbin.org/delay/4 parallel_1:http://httpbin.org/delay/5 parallel_2:http://httpbin.org/delay/7 series_0:http://httpbin.org/delay/1 series_1:http://httpbin.org/delay/2 series_2:http://httpbin.org/delay/3]
```

示例中`Node1-3`是串行同步执行，输出结果可看到，Node2/Node3分别获取到上个Node的运行结果.

并行耗时最多的是`Node7`，故程序总耗时7s多一点

完整示例见`example_test.go`