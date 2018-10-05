package trafficmanager

import (
	"fmt"
	"sync"
	"time"
)

// TrafficManager -
type TrafficManager struct {
	trafficLights []*TrafficLight
}

// NewTrafficManager -
func NewTrafficManager(trafficLights []*TrafficLight) *TrafficManager {
	tm := &TrafficManager{trafficLights: trafficLights}
	return tm
}

// TestMode -
func (tm *TrafficManager) TestMode() {
	var flip bool
	for i := range tm.trafficLights {
		if flip {
			tm.trafficLights[i].Load(2, TestProgram)
		} else {
			tm.trafficLights[i].Load(0, TestProgram)
		}
		flip = !flip
	}
}

// Start -
func (tm *TrafficManager) Start() {
	var flip bool
	for i := range tm.trafficLights {
		if flip {
			tm.trafficLights[i].Load(2, NormalProgram)
		} else {
			tm.trafficLights[i].Load(0, NormalProgram)
		}
		flip = !flip
	}
}

// On -
func (tm *TrafficManager) On(wg *sync.WaitGroup) {

	d := make(chan int)

	go func(update chan int) {
		var cnt int
		for {
			<-update
			cnt++
			if cnt >= len(tm.trafficLights) {
				for i := range tm.trafficLights {
					if monitor {
						fmt.Printf("%s   ", tm.trafficLights[i].String())
					}
					if enableLEDs {
						tm.trafficLights[i].FlashLEDs()
					}
				}
				if monitor {
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
			time.Sleep(time.Second * 1)
		}
	}()
}
