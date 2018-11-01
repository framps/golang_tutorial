package classes

// Samples used in a small go tutorial
//
// Copyright (C) 2017,2018 framp at linux-tips-and-tricks dot de
//
// Samples for go - simple trafficlight simulation using go channels and go routines
//
// See github.com/framps/golang_tutorial for latest code

import (
	"fmt"
	"time"

	"github.com/framps/golang_tutorial/trafficLight/v1/globals"
)

// TrafficManager -
type TrafficManager struct {
	trafficLights []*TrafficLight
	program       *Program
	onoff         chan bool
	ledController *LEDController // controller to driver LEDs
}

// NewTrafficManager -
func NewTrafficManager(ledController *LEDController, trafficLights []*TrafficLight) *TrafficManager {
	tm := &TrafficManager{trafficLights: trafficLights, ledController: ledController}
	tm.LoadProgram(ProgramTest)
	tm.onoff = make(chan bool)
	return tm
}

// LoadProgram - load new program in trafficlights
func (tm *TrafficManager) LoadProgram(program *Program) {
	tm.program = program
	idxint := 0
	for i := range tm.trafficLights {
		tm.trafficLights[i].Load(idxint, *tm.program)
		idxint = (idxint + len(tm.program.Phases)/2) % len(tm.program.Phases)
	}
}

// Start - Start trafficmanager and manage trafficlights
func (tm *TrafficManager) Start() {
	d := make(chan int)

	// Display trafficlights on terminal if enabled
	go func() {
		cnt := 0
		for {
			<-d
			if globals.Monitor {
				cnt++
				if cnt >= len(tm.trafficLights) {
					for i := range tm.trafficLights {
						fmt.Printf("%s   ", tm.trafficLights[i].String())
					}
					fmt.Println()
					cnt = 0
				}
			}
		}
	}()

	// start all trafficlights to run in parallel
	for i := range tm.trafficLights {
		tm.trafficLights[i].On(d, tm.ledController)
	}

	// send ticks to traffic lights
	go func() {
		for {
			for i := range tm.trafficLights {
				tm.trafficLights[i].c <- struct{}{} // send new tick
			}
			time.Sleep(tm.program.clockSpeed)
		}
	}()
}

// Stop - Stop trafficmanager and trafficlights
func (tm *TrafficManager) Stop() {
	for _, l := range tm.trafficLights {
		l.Off()
	}
	tm.ledController.Close()
}
