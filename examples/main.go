package main

import (
	"context"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/blessli/monitor/pkg"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	engine *gin.Engine
	CustomCollector *prometheus.HistogramVec
)
func init() {
	engine = gin.New()
	engine.GET("/metrics", func(c *gin.Context) {
		h := promhttp.Handler()
		h.ServeHTTP(c.Writer, c.Request)
	})
	NewColleactor()
	pkg.InitPuller(CustomCollector)
}
func NewColleactor() {
	dBuckets := []float64{5, 10, 25, 50, 100, 250, 500, 1000, 2500, 5000, 10000}
	CustomCollector = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "redis_storage_exec_cost",
			Help:    "redis存储执行耗时",
			Buckets: dBuckets,
		},
		[]string{
			"hostname",
			"method",
			"status",
			"rmsp",
		},
	)
}

func main() {
	engine.GET("/ping", func(c *gin.Context) {
		begin := time.Now()
		ctx := context.TODO()
		client := &http.Client{}
		// mock third call
		req, _ := http.NewRequestWithContext(ctx, "GET", "https://www.baidu.com/s?wd=%E4%B8%89%E4%B8%AA%E8%AF%8D%E8%AF%BB%E6%87%82%E2%80%9C%E4%B8%89%E7%8E%AF%E5%B3%B0%E4%BC%9A%E2%80%9D&sa=fyb_n_homepage&rsv_dl=fyb_n_homepage&from=super&cl=3&tn=baidutop10&fr=top1000&rsv_idx=2&hisfilter=1", strings.NewReader(""))
		resp, err := client.Do(req)
		if err!=nil {
			StorageExecReport("local", "get", "unknown", 50000, time.Since(begin))
			c.JSON(200, gin.H{
				"message": "error",
				"status":  "50000",
			})
			return
		}
		// test
		if resp != nil {
			StorageExecReport("local", "get", "unknown", resp.StatusCode, time.Since(begin))
		}
		c.JSON(200, gin.H{
			"message": "pong",
			"status":  resp.StatusCode,
		})
	})
	engine.Run(":8686")
}
func StorageExecReport(hostname, method, rmsp string, code int, duration time.Duration) {
	//	统计直方图histogram指标
	CustomCollector.With(prometheus.Labels{
		"hostname": hostname,
		"method":   method,
		"status":   strconv.Itoa(code),
		"rmsp":     rmsp,
	}).Observe(float64(duration.Milliseconds()))
}