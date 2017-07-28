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

// enable do print debug messages
const debug = false

const enableLEDs = false

// lamp colors
const (
	green = iota
	yellow
	red
	redyellow
)

// ascii representation of phase lamps (Green, Yellow, Red, Read and Yellow )
var phaseString = []string{". . G", ". Y .", "R . .", "R Y ."}

var t1LEDs = LEDs{[...]int{11, 12, 13}}
var t2LEDs = LEDs{[...]int{23, 24, 25}}

// LEDs - LED pin numbers for lights of one traffic light
type LEDs struct {
	pin [3]int // red, yellow, green
}

// Phase consists of lights and number of ticks to flash the linghts
type Phase struct {
	Lights int
	Ticks  int
}

// Program - Has phases and a state (active phase)
type Program struct {
	Phases []Phase
	state  int
}

// NormalProgram - Just the common traffic light
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
	leds    LEDs           // LEDs to use
	c       *chan struct{} // tick channel to liston on
}

// NewTrafficLight -- Create a new trafficlight
func NewTrafficLight(name string, startPhase int, program Program, c *chan struct{}, leds LEDs) (t *TrafficLight) {
	t = &TrafficLight{name, 0, &program, leds, c}
	program.state = startPhase
	return t
}

// Implement Stringer interface to display a readable form of the traffic light
func (t *TrafficLight) String() string {
	return fmt.Sprintf("%v: %v", t.name, phaseString[t.program.Phases[t.program.state].Lights])
}

// FlashLEDs -
func (t *TrafficLight) FlashLEDs() {
	l := phaseString[t.program.Phases[t.program.state].Lights]
	for i := 0; i < len(l); i += 2 {
		if l[i] == byte('.') {
			fmt.Printf("off %d ", t.leds.pin[i/2])
		} else {
			fmt.Printf("on %d ", t.leds.pin[i/2])
		}
	}
}

// Run - run traffic light program
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
	}
	fmt.Printf("%v ", t)
	if enableLEDs {
		t.FlashLEDs()
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

	trafficLight1 := NewTrafficLight("T1", 1, NormalProgram, &tick1, t1LEDs)
	trafficLight2 := NewTrafficLight("T2", 3, NormalProgram, &tick2, t2LEDs)

	trafficLights := []*TrafficLight{trafficLight1, trafficLight2}

	// start all trafficlights to run parallel as a go routine
	for i := range trafficLights {
		go trafficLights[i].Run()
	}

	for {
		for i := range trafficLights {
			*trafficLights[i].c <- struct{}{} // send new tick
		}
		fmt.Println()
		time.Sleep(time.Second * 1)
	}
}
