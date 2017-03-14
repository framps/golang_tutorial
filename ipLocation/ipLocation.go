// Samples used in a small go tutorial
//
// Copyright (C) 2017 framp at linux-tips-and-tricks dot de
//
// Samples for go syntax - simple REST client which returns the city the client is located from the client's IP

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

// Location - JSON payload returned by REST API
type Location struct {
	As           string  `json:"as"` // json element 'as' will be unmarshalled into struct element 'As'
	City         string  `json:"city"`
	Country      string  `json:"country"`
	CountryCode  string  `json:"countrycode"`
	Isp          string  `json:"isp"`
	Latitude     float32 `json:"lat"`
	Longitude    float32 `json:"lon"`
	Organization string  `json:"org"`
	Query        string  `json:"query"`
	Region       string  `json:"region"`
	RegionName   string  `json:"regionname"`
	Status       string  `json:"status"`
	Timezone     string  `json:"timezone"`
	Zip          string  `json:"zip"`
}

const apiURL = "http://ip-api.com/json"

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
	location := new(Location)
	err = json.Unmarshal(body, location)
	abortIfError(err)

	// print information
	fmt.Printf("Information retrived for IP %v\nCountry code: %v\nCity: %v\nProvider: %v\n",
		location.Query, location.CountryCode, location.City, location.Isp)

}
