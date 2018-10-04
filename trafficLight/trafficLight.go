// Samples used in a small go tutorial
//
// Copyright (C) 2017,2018 framp at linux-tips-and-tricks dot de
//
// Samples for go - simple trafficlight simulation using go channels and go routines
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
	off = iota
	green
	yellow
	red
	redyellow
)

// ascii representation of phase lamps (Green, Yellow, Red, Read and Yellow )
var phaseString = []string{". . .", ". . G", ". Y .", "R . .", "R Y ."}

var t1LEDs = LEDs{[...]int{11, 12, 13}}
var t2LEDs = LEDs{[...]int{21, 22, 23}}
var t3LEDs = LEDs{[...]int{31, 32, 33}}
var t4LEDs = LEDs{[...]int{41, 42, 43}}

// LEDs - LED pin numbers for lights of one traffic light
type LEDs struct {
	pin [3]int // red, yellow, green
}

// Phase consists of light and number of ticks to flash the lights
type Phase struct {
	Lights int
	Ticks  int
}

// Program - Has phases and a state (active phase)
type Program struct {
	Phases []Phase
	state  int
}

// TestProgram - Turn every lamp on
var TestProgram = Program{
	Phases: []Phase{
		Phase{red, 1},
		Phase{yellow, 1},
		Phase{green, 1},
		Phase{off, 1},
		Phase{green, 1},
		Phase{yellow, 1},
		Phase{red, 1},
		Phase{off, 1},
	},
}

// WarningProgram - Traffic light is not working, just blink
var WarningProgram = Program{
	Phases: []Phase{
		Phase{yellow, 1},
		Phase{off, 1},
	},
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

// TrafficManager -
type TrafficManager struct {
	trafficLights []TrafficLight
}

// NewTrafficManager -
func NewTrafficManager(trafficLights []TrafficLight) *TrafficManager {
	tm := &TrafficManager{trafficLights: trafficLights}
	return tm
}

// Test -
func (tm *TrafficManager) Test() {
	var flip bool
	for i := range tm.trafficLights {
		if flip {
			tm.trafficLights[i].Load(5, TestProgram)
		} else {
			tm.trafficLights[i].Load(1, TestProgram)
		}
		flip = !flip
	}
}

// Start -
func (tm *TrafficManager) Start() {
	var flip bool
	for i := range tm.trafficLights {
		if flip {
			tm.trafficLights[i].Load(3, NormalProgram)
		} else {
			tm.trafficLights[i].Load(1, NormalProgram)
		}
		flip = !flip
	}
}

// On -
func (tm *TrafficManager) On() {

	d := make(chan int)

	go func(update chan int) {
		var cnt int
		for {
			<-update
			cnt++
			if cnt >= len(tm.trafficLights) {
				for i := range tm.trafficLights {
					fmt.Printf("%s   ", tm.trafficLights[i].String())
				}
				fmt.Println()
				cnt = 0
			}
		}
	}(d)

	// start all trafficlights to run parallel as a go routine
	for i := range tm.trafficLights {
		go tm.trafficLights[i].Run(d)
	}

	for {
		for i := range tm.trafficLights {
			tm.trafficLights[i].c <- struct{}{} // send new tick
		}
		time.Sleep(time.Second * 1)
	}
}

// TrafficLight -
type TrafficLight struct {
	number  int           // light number
	ticks   int           // ticks received
	program Program       // program to execute
	leds    LEDs          // LEDs to use
	c       chan struct{} // tick channel to liston on
}

// NewTrafficLight -- Create a new trafficlight
func NewTrafficLight(number int, leds LEDs) (t *TrafficLight) {
	c := make(chan struct{})
	t = &TrafficLight{number, 0, WarningProgram, leds, c}
	t.program.state = 1
	return t
}

// Load -
func (t *TrafficLight) Load(startPhase int, program Program) {
	t.program = program
	t.program.state = startPhase
}

// Implement Stringer interface to display a readable form of the traffic light
func (t *TrafficLight) String() string {
	return fmt.Sprintf("<%d>: %s |", t.number, phaseString[t.program.Phases[t.program.state].Lights])
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
func (t *TrafficLight) Run(callBack chan int) {

	for {
		debugMessage("%v: Waiting ...\n", t.number)
		<-t.c // wait for tick
		debugMessage("%v: Advancing ...\n", t.number)
		t.Advance() // next trafficlight phase
		callBack <- t.number
	}
}

// Advance -- Advances a trafficlight to next phase
func (t *TrafficLight) Advance() {
	debugMessage("%v: Got tick\n", t.number)
	if t.ticks++; t.ticks >= t.program.Phases[t.program.state].Ticks {
		debugMessage("%v: Next phase\n", t.number)
		t.program.state = (t.program.state + 1) % len(t.program.Phases)
		t.ticks = 0
	}
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

	trafficLight1 := NewTrafficLight(0, t1LEDs)
	trafficLight2 := NewTrafficLight(1, t2LEDs)
	trafficLight3 := NewTrafficLight(2, t3LEDs)
	trafficLight4 := NewTrafficLight(3, t4LEDs)

	trafficLights := []TrafficLight{*trafficLight1, *trafficLight2, *trafficLight3, *trafficLight4}

	tm := NewTrafficManager(trafficLights)

	go func() {
		time.Sleep(time.Second * 10)
		tm.Test()
		time.Sleep(time.Second * 10)
		tm.Start()
	}()

	tm.On()

}
