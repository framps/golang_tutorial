// Samples used in a small go tutorial
//
// Copyright (C) 2017 framp at linux-tips-and-tricks dot de
//
// Samples for go syntax - simple REST client which calculates the sunrise and sunset time of a location

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

// Payload - JSON payload returned by REST API
type Payload struct {
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

// Payload2 - JSON payload returned by REST API
type Payload2 struct {
	Results struct {
		Sunrise string `json:"sunrise"`
		Sunset  string `json:"sunset"`
	} `json:"results"`
	Status string `json:"status"`
}

const apiURL = "http://maps.googleapis.com/maps/api/geocode/json?address=%s,%s,%d"
const apiURL2 = "http://api.sunrise-sunset.org/json?lat=%f&lng=%f&date=today"

// helperfunction for errors
func abortIfError(err error) {
	if err != nil {
		log.Fatalf("%+v\n", err)
		os.Exit(42)
	}
}

func main() {

	// parse command arguments
	city := flag.String("city", "Berlin", "City")
	street := flag.String("street", "Bahnhofstrasse", "Street")
	number := flag.Int("number", 1, "Number")
	flag.Parse()

	// setup API url with parameters
	completeAPIURL := fmt.Sprintf(apiURL, *city, *street, *number)

	fmt.Printf("Retrieving geolocation information for %s, %s %d\n", *city, *street, *number)

	// establish connection
	resp, err := http.Get(completeAPIURL)
	defer resp.Body.Close()
	abortIfError(err)

	// retrieve json
	body, err := ioutil.ReadAll(resp.Body)
	abortIfError(err)

	// unmarshall the json payload into go struct
	payload := new(Payload)
	err = json.Unmarshal(body, payload)
	abortIfError(err)

	longitude := payload.Results[0].Geometry.Location.Logitude
	latitude := payload.Results[0].Geometry.Location.Lattitude

	fmt.Printf("Status: %v - Longitude: %v - Latitude: %v\n", payload.Status, longitude, latitude)

	// get sunset and sunrise time
	// http://sunrise-sunset.org/api

	completeAPIURL = fmt.Sprintf(apiURL2, latitude, longitude)

	fmt.Printf("Retrieving UTC sunrise and sunset information for %f, %f\n", latitude, longitude)

	// establish connection
	resp, err = http.Get(completeAPIURL)
	defer resp.Body.Close()
	abortIfError(err)

	// retrieve json
	body, err = ioutil.ReadAll(resp.Body)
	abortIfError(err)

	// unmarshall the json payload into go struct
	payload2 := new(Payload2)
	err = json.Unmarshal(body, payload2)
	abortIfError(err)

	sunrise := payload2.Results.Sunrise
	sunset := payload2.Results.Sunset

	fmt.Printf("Retrieved UTC sunrise %s and sunset %s for %f, %f\n", sunrise, sunset, latitude, longitude)

}
