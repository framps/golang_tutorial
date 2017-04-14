// Samples used in a small go tutorial
//
// Copyright (C) 2017 framp at linux-tips-and-tricks dot de
//
// Samples for go - simple traffix light simulation using channels go go routines
//
// See github.com/framps/golang_tutorial for latest code

package main

import (
	"fmt"
	"time"
)

const debug = false

// lamp colors
const (
	green = iota
	yellow
	red
	redyellow
)

// ascii representation of phase lamps (Green, Yellow, Red, )
var phaseString = []string{". . G", ". Y .", "R . .", "R Y ."}

// Phase --
type Phase struct {
	Lamps int
	Ticks int
}

// Program - Has phases and a state (active phase)
type Program struct {
	Phases []Phase
	state  int
}

// NormalProgram -
var NormalProgram = Program{
	Phases: []Phase{
		Phase{green, 3},
		Phase{yellow, 1},
		Phase{red, 3},
		Phase{redyellow, 1},
	},
}

// TrafficLight -
type TrafficLight struct {
	name    string         // name
	ticks   int            // ticks received
	program *Program       // program to execute
	c       *chan struct{} // tick channel to liston on
}

// NewTrafficLight -- Create a new trafficlight
func NewTrafficLight(name string, startPhase int, program Program, c *chan struct{}) (t *TrafficLight) {
	t = &TrafficLight{name, 0, &program, c}
	program.state = startPhase
	return t
}

// Implement Stringer interface
func (t *TrafficLight) String() string {
	return fmt.Sprintf("%v: %v", t.name, phaseString[t.program.Phases[t.program.state].Lamps])
}

// Run --
func (t *TrafficLight) Run() {

	for {
		debugMessage("%s: Waiting ...\n", t.name)
		<-*t.c // wait for tick
		debugMessage("%s: Advancing ...\n", t.name)
		t.Advance() // next trafficlight phase
	}
}

// Advance -- Advances a trafficlight to next phase
func (t *TrafficLight) Advance() {
	debugMessage("%s: Got tick\n", t.name)
	if t.ticks++; t.ticks >= t.program.Phases[t.program.state].Ticks {
		t.program.state = (t.program.state + 1) % len(t.program.Phases)
		t.ticks = 0
		fmt.Printf("%v\n", t)
	}
}

// Debugging helper
func debugMessage(f string, p ...interface{}) {
	if debug {
		fmt.Printf(f, p...)
	}
}

func main() {

	tick1 := make(chan struct{})
	tick2 := make(chan struct{})

	trafficLight1 := NewTrafficLight("Trafficlight1", 1, NormalProgram, &tick1)
	trafficLight2 := NewTrafficLight("Trafficlight2", 3, NormalProgram, &tick2)

	trafficLights := []*TrafficLight{trafficLight1, trafficLight2}

	// start all trafficlights to run parallel as a go routine
	for i := range trafficLights {
		go trafficLights[i].Run()
	}

	for {
		for i := range trafficLights {
			*trafficLights[i].c <- struct{}{} // send new tick
		}
		time.Sleep(time.Second * 1)
	}
}
