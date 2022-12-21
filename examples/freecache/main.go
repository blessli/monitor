package main

import (
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"time"

	"github.com/blessli/monitor/pkg"
	"github.com/coocood/freecache"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	engine           *gin.Engine
	CustomCollector  *prometheus.HistogramVec
	CustomCollector2 *prometheus.GaugeVec
	cacheInstance    *freecache.Cache
	cacheValue       = getTestCacheValue()
	cacheSize        = 0
)

func getRandNum() string {
	rand.Seed(time.Now().Unix())
	aaaKey := "aaa:%d:buy:%d:%d:%d:%d"
	key := fmt.Sprintf(aaaKey, rand.Int(), rand.Int(), rand.Int(), rand.Int(), rand.Int())
	return key
}

func getTestCacheValue() []byte {
	ss := ""
	for i := 0; i < 1e2; i++ {
		ss += uuid.New().String()
	}
	cacheValue := []byte(ss)
	cacheSize = len(cacheValue)
	log.Println("cacheValue size: ", cacheSize)
	return cacheValue
}
func init() {
	cacheInstance = freecache.NewCache(5 * 1024 * 1024)
	engine = gin.New()
	engine.GET("/metrics", func(c *gin.Context) {
		h := promhttp.Handler()
		h.ServeHTTP(c.Writer, c.Request)
	})
	NewColleactor()
	NewColleactor2()
	pkg.InitPuller(CustomCollector,CustomCollector2)
}
func NewColleactor() {
	dBuckets := []float64{5, 10, 25, 50, 100, 250, 500, 1000, 2500, 5000, 10000}
	CustomCollector = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "freecache_hit_rate",
			Help:    "freecache命中率",
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

func NewColleactor2() {
	CustomCollector2 = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:      "freecache_hit_rate2",
			Namespace: "freecache",
			Help:      "This represent the number of items in cache",
		},
		[]string{
			"hostname",
			"method",
			"status",
			"rmsp",
		},
	)
}

func PingHandler(c *gin.Context) {
	key := getRandNum()
	v, err := cacheInstance.Get([]byte(key))
	if err != nil && err.Error() != freecache.ErrNotFound.Error() {
		log.Println("freecache get error: ", err)
		c.JSON(200, gin.H{
			"message": "error",
			"status":  "50000",
		})
		return
	}
	if len(v) == 0 {
		if err := cacheInstance.Set([]byte(key), cacheValue, 2); err != nil {
			log.Println("ser error: ", err)
		}
	}
	c.JSON(200, gin.H{
		"message": "ok",
		"status":  "0",
	})
}

func main() {
	engine.GET("/ping", PingHandler)
	go func() {
		ticker := time.NewTicker(time.Second)
		for {
			select {
			case <-ticker.C:
				report()
			}
		}
	}()
	engine.Run(":8686")
}

func report() {
	hitRate := cacheInstance.HitRate()
	// now := time.Now().Unix()
	aat := cacheInstance.AverageAccessTime()
	log.Printf("hit rate is %v, evacuates %v, entries %v, average cost %v, expire count %v, cachevalue size %v\n",
		cacheInstance.HitRate(), cacheInstance.EvacuateCount(), cacheInstance.EntryCount(), aat/(int64(time.Millisecond)/int64(time.Nanosecond)), cacheInstance.ExpiredCount(), cacheSize)
	HitRateReport("local", "get", "unknown", 200, hitRate*100)
	HitRateReport2("local", "get", "unknown", 200, hitRate*100)
}
func HitRateReport(hostname, method, rmsp string, code int, hitRate float64) {
	//	统计直方图histogram指标
	CustomCollector.With(prometheus.Labels{
		"hostname": hostname,
		"method":   method,
		"status":   strconv.Itoa(code),
		"rmsp":     rmsp,
	}).Observe(hitRate)
}

func HitRateReport2(hostname, method, rmsp string, code int, hitRate float64) {
	//	统计直方图histogram指标
	CustomCollector2.With(prometheus.Labels{
		"hostname": hostname,
		"method":   method,
		"status":   strconv.Itoa(code),
		"rmsp":     rmsp,
	}).Set(hitRate)
}
