// Samples used in a small go tutorial
//
// Copyright (C) 2017 framp at linux-tips-and-tricks dot de
//
// Samples for go syntax - simple REST client which calculates the geolocation, sunrise and sunset time of a location
//
// Sample locations:
// go run geoLocationSunriseSunset.go -location "Stuttgart, Königstrasse, 1"
// Output:
// --- REST API call: Retrieving geolocation for Stuttgart, Königstrasse, 1
// --- URL: http://maps.googleapis.com/maps/api/geocode/json?address=Stuttgart%2C+K%C3%B6nigstrasse%2C+1
// Longitude: 9.182134
// Latitude : 48.782154
// --- REST API call: Retrieving UTC sunrise and sunset for 48.782154, 9.182134
// --- URL: http://api.sunrise-sunset.org/json?date=today&formatted=0&lat=48.782154&lng=9.182134
// sunrise: 2017-03-24T05:15:58+00:00 UTC
// sunset : 2017-03-24T17:42:42+00:00 UTC
// --- Converting UTC into local time
// sunrise: 2017-03-24 06:15:58 +0100 CET
// sunset : 2017-03-24 18:42:42 +0100 CET
//
// Other locations
// go run geoLocationSunriseSunset.go -location "Canberra, Batman Street, 1"
// go run geoLocationSunriseSunset.go -location "Berlin"
// go run geoLocationSunriseSunset.go -location "Paris"
// go run geoLocationSunriseSunset.go -location "New York"

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"
)

var debug bool // debug flag

// REST api url
const googleLocationURL = "http://maps.googleapis.com/maps/api/geocode/json"

// GoogleLocationResponse - JSON response returned by google location REST API
type GoogleLocationResponse struct {
	Results []struct {
		Geometry struct {
			Location struct {
				Lattitude float32 `json:"lat"`
				Logitude  float32 `json:"lng"`
			} `json:"location"`
		} `json:"geometry"`
	} `json:"results"`
	Status string `json:"status"`
}

// REST api url
const sunriseSunsetOrgURL = "http://api.sunrise-sunset.org/json"

// SunriseSunsetOrgResponse - JSON response returned by sunrise-sunset.org REST API
type SunriseSunsetOrgResponse struct {
	Results struct {
		Sunrise string `json:"sunrise"`
		Sunset  string `json:"sunset"`
	} `json:"results"`
	Status string `json:"status"`
}

// helperfunction for errors
func abortIfError(err error) {
	if err != nil {
		log.Fatalf("%+v\n", err)
		os.Exit(42)
	}
}

// helper to create encoded query parms for url
func getEncodedParms(values map[string]string) string {
	vals := url.Values{}
	for k, v := range values { // loop for map elements, return key, value
		vals.Add(k, v) // append value to key if key was already used
	}
	return vals.Encode()
}

// helper to calculate local time
func localTime(t time.Time) time.Time {
	loc, err := time.LoadLocation("Local")
	abortIfError(err)
	return t.In(loc)
}

// helper to execute get REST calls
func retrievePage(url string) *[]byte {
	fmt.Printf("--- URL: %v\n", url)
	// establish connection
	resp, err := http.Get(url)
	defer resp.Body.Close() // close connection at end of func
	abortIfError(err)

	// retrieve json payload
	body, err := ioutil.ReadAll(resp.Body)
	abortIfError(err)
	if debug {
		fmt.Printf("*** Retrieved ***\n%v\n", string(body))
	}
	return &body
}

func main() {

	// parse command arguments, set default parms and help text for -h flag
	location := flag.String("location", "Stuttgart, Bahnhofstraße, 1", "Location to query sunrise and sunset")
	flag.BoolVar(&debug, "debug", false, "enable debug output")
	flag.Parse()

	// retrieve geolocation information
	fmt.Printf("--- REST API call: Retrieving geolocation for %s\n", *location)

	// set url encoded REST call queries
	// url?address=<location>
	parms := getEncodedParms(map[string]string{
		"address": *location,
	})
	body := retrievePage(googleLocationURL + "?" + parms)

	// unmarshall the json response into go struct
	response := new(GoogleLocationResponse)
	err := json.Unmarshal(*body, response)
	abortIfError(err)

	if response.Status != "OK" {
		fmt.Printf("Failed to retrieve geolocation. %s", response.Status)
		os.Exit(42)
	}
	longitude := response.Results[0].Geometry.Location.Logitude
	latitude := response.Results[0].Geometry.Location.Lattitude

	fmt.Printf("Longitude: %v\nLatitude : %v\n", longitude, latitude)

	// retrieve sunset and sunrise time in UTC time

	fmt.Printf("--- REST API call: Retrieving UTC sunrise and sunset for %f, %f\n", latitude, longitude)

	// set url encoded REST call queries
	// url?lat=<lat>?lng=<lng>?date=today?formatted=0
	parms = getEncodedParms(map[string]string{
		"lat":       fmt.Sprintf("%v", latitude),
		"lng":       fmt.Sprintf("%v", longitude),
		"date":      "today",
		"formatted": "0",
	})
	body = retrievePage(sunriseSunsetOrgURL + "?" + parms)

	// unmarshall the json response into go struct
	response2 := new(SunriseSunsetOrgResponse)
	err = json.Unmarshal(*body, response2)
	abortIfError(err)

	// 2017-03-24T04:57:33+00:00
	// ISO 8601
	sunrise := response2.Results.Sunrise
	sunset := response2.Results.Sunset

	fmt.Printf("sunrise: %v UTC\nsunset : %v UTC\n", sunrise, sunset)

	// convert sunset and sunrise time into local time
	// 2017-03-24 18:28:47.005227857 +0000 UTC

	fmt.Printf("--- Converting UTC into local time\n")

	sr, err := time.Parse(time.RFC3339, sunrise)
	abortIfError(err)
	ss, err := time.Parse(time.RFC3339, sunset)
	abortIfError(err)
	fmt.Printf("sunrise: %v\nsunset : %v\n", localTime(sr), localTime(ss))
}
