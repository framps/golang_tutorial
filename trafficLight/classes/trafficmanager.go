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
	program       *Program
}

// NewTrafficManager -
func NewTrafficManager(trafficLights []*TrafficLight, ledController *LEDController) *TrafficManager {
	tm := &TrafficManager{trafficLights: trafficLights, lc: ledController}
	tm.StartProgram(ProgramWarning)
	return tm
}

// StartProgram -
func (tm *TrafficManager) StartProgram(program *Program) {
	tm.program = program
	idxint := 0
	for i := range tm.trafficLights {
		tm.trafficLights[i].Load(idxint, *tm.program)
		idxint = (idxint + len(tm.trafficLights)/2) % len(tm.trafficLights)
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
