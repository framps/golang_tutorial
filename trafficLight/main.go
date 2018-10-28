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
	// GPIO numbers to BCM GPIO numbers mapping defined in GPIO.json
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

	ctrlc := make(chan os.Signal, 1)
	signal.Notify(ctrlc, os.Interrupt, syscall.SIGHUP,
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	// watch for CTRLC
	go func() {
		<-ctrlc
		done <- struct{}{}
	}()

	type ProgramChunk struct {
		program  *classes.Program
		duration time.Duration
	}

	// programs to run
	programs := []ProgramChunk{
		ProgramChunk{classes.ProgramWarning, time.Second * 5},
		ProgramChunk{classes.ProgramNormal1, time.Second * 15},
		ProgramChunk{classes.ProgramWarning, time.Second * 5},
		ProgramChunk{classes.ProgramNormal2, time.Second * 15},
		ProgramChunk{classes.ProgramWarning, time.Second * 5},
		ProgramChunk{classes.ProgramNormal3, time.Second * 15},
		ProgramChunk{classes.ProgramWarning, time.Second * 5},
		ProgramChunk{classes.ProgramNormal4, time.Second * 15},
	}

	// start manager
	tm.Start()

	// loop though list of programs
	for {
		for _, p := range programs {
			if globals.Monitor {
				fmt.Printf("Running program %s for %s\n", p.program.Name, p.duration)
			}
			tm.LoadProgram(p.program)
			time.Sleep(p.duration)
			select {
			case <-done:
				lc.Close()
				os.Exit(1)
			default:
			}
		}
	}

}
