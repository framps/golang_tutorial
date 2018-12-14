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
)

// ascii representation of phase lamps (Green, Yellow, Red, Read and Yellow )
var phaseString = []string{". . .", ". . G", ". Y .", "R . .", "R Y ."}

// TrafficLight -
type TrafficLight struct {
	number  int            // light number
	ticks   int            // ticks received
	program Program        // program to execute
	leds    LEDs           // LEDs to use
	c       chan struct{}  // tick channel to liston on
	lc      *LEDController // controller to driver LEDs
	on      bool           // on/off flag
}

// NewTrafficLight -- Create a new trafficlight
func NewTrafficLight(number int, leds LEDs) (t *TrafficLight) {
	c := make(chan struct{})
	t = &TrafficLight{number: number, ticks: 0,
		program: *ProgramTest, leds: leds, c: c}
	t.program.state = 1
	return t
}

// Load - Load new program
func (t *TrafficLight) Load(startPhase int, program Program) {
	t.program = program
	t.program.state = startPhase
}

// Implement Stringer interface to display a readable form of the traffic light
func (t *TrafficLight) String() string {
	return fmt.Sprintf("<%d>: %s |", t.number, phaseString[t.program.Phases[t.program.state].Lights])
}

// On - Turn trafficlight on
func (t *TrafficLight) On(callBack chan int, lc *LEDController) {
	t.lc = lc
	t.on = true
	go func() {
		for {
			<-t.c // wait for tick
			if !t.on {
				break
			}
			t.Advance() // next trafficlight phase
			if EnableLEDs {
				t.lc.FlashLEDs(t)
			}
			callBack <- t.number
		}
	}()
}

// Off - Turn trafficlight off
func (t *TrafficLight) Off() {
	t.on = false
}

// Advance -- Advances a trafficlight to next phase
func (t *TrafficLight) Advance() {
	if t.ticks++; t.ticks >= t.program.Phases[t.program.state].Ticks {
		t.program.state = (t.program.state + 1) % len(t.program.Phases)
		t.ticks = 0
	}
}
