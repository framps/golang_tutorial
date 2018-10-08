package classes

import (
	"fmt"
	"os"

	"github.com/framps/golang_tutorial/trafficLight/globals"
	rpio "github.com/stianeikeland/go-rpio"
)

// LEDs - LED GPIO numbers for lights of one traffic light
type LEDs struct {
	Pin [3]int // red, yellow, green
}

// map GPIO numbers to BCM GPIO numbers
//                     0   1   2   3   4   5   6   7
var gpio2bcm = [8]int{17, 18, 27, 22, 23, 24, 25, 4}

type LEDController struct {
	enabled bool
}

func NewLEDController() *LEDController {
	l := &LEDController{enabled: globals.EnableLEDs}
	l.Open()
	return l
}

func (l *LEDController) ClearAll() {
	if l.enabled {
		for _, p := range gpio2bcm {
			pin := rpio.Pin(p)
			pin.Output()
			pin.Low()
		}
	}
}

func (l *LEDController) Open() {
	if l.enabled {
		err := rpio.Open()
		if err != nil {
			fmt.Printf("Error accessing GPIO: %s\n", err.Error())
			os.Exit(42)
		}
	}
}

func (l *LEDController) Close() {
	if l.enabled {
		l.ClearAll()
		rpio.Close()
	}
}

func (l *LEDController) On(gpio int) {
	pin := rpio.Pin(gpio2bcm[gpio])
	pin.Output()
	pin.High()
}

func (l *LEDController) Off(gpio int) {
	pin := rpio.Pin(gpio2bcm[gpio])
	pin.Output()
	pin.Low()
}
