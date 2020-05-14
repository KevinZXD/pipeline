package pipeline_test

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"pipeline"
)

type URL struct {
	name  string
	url   string
	value string
}

func newURL(name, url string) *URL {
	return &URL{
		url:  url,
		name: name,
	}
}

func (e *URL) Run(ctx context.Context) {
	res := pipeline.UpstreamValue(ctx)
	if res != nil {
		fmt.Println("upstream value:", res)
	}
	// real work
	http.Get(e.url)
	e.value = e.url
}

func (e *URL) Values() map[string]interface{} {
	return pipeline.NodeValue(e.name, e.value)
}

func ExamplePipeline() {
	seriesURLs := []string{
		"http://httpbin.org/delay/1",
		"http://httpbin.org/delay/2",
		"http://httpbin.org/delay/3",
	}

	parallelURLs := []string{
		"http://httpbin.org/delay/4",
		"http://httpbin.org/delay/6",
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
	fmt.Println(parallel.Values())

	// Output:
	// upstream value: map[series_0:http://httpbin.org/delay/1]
	// upstream value: map[series_1:http://httpbin.org/delay/2]
	// map[parallel_0:http://httpbin.org/delay/4 parallel_1:http://httpbin.org/delay/6 parallel_2:http://httpbin.org/delay/7 series_0:http://httpbin.org/delay/1 series_1:http://httpbin.org/delay/2 series_2:http://httpbin.org/delay/3]
}
