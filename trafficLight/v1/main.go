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

	"github.com/framps/golang_tutorial/trafficLight/v1/classes"
	"github.com/framps/golang_tutorial/trafficLight/v1/globals"
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
		classes.NewTrafficLight(0, T1LEDs),
		classes.NewTrafficLight(1, T2LEDs)}

	tm := classes.NewTrafficManager(lc, trafficLights)

	ctrlc := make(chan os.Signal, 1)
	signal.Notify(ctrlc, os.Interrupt, syscall.SIGTERM, syscall.SIGHUP)

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
	}

	// start traffic manager
	tm.Start()

	// listen for CTRLC
	done := make(chan struct{})
	go func() {
		<-ctrlc
		fmt.Printf("CTRLC")
		done <- struct{}{}
	}()

	// loop though list of programs
	go func() {
		for {
			for _, p := range programs {
				if globals.Monitor {
					fmt.Printf("Running program %s for %s\n", p.program.Name, p.duration)
				}
				tm.LoadProgram(p.program)
				time.Sleep(p.duration)
			}
		}
	}()

	// wait for CTRLC
	<-done
	fmt.Print("Done received")
	tm.Stop()
	os.Exit(0)
}
