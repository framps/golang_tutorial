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
	"os/signal"
	"syscall"
	"time"

	"github.com/framps/golang_tutorial/trafficLight/classes"
	"github.com/framps/golang_tutorial/trafficLight/globals"
)

func main() {

	// GPIO#s: red, yellow, green
	var (
		T1LEDs = classes.LEDs{[...]int{2, 3, 4}}
		T2LEDs = classes.LEDs{[...]int{5, 6, 7}}
	)

	flag.BoolVar(&globals.Debug, "debug", false, "Write debug messages")
	flag.BoolVar(&globals.Monitor, "monitor", true, "Monitor LEDs on screen")
	flag.BoolVar(&globals.EnableLEDs, "leds", false, "Drive LEDs")
	flag.Parse()

	lc := classes.NewLEDController()

	trafficLights := []*classes.TrafficLight{
		classes.NewTrafficLight(0, T1LEDs, lc),
		classes.NewTrafficLight(1, T2LEDs, lc)}

	tm := classes.NewTrafficManager(trafficLights)

	done := make(chan struct{})

	c := make(chan os.Signal, 1)
	//	signal.Notify(c, os.Interrupt, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	signal.Notify(c, syscall.SIGTERM)
	go func() {
		<-c
		done <- struct{}{}
	}()

	type ProgramChunk struct {
		program  *classes.Program
		duration time.Duration
	}

	programs := []ProgramChunk{
		ProgramChunk{classes.ProgramNormal2, time.Second * 5},
		ProgramChunk{classes.ProgramWarning, time.Second * 3},
		ProgramChunk{classes.ProgramNormal3, time.Second * 5},
		ProgramChunk{classes.ProgramWarning, time.Second * 3},
	}

	tm.On()

	for {
		for _, p := range programs {
			fmt.Printf("Program %s: Sleeping %s\n", p.program.Name, p.duration)
			time.Sleep(p.duration)
			tm.LoadProgram(p.program)
			select {
			case <-done:
				lc.Close()
				os.Exit(1)
			}
		}
	}

}
