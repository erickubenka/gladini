package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"net/http"
	"sync"
	"time"
	
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/log"
)

const (
	nameSpace = "selenium_grid"
	subSystem = "hub"
)

var (
	listenAddress = flag.String("listen-address", ":8080", "Address on which to expose metrics.")
	metricsPath   = flag.String("telemetry-path", "/metrics", "Path under which to expose metrics.")
	gridUri     = flag.String("grid-uri", "http://localhost:30020", "URI on which to scrape Selenium Grid.")
)

type Exporter struct {
	URI, message                         												string
	mutex                                												sync.RWMutex
	up, ready, registeredNodes															prometheus.Gauge
	chromeNodes, chromeSessionsInUse, chromeSessionsTotalAvailable, chromeSessionsFree 	prometheus.Gauge
}

type HubResponse struct {
	Value HubResponseValue `json:"value"`
}

type HubResponseValue struct {
	Ready bool `json:"ready"`
	Message string `json:"message"`
	Nodes []Node `json:"nodes"`
}

type Node struct {
	Id string `json:"id"`
	MaxSessions float64 `json:"maxSessions"`
	StereoTypes []StereoType `json:"stereotypes"`
	Sessions []Session `json:"sessions"`
}

type StereoType struct {

	Capabilities Capability `json:"capabilities"`
	Count float64 `json:"count"`
}

type Capability struct {
	BrowserName string `json:"browserName"`
}

type Session struct {

}

func NewExporter(uri string) *Exporter {
	log.Infoln("Collecting data from:", uri)

	return &Exporter{
		URI: uri,
		up: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: nameSpace,
			Name:      "up",
			Help:      "was the last scrape of Selenium Grid successful.",
		}),
		ready: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: nameSpace,
			Subsystem: subSystem,
			Name:      "ready",
			Help:      "selenium grid ready state",
		}),
		registeredNodes: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: nameSpace,
			Subsystem: subSystem,
			Name:      "registeredNodes",
			Help:      "total number of registered nodes",
		}),
		chromeNodes: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: nameSpace,
			Subsystem: subSystem,
			Name:      "chromeNodes",
			Help:      "total number of registered chrome nodes",
		}),
		chromeSessionsInUse: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: nameSpace,
			Subsystem: subSystem,
			Name:      "chromeSessionsInUse",
			Help:      "total number of chrome sessions in use",
		}),
		chromeSessionsTotalAvailable: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: nameSpace,
			Subsystem: subSystem,
			Name:      "chromeSessionsTotalAvailable",
			Help:      "total number of chrome sessions available regardless of use state",
		}),
		chromeSessionsFree: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: nameSpace,
			Subsystem: subSystem,
			Name:      "chromeSessionsFree",
			Help:      "total number of chrome sessions free",
		}),
	}
}

func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	
	e.up.Describe(ch)
	e.ready.Describe(ch)
	e.registeredNodes.Describe(ch)

	e.chromeNodes.Describe(ch)
	e.chromeSessionsInUse.Describe(ch)
	e.chromeSessionsTotalAvailable.Describe(ch)
	e.chromeSessionsFree.Describe(ch)
}

func (e *Exporter) Collect(ch chan<- prometheus.Metric) {

	e.mutex.Lock()
	defer e.mutex.Unlock()

	e.scrape()

	ch <- e.up
	ch <- e.ready
	ch <- e.registeredNodes

	ch <- e.chromeNodes
	ch <- e.chromeSessionsInUse
	ch <- e.chromeSessionsTotalAvailable
	ch <- e.chromeSessionsFree

	return
}


func (e *Exporter) scrape() {

	e.ready.Set(0)
	e.up.Set(0)
	e.registeredNodes.Set(0)

	e.chromeNodes.Set(0)
	e.chromeSessionsInUse.Set(0)
	e.chromeSessionsTotalAvailable.Set(0)
	e.chromeSessionsFree.Set(0)

	body, err := e.fetch()
	if err != nil {
		e.up.Set(0)
		log.Errorf("Can't scrape Selenium Grid: %v", err)
		return
	}

	e.up.Set(1)

	log.Infoln(string(body))

	var hResponse HubResponse
	if err := json.Unmarshal(body, &hResponse); err != nil {

		log.Errorf("Can't decode Selenium Grid response: %v", err)
		return
	}

	// Set ready state of grid...
	if hResponse.Value.Ready {
		e.ready.Set(1)
	}
	
	// set registered node count.
	e.registeredNodes.Set(float64(len(hResponse.Value.Nodes)))
	
	// set registered node count for chrome...
	var chromeNodeCount float64 = 0
	var chromeSessionCount float64 = 0
	var chromeSessionMax float64 = 0
	for _, node := range hResponse.Value.Nodes {
		for _, stereoType := range node.StereoTypes {
			if stereoType.Capabilities.BrowserName == "chrome" {
				chromeNodeCount++
			}
		}

		chromeSessionCount = chromeSessionCount + float64(len(node.Sessions))
		chromeSessionMax = chromeSessionMax + float64(node.MaxSessions)
	}

	e.chromeNodes.Set(chromeNodeCount)
	e.chromeSessionsInUse.Set(chromeSessionCount)
	e.chromeSessionsTotalAvailable.Set(chromeSessionMax)
	e.chromeSessionsFree.Set(chromeSessionMax-chromeSessionCount)
}

func (e Exporter) fetch() (output []byte, err error) {

	client := http.Client{
		Timeout: 3 * time.Second,
	}

	response, err := client.Get(e.URI + "/status")
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	output, err = ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return
}


func main() {
	flag.Parse()

	log.Infoln("Starting selenium_grid_exporter")

	prometheus.MustRegister(NewExporter(*gridUri))
	prometheus.Unregister(prometheus.NewGoCollector())
	prometheus.Unregister(prometheus.NewProcessCollector(prometheus.ProcessCollectorOpts{}))

	http.Handle(*metricsPath, promhttp.Handler())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, *metricsPath, http.StatusMovedPermanently)
	})

	log.Fatal(http.ListenAndServe(*listenAddress, nil))
}
