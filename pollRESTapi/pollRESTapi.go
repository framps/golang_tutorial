// Samples used in a small go tutorial
//
// Copyright (C) 2017 framp at linux-tips-and-tricks dot de
//
// Samples for go syntax - simple REST client which polls a REST API server
//
// See github.com/framps/golang_tutorial for latest code

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"sync"
	"time"
)

const serverPort = 8080

type SimpleRESTResponse struct {
	TimeNow string `json:"timenow"`
}

var actualTime time.Time

func startServer(port int) {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		rsp := SimpleRESTResponse{actualTime.Truncate(time.Second).Format("15:04:05")}

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(rsp); err != nil {
			log.Fatal(err)
		}
	})

	fmt.Printf("Server listening on %s\n", fmt.Sprintf("localhost:%d", port))
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}

func startUpdater() {

	for {
		actualTime = time.Now()
		fmt.Printf("<-: %s\n", actualTime.Truncate(time.Second).Format("15:04:05"))
		time.Sleep(time.Second)
	}

}

func startClient(serverPort int, id int) {

	fmt.Printf("Client %d polling %s\n", id, fmt.Sprintf("http://localhost:%d", serverPort))

	for {

		rsp, err := http.Get(fmt.Sprintf("http://localhost:%d", serverPort))
		if err != nil {
			log.Fatal(err)
		}
		defer rsp.Body.Close()
		body, err := ioutil.ReadAll(rsp.Body)

		var response SimpleRESTResponse
		json.Unmarshal(body, &response)

		fmt.Printf(" %d: %s ->\n", id, response.TimeNow)

		o := time.Millisecond * time.Duration(rand.Intn(500)) // add some random delay

		time.Sleep(time.Second*3 + o)
	}

}

func main() {

	const numClients = 3
	var wg sync.WaitGroup

	fmt.Println("Starting time updater")
	go startUpdater()
	wg.Add(1)

	fmt.Println("Starting REST API server")
	go startServer(serverPort)
	wg.Add(1)
	time.Sleep(time.Second * 3)

	for i := 1; i <= numClients; i++ {
		fmt.Println("Starting REST client %d", i)
		go startClient(serverPort, i)
		wg.Add(1)
	}

	wg.Wait()

}
