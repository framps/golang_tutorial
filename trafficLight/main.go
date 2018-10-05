// Samples used in a small go tutorial
//
// Copyright (C) 2017,2018 framp at linux-tips-and-tricks dot de
//
// Samples for go - simple trafficlight simulation using go channels and go routines
//
// See github.com/framps/golang_tutorial for latest code

package main

import (
	"flag"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/framps/golang_tutorial/trafficLight/classes"
	"github.com/framps/golang_tutorial/trafficLight/globals"
	rpio "github.com/stianeikeland/go-rpio"
)

func main() {

	flag.BoolVar(&globals.Debug, "debug", false, "Write debug messages")
	flag.BoolVar(&globals.Monitor, "monitor", true, "Monitor LEDs on screen")
	flag.BoolVar(&globals.EnableLEDs, "leds", false, "Drive LEDs")

	flag.Parse()

	if globals.EnableLEDs {
		err := rpio.Open()
		if err != nil {
			fmt.Printf("Error accessing GPIO: %s\n", err.Error())
			os.Exit(42)
		}
	}

	trafficLight1 := classes.NewTrafficLight(0, classes.T1LEDs)
	trafficLight2 := classes.NewTrafficLight(1, classes.T2LEDs)

	trafficLights := []*classes.TrafficLight{trafficLight1, trafficLight2}

	tm := classes.NewTrafficManager(trafficLights)

	var wg sync.WaitGroup
	wg.Add(1)
	tm.On(&wg)

	go func() {
		time.Sleep(time.Second * 5)
		tm.TestMode()
		time.Sleep(time.Second * 5)
		tm.Start()
	}()

	wg.Wait()
}
