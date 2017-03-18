package main

import (
	"fmt"
	"time"
)

const debug = true

const (
	green = iota
	yellow
	red
	redyellow
)

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
}

// Advance --
func (t *TrafficLight) Advance() {
	fmt.Printf("Got tick %v\n", t.name)
	if t.ticks++; t.ticks >= t.program.Phases[t.state].Ticks {
		t.state = (t.state + 1) % len(t.program.Phases)
		t.ticks = 0
		fmt.Printf("Advanced %v to %v\n", t.name, t.state)
	}
}

/*
// Controller -
type Controller struct {
	TrafficLights []TrafficLight
}

func (c *Controller) Start() {

	for _,light := c.TrafficLights range {
		go
	}
}
*/

func main() {

	var trafficLight1 = TrafficLight{"Trafficlight1", 0, 0, NormalProgram}
	var trafficLight2 = TrafficLight{"Trafficlight2", 2, 0, NormalProgram}
/*
	var controller = Controller{[]TrafficLight{trafficLight1, trafficLight2}}
*/
	for {
		time.Sleep(time.Second * 3)
		trafficLight1.Advance()
		trafficLight2.Advance()

	}
}
