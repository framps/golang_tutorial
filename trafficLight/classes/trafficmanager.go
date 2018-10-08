package classes

import (
	"fmt"
	"time"

	"github.com/framps/golang_tutorial/trafficLight/globals"
)

// TrafficManager -
type TrafficManager struct {
	trafficLights []*TrafficLight
	lc            *LEDController
	program       Program
}

// NewTrafficManager -
func NewTrafficManager(trafficLights []*TrafficLight, ledController *LEDController, program *Program) *TrafficManager {
	tm := &TrafficManager{trafficLights: trafficLights, lc: ledController}
	tm.WarningMode()
	return tm
}

// WarningMode -
func (tm *TrafficManager) WarningMode() {
	tm.program = WarningProgram
	var flip bool
	for i := range tm.trafficLights {
		if flip {
			tm.trafficLights[i].Load(0, tm.program)
		} else {
			tm.trafficLights[i].Load(0, tm.program)
		}
		flip = !flip
	}
}

// TestMode -
func (tm *TrafficManager) TestMode() {
	tm.program = TestProgram
	var flip bool
	for i := range tm.trafficLights {
		if flip {
			tm.trafficLights[i].Load(0, tm.program)
		} else {
			tm.trafficLights[i].Load(0, tm.program)
		}
		flip = !flip
	}
}

// Start -
func (tm *TrafficManager) Start() {
	tm.program = NormalProgram
	var flip bool
	for i := range tm.trafficLights {
		if flip {
			tm.trafficLights[i].Load(2, tm.program)
		} else {
			tm.trafficLights[i].Load(0, tm.program)
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
					if globals.Monitor {
						fmt.Printf("%s   ", tm.trafficLights[i].String())
					}
					if globals.EnableLEDs {
						tm.trafficLights[i].FlashLEDs(tm.lc)
					}
				}
				if globals.Monitor {
					fmt.Println()
				}
				cnt = 0
			}
		}
	}(d)

	// start all trafficlights to run parallel as a go routine
	for i := range tm.trafficLights {
		go tm.trafficLights[i].Run(d)
	}

	go func() {
		for {
			for i := range tm.trafficLights {
				tm.trafficLights[i].c <- struct{}{} // send new tick
			}
			time.Sleep(tm.program.clockSpeed)
		}
	}()
}
