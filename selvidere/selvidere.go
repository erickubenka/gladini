package main

import (
	"encoding/json"
	"flag"
	"strconv"
	"io/ioutil"
	"net/http"
	//"sync"
	"time"

	"github.com/prometheus/common/log"
)

var (
	interval   = flag.String("interval", "5", "Time interval to scrape in seconds")
	gridUri     = flag.String("grid-uri", "http://localhost:30020", "URI on which to scrape Selenium Grid.")
)

func scrape() {

	body, err := fetch()
	if err != nil {
		log.Errorf("Can't scrape Selenium Grid: %v", err)
		return
	}

	// get complete response as json raw map
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
		log.Infoln("Grid ready.")
	}

	// get {hubResponseMap.value.nodes[]}
	var nodesMap []map[string]json.RawMessage
	err = json.Unmarshal(valueMap["nodes"], &nodesMap)
	log.Infoln("Registered nodes: ", len(nodesMap))

	// iterate every node response
	var tmpNodeList []string

	for _, nodeMapEntry := range nodesMap {

		var warning string
		var id string
		json.Unmarshal(nodeMapEntry["id"], &id)

		if(contains(tmpNodeList, id)) {
			log.Warnln("Node exists twice, removing: ", id)
			delete(id)
			continue;
		}

		tmpNodeList = append(tmpNodeList, id)

		if err := json.Unmarshal(nodeMapEntry["warning"], &warning); err == nil {
			log.Warnln("Node has warnings present, removing: ", id)
			delete(id)
		} else {
			log.Infoln("Node works properly:", id)
		}
	}
}

func contains(arr []string, str string) bool {
	for _, a := range arr {
	   if a == str {
		  return true
	   }
	}
	return false
 }

func delete(id string) {

	var url string = *gridUri + "/se/grid/distributor/node/" + id

	client := http.Client{
		Timeout: 3 * time.Second,
	}

	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		log.Fatal(err)
		return
	}

	req.Header.Set("X-REGISTRATION-SECRET", "")

	resp, err := client.Do(req)
    if err != nil {
        log.Fatal(err)
        return
    }
    defer resp.Body.Close()
}

func fetch() (output []byte, err error) {

	client := http.Client{
		Timeout: 3 * time.Second,
	}

	response, err := client.Get(*gridUri + "/status")
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

	log.Infoln("Starting selvidere")

	for {
		log.Infoln("Scraping...")
		scrape()
		timespan, _ := strconv.Atoi(*interval);
		time.Sleep(time.Duration(timespan) * time.Second)
	}
}
