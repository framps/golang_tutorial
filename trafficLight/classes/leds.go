package classes

import (
	"fmt"
	"os"

	"github.com/framps/golang_tutorial/trafficLight/globals"
	rpio "github.com/stianeikeland/go-rpio"
)

// LEDs - LED pin numbers for lights of one traffic light
type LEDs struct {
	Pin [3]int // red, yellow, green
}

type LEDController struct {
	enabled bool
	pins    [8]int
}

func NewLEDController(pins [8]int) *LEDController {
	l := &LEDController{pins: pins, enabled: globals.EnableLEDs}
	l.Open()
	return l
}

func (l *LEDController) ClearAll() {
	if l.enabled {
		for _, p := range l.pins {
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
	pin := rpio.Pin(l.pins[gpio])
	pin.Output()
	pin.High()
}

func (l *LEDController) Off(gpio int) {
	pin := rpio.Pin(l.pins[gpio])
	pin.Output()
	pin.Low()
}
