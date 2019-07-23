package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type Info struct {
	Region string `json:"region_name"`
	City   string `json:"city"`
}

// Echo headers and info
type Echo struct {
	Headers http.Header `json:"headers"`
	IP      string      `json:"ip"`
	Info    Info        `json:"info"`
}

var apiURL = "api.ipstack.com"
var APIKEY = "xxx"

func startHTTPServer() {

}

// GetUserLocation get a users region and city
func GetUserLocation(remote string) Info {
	ip := strings.Split(remote, ":")
	info := Info{}
	url := fmt.Sprintf("http://%s/%s?access_key=%s&format=1", apiURL, ip, APIKEY)
	resp, err := http.Get(url)
	if err != nil {
		log.Printf("error getting info from %s", url)
	}
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("error reading info from %v", resp)
	}
	err = json.Unmarshal(respBody, &info)
	return info
}
