# monitor
based on prometheus

# how to use
1. ./prometheus --web.listen-address="0.0.0.0:9999"(prometheus.yaml添加localhost:8686这个job)
1. go run examples/main.go
2. ./wrk -t 2 -c 50 -d 1m --latency http://localhost:8686/ping
3. histogram_quantile(0.90,sum(rate(redis_storage_exec_cost_bucket{status!="404"}[1m])) by (handler,le))
# grafana注意点
在panel里编辑prom语句，每次必现"1:19: parse error: unexpected character: '\ufeff'"<br>
解决方法：在prometheus web后台execute这个prom语句ok后，粘贴到grafana中