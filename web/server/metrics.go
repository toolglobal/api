package server

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"strconv"
	"time"

	"github.com/huzhongqing/ginprom"
)

var (
	PromConfig = ginprom.Config{
		FixedPath: []string{
			"/v1/actions",
			"/v1/transactions",
			"/v1/ledgers",

			"/v2/contract/transactions",
			"/v2/transactions",
		},

		ParamsPath: map[string]int{
			"/v2/accounts/:address/transactions": 3,
			"/v1/accounts/:address":              3,
			"/v1/transactions/:address":          3,
		},
	}
)

func NewMetrics(e *gin.Engine) *ginprom.GinPrometheus {
	p := ginprom.New(e, &PromConfig)
	p.DefaultRegister("api")
	// 开启 ／metrics 接口
	p.Use(e)

	// 开启采集
	//e.Use(p.DefaultHandlerFunc(), gin.Recovery())

	iMetric := ginprom.Metric{
		ID:          "reqElapsed",
		Name:        "http_request_elapsed_second",
		Description: "http api request elapsed",
		Type:        "histogram_vec",
		Args:        []string{"method", "path"},
	}

	if err := p.AddMetrics(&iMetric, "api"); err != nil {
		panic("ginprom add metrics " + err.Error())
	}

	// 设置自定义的采集
	// 加入接口直方图
	handlerF := func(gp *ginprom.GinPrometheus) gin.HandlerFunc {
		return func(c *gin.Context) {
			path := c.Request.URL.Path
			path, ok := gp.HitPath(path)
			if !ok {
				c.Next()
				return
			}

			start := time.Now()
			requestSize := gp.ReqSize(c.Request)

			c.Next()

			status := strconv.Itoa(c.Writer.Status())
			elapsed := float64(time.Since(start)) / float64(time.Second)
			responseSize := float64(c.Writer.Size())

			reqCnt, ok := gp.MetricsMap.Load("reqCnt")
			if ok {
				reqCnt.(ginprom.Metric).Collector.(*prometheus.CounterVec).WithLabelValues(status, c.Request.Method, c.HandlerName(), c.Request.Host, path).Inc()
			}

			reqSz, ok := gp.MetricsMap.Load("reqSz")
			if ok {
				reqSz.(ginprom.Metric).Collector.(prometheus.Summary).Observe(float64(requestSize))
			}

			resSz, ok := gp.MetricsMap.Load("resSz")
			if ok {
				resSz.(ginprom.Metric).Collector.(prometheus.Summary).Observe(float64(responseSize))
			}

			reqDur, ok := gp.MetricsMap.Load("reqDur")
			if ok {
				reqDur.(ginprom.Metric).Collector.(prometheus.Summary).Observe(float64(elapsed))
			}
			reqElapsed, ok := gp.MetricsMap.Load("reqElapsed")
			if ok {
				reqElapsed.(ginprom.Metric).Collector.(*prometheus.HistogramVec).WithLabelValues(c.Request.Method, path).Observe(float64(elapsed))
			}
		}
	}

	e.Use(handlerF(p), gin.Recovery())

	return p
}
