package classes

import (
	"testing"
	"time"

	"github.com/framps/golang_tutorial/trafficLight/globals"
)

func TestCreateAndIdle(t *testing.T) {

	T1LEDs := LEDs{[...]int{2, 3, 4}}

	globals.Debug = true

	lc := NewLEDController()
	lc.Enable()
	tl := NewTrafficLight(0, T1LEDs, lc)

	var stop chan struct{}
	t.Log("Starting trafficlight")
	tl.Start(stop)

	debugMessage("Starting timer")
	timer := time.NewTimer(time.Second * 10)
	<-timer.C
	debugMessage("Timer fired")
	t.Log("Stopping trafficlight")
	stop <- struct{}{}
	debugMessage("Stop sent")
}
