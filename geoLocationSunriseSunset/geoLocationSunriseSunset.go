// Samples used in a small go tutorial
//
// Copyright (C) 2017 framp at linux-tips-and-tricks dot de
//
// Samples for go syntax - simple REST client which calculates the sunrise and sunset time of a location

package main

import (
	"encoding/json"
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

const apiURL = "http://maps.googleapis.com/maps/api/geocode/json?address="

// helperfunction for errors
func abortIfError(err error) {
	if err != nil {
		log.Fatalf("%+v\n", err)
		os.Exit(42)
	}
}

func main() {

	// establish connection
	resp, err := http.Get(apiURL)
	defer resp.Body.Close()
	abortIfError(err)

	// retrieve json
	body, err := ioutil.ReadAll(resp.Body)
	abortIfError(err)

	// unmarshall the json payload into go struct
	payload := new(Payload)
	err = json.Unmarshal(body, payload)
	abortIfError(err)

	fmt.Printf("Status: %v - Long: %v - Lat: %v\n", payload.Status, payload.Results[0].Geometry.Location.Logitude, payload.Results[0].Geometry.Location.Lattitude)
}
