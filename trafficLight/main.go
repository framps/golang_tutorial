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
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/framps/golang_tutorial/trafficLight/classes"
	"github.com/framps/golang_tutorial/trafficLight/globals"
)

func main() {

	// GPIOs: red, yellow, green
	var (
		T1LEDs = classes.LEDs{[...]int{2, 3, 4}}
		T2LEDs = classes.LEDs{[...]int{5, 6, 7}}
	)

	flag.BoolVar(&globals.Debug, "debug", false, "Write debug messages")
	flag.BoolVar(&globals.Monitor, "monitor", true, "Monitor LEDs on screen")
	flag.BoolVar(&globals.EnableLEDs, "leds", false, "Drive LEDs")

	flag.Parse()

	trafficLight1 := classes.NewTrafficLight(0, T1LEDs)
	trafficLight2 := classes.NewTrafficLight(1, T2LEDs)
	trafficLights := []*classes.TrafficLight{trafficLight1, trafficLight2}

	lc := classes.NewLEDController()
	tm := classes.NewTrafficManager(trafficLights, lc)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		<-c
		lc.Close()
		os.Exit(1)
	}()

	tm.On()

	for {
		time.Sleep(time.Second * 5)
		tm.StartProgram(classes.ProgramTest)
		time.Sleep(time.Second * 5)
		tm.StartProgram(classes.ProgramNormal)
		time.Sleep(time.Second * 30)
		tm.StartProgram(classes.ProgramWarning)

	}

}
