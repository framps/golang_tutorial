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

	"github.com/framps/golang_tutorial/trafficLight/globals"
)

// TrafficManager -
type TrafficManager struct {
	trafficLights []*TrafficLight
	program       *Program
}

// NewTrafficManager -
func NewTrafficManager(trafficLights []*TrafficLight) *TrafficManager {
	tm := &TrafficManager{trafficLights: trafficLights}
	tm.LoadProgram(ProgramNormal)
	return tm
}

// LoadProgram -
func (tm *TrafficManager) LoadProgram(program *Program) {
	debugMessage("Loading program %s\n", program.Name)
	tm.program = program
	idxint := 0
	for i := range tm.trafficLights {
		debugMessage("%d: Loading %d - Phase: %d\n", i, idxint, len(tm.program.Phases))
		tm.trafficLights[i].Load(idxint, *tm.program)
		idxint = (idxint + len(tm.program.Phases)/2) % len(tm.program.Phases)
	}
}

// On -
func (tm *TrafficManager) On() {

	d := make(chan int)

	// Display trafficlights
	go func() {
		cnt := 0
		for {
			n := <-d
			cnt++
			debugMessage("TM: Got update from %d (%d)\n", n, cnt)
			if cnt >= len(tm.trafficLights) {
				for i := range tm.trafficLights {
					if globals.Monitor {
						fmt.Printf("%s   ", tm.trafficLights[i].String())
					}
				}
				if globals.Monitor {
					fmt.Println()
				}
				cnt = 0
			}
		}
	}()

	// start raffigLightall trafficlights to run in parallel
	for i := range tm.trafficLights {
		go tm.trafficLights[i].Run(d)
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
