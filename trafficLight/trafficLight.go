package main

import (
	"fmt"
	"time"
)

const debug = false

const (
	green = iota
	yellow
	red
	redyellow
)

var phaseString = []string{". . N", ". G .", "R . .", "R G ."}

// Phase --
type Phase struct {
	Lamp  int
	Ticks int
}

// Program -
type Program struct {
	Phases []Phase
}

// NormalProgram -
var NormalProgram = Program{
	Phases: []Phase{
		Phase{
			Lamp:  green,
			Ticks: 3,
		},
		Phase{
			Lamp:  yellow,
			Ticks: 1,
		},
		Phase{
			Lamp:  red,
			Ticks: 3,
		},
		Phase{
			Lamp:  redyellow,
			Ticks: 1,
		},
	},
}

// TrafficLight -
type TrafficLight struct {
	name    string
	state   int
	ticks   int
	program Program
	c       *chan struct{}
}

// Run --
func (t *TrafficLight) Run() {

	for {
		debugMessage("%s: Waiting ...\n", t.name)
		<-*t.c
		debugMessage("%s: Advancing ...\n", t.name)
		t.Advance()
	}
}

// Advance --
func (t *TrafficLight) Advance() {
	debugMessage("%s: Got tick\n", t.name)
	if t.ticks++; t.ticks >= t.program.Phases[t.state].Ticks {
		t.state = (t.state + 1) % len(t.program.Phases)
		t.ticks = 0
		fmt.Printf("%v: %v\n", t.name, phaseString[t.state])
	}
}

func debugMessage(f string, p ...interface{}) {
	if debug {
		fmt.Printf(f, p...)
	}
}

func main() {

	tick1 := make(chan struct{})
	tick2 := make(chan struct{})

	var trafficLight1 = TrafficLight{"Trafficlight1", 1, 0, NormalProgram, &tick1}
	var trafficLight2 = TrafficLight{"Trafficlight2", 3, 0, NormalProgram, &tick2}

	trafficLights := []TrafficLight{trafficLight1, trafficLight2}

	for i := range trafficLights {
		go trafficLights[i].Run()
	}

	for {
		for i := range trafficLights {
			*trafficLights[i].c <- struct{}{}
		}
		time.Sleep(time.Second * 1)
	}
}
