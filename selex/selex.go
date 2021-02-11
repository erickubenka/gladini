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
URI, message                         																		string
	mutex                                																	sync.RWMutex
	up, ready, registeredNodes																				prometheus.Gauge
	chromeNodes																								prometheus.Gauge
	chromeSessionsInUse, chromeSessionsTotalAvailable, chromeSessionsFree, chromeSessionsInUsePercent 		prometheus.Gauge
	firefoxNodes																							prometheus.Gauge
	firefoxSessionsInUse, firefoxSessionsTotalAvailable, firefoxSessionsFree, firefoxSessionsInUsePercent 	prometheus.Gauge
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
	StereoTypes []StereoType `json:"stereotype"`
	Sessions []Session `json:"sessions"`
	
	// only present on fault
	Warning string `json:warning`
}

type StereoType struct {

	// Capabilities Capability `json:"capabilities"`
	// Count float64 `json:"count"`
	BrowserName string `json:"browserName"`
}

type Slot struct {
	LastStarted string `json:"lastStarted"`
	StereoType StereoType `json:"stereotype"`
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
		chromeSessionsInUsePercent: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: nameSpace,
			Subsystem: subSystem,
			Name:      "chromeSessionsInUsePercent",
			Help:      "Percentage of used chrome sessions",
		}),
		firefoxNodes: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: nameSpace,
			Subsystem: subSystem,
			Name:      "firefoxNodes",
			Help:      "total number of registered firefox nodes",
		}),
		firefoxSessionsInUse: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: nameSpace,
			Subsystem: subSystem,
			Name:      "firefoxSessionsInUse",
			Help:      "total number of firefox sessions in use",
		}),
		firefoxSessionsTotalAvailable: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: nameSpace,
			Subsystem: subSystem,
			Name:      "firefoxSessionsTotalAvailable",
			Help:      "total number of firefox sessions available regardless of use state",
		}),
		firefoxSessionsFree: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: nameSpace,
			Subsystem: subSystem,
			Name:      "firefoxSessionsFree",
			Help:      "total number of firefox sessions free",
		}),
		firefoxSessionsInUsePercent: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: nameSpace,
			Subsystem: subSystem,
			Name:      "firefoxSessionsInUsePercent",
			Help:      "Percentage of used firefox sessions",
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
	e.chromeSessionsInUsePercent.Describe(ch)

	e.firefoxNodes.Describe(ch)
	e.firefoxSessionsInUse.Describe(ch)
	e.firefoxSessionsTotalAvailable.Describe(ch)
	e.firefoxSessionsFree.Describe(ch)
	e.firefoxSessionsInUsePercent.Describe(ch)
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
	ch <- e.chromeSessionsInUsePercent

	ch <- e.firefoxNodes
	ch <- e.firefoxSessionsInUse
	ch <- e.firefoxSessionsTotalAvailable
	ch <- e.firefoxSessionsFree
	ch <- e.firefoxSessionsInUsePercent

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
	e.chromeSessionsInUsePercent.Set(0)

	e.firefoxNodes.Set(0)
	e.firefoxSessionsInUse.Set(0)
	e.firefoxSessionsTotalAvailable.Set(0)
	e.firefoxSessionsFree.Set(0)
	e.firefoxSessionsInUsePercent.Set(0)

	body, err := e.fetch()
	if err != nil {
		e.up.Set(0)
		log.Errorf("Can't scrape Selenium Grid: %v", err)
		return
	}

	e.up.Set(1)

	// get complete respone as json raw map
	var hubResponseMap map[string]json.RawMessage
	if err := json.Unmarshal(body, &hubResponseMap); err != nil {
		log.Errorf("Can't decode Selenium Grid response: %v", err)
		return
	}

	// get {hubResponseMap.value}
	var valueMap map[string]json.RawMessage
	err = json.Unmarshal(hubResponseMap["value"], &valueMap)
	
	// Set ready state of grid...
	var readyState bool
	err = json.Unmarshal(valueMap["ready"], &readyState)

	if readyState {
		e.ready.Set(1)
	}

	// get {hubResponseMap.value.nodes[]}
	var nodesMap []map[string]json.RawMessage
	err = json.Unmarshal(valueMap["nodes"], &nodesMap)

	// set registered node count.
	e.registeredNodes.Set(float64(len(nodesMap)))
	
	// set registered node count for chrome...
	var chromeNodeCount float64 = 0
	var chromeSessionCount float64 = 0
	var chromeSessionMax float64 = 0

	// set registered node count for firefox
	var firefoxNodeCount float64 = 0
	var firefoxSessionCount float64 = 0
	var firefoxSessionMax float64 = 0

	// iterate every node response
	for _, nodeMapEntry := range nodesMap {

		var warning string
		var id string
		json.Unmarshal(nodeMapEntry["id"], &id)

		if err := json.Unmarshal(nodeMapEntry["warning"], &warning); err == nil {
			log.Warnln("Node has warnings present:", id)
		} else {
			log.Infoln("Node works properly:", id)
			var slots []Slot
			json.Unmarshal(nodeMapEntry["slots"], &slots)

			for _, slot := range slots {

				// counting chrome nodes here...
				if slot.StereoType.BrowserName == "chrome" {
					// increase chrome node count here...
					chromeNodeCount++
					
					// try to get sessions count here...
					var sessionsMap []map[string]json.RawMessage
					json.Unmarshal(nodeMapEntry["sessions"], &sessionsMap)
					chromeSessionCount = chromeSessionCount + float64(len(sessionsMap))
					
					var maxSessions float64
					json.Unmarshal(nodeMapEntry["maxSessions"], &maxSessions)
					chromeSessionMax = chromeSessionMax + float64(maxSessions)
				}

				// counting firefox nodes here...
				if slot.StereoType.BrowserName == "firefox" {
					// increase firefox node count here...
					firefoxNodeCount++
					
					// try to get sessions count here...
					var sessionsMap []map[string]json.RawMessage
					json.Unmarshal(nodeMapEntry["sessions"], &sessionsMap)
					firefoxSessionCount = firefoxSessionCount + float64(len(sessionsMap))
					
					var maxSessions float64
					json.Unmarshal(nodeMapEntry["maxSessions"], &maxSessions)
					firefoxSessionMax = firefoxSessionMax + float64(maxSessions)
				}

				// counting other nodes here...
				// todo
			}
		}
	}

	e.chromeNodes.Set(chromeNodeCount)
	e.chromeSessionsInUse.Set(chromeSessionCount)
	e.chromeSessionsTotalAvailable.Set(chromeSessionMax)
	e.chromeSessionsFree.Set(chromeSessionMax-chromeSessionCount)
	e.chromeSessionsInUsePercent.Set((chromeSessionCount/chromeSessionMax)*100)

	e.firefoxNodes.Set(firefoxNodeCount)
	e.firefoxSessionsInUse.Set(firefoxSessionCount)
	e.firefoxSessionsTotalAvailable.Set(firefoxSessionMax)
	e.firefoxSessionsFree.Set(firefoxSessionMax-firefoxSessionCount)
	e.firefoxSessionsInUsePercent.Set((firefoxSessionCount/firefoxSessionMax)*100)
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
