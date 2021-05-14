package main

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"strconv"
)

//Define a struct for you collector that contains pointers
type Collector struct {
	webMetric *prometheus.Desc
}

type ServiceInfo struct {
	servicename string
	port        string
	api         string
}

var serviceinfos = []ServiceInfo{
	{
		servicename: "www.baidu.com",
		port:        "80",
		api:         "/",
	},
	{
		servicename: "www.wework.cn",
		port:        "80",
		api:         "/",
	},
}

//initializes every descriptor and returns a pointer to the collector
func newWebCollector() *Collector {
	return &Collector{
		webMetric: prometheus.NewDesc(
			"web_code_metric",
			"Shows web code metric",
			[]string{"servicename"},
			prometheus.Labels{},
		),
	}
}

//Describe function.
func (collector *Collector) Describe(ch chan<- *prometheus.Desc) {
	//Update this section with the each metric you create for a given collector
	ch <- collector.webMetric
}

//collect function
func (collector *Collector) Collect(ch chan<- prometheus.Metric) {
	metricValues := WebCheckGet(serviceinfos)
	for k, v := range metricValues {
		//Write latest value for each metric in the prometheus metric channel.
		//Note that you can pass CounterValue, GaugeValue, or UntypedValue types here.
		strv := strconv.Itoa(v)
		log.Info("调用地址：" + k + "，响应状态码：" + strv)
		ch <- prometheus.MustNewConstMetric(collector.webMetric, prometheus.GaugeValue, float64(v), k)
	}
}

//web check code function
func WebCheckGet(infos []ServiceInfo) map[string]int {
	webCode := map[string]int{}
	for _, v := range infos {
		resp, err := http.Get("http://" + v.servicename + ":" + v.port + v.api)
		if err != nil {
			log.Error(err)
			webCode[v.servicename] = 0
		} else {
			webCode[v.servicename] = resp.StatusCode
		}
	}
	return webCode
}

func init() {
	log.SetFormatter(&log.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})
	log.SetOutput(os.Stdout)
}

func main() {
	reg := prometheus.NewPedanticRegistry()
	web := newWebCollector()
	reg.MustRegister(web)

	//any metrics on the /metrics endpoint.
	http.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{}))
	log.Info("Beginning to serve on port :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
