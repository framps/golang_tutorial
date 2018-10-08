package classes

import (
	"fmt"

	"github.com/framps/golang_tutorial/trafficLight/constants"
)

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
	t = &TrafficLight{number: number, ticks: 0, program: ProgramWarning, leds: leds, c: c}
	t.program.state = 1
	return t
}

// Load -
func (t *TrafficLight) Load(startPhase int, program *Program) {
	t.program = program
	t.program.state = startPhase
}

// Implement Stringer interface to display a readable form of the traffic light
func (t *TrafficLight) String() string {
	return fmt.Sprintf("<%d>: %s |", t.number, constants.PhaseString[t.program.Phases[t.program.state].Lights])
}

// FlashLEDs -
func (t *TrafficLight) FlashLEDs(lightController *LEDController) {
	l := constants.PhaseString[t.program.Phases[t.program.state].Lights]
	for i := 0; i < len(l); i += 2 {
		if l[i] == byte('.') {
			lightController.Off(t.leds.Pin[i/2])
		} else {
			lightController.On(t.leds.Pin[i/2])
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
}
