package classes

// Samples used in a small go tutorial
//
// Copyright (C) 2017,2018 framp at linux-tips-and-tricks dot de
//
// Samples for go - simple trafficlight simulation using go channels and go routines
//
// See github.com/framps/golang_tutorial for latest code

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/framps/golang_tutorial/trafficLight/v1/globals"
	rpio "github.com/stianeikeland/go-rpio"
)

const gPIOFile = "GPIO.json"

// LEDs - LED GPIO numbers for lights of one traffic light
type LEDs struct {
	Pin [3]int // red, yellow, green
}

// default map of GPIO numbers to BCM GPIO numbers
// change GPIO.json to use other another mapping
// GPIO#                  --->  0   1   2   3   4   5   6   7
var defaultgpio2bcm = [...]int{17, 18, 27, 22, 23, 24, 25, 04}

// LEDController -
type LEDController struct {
	enabled  bool
	gpio2bcm [8]int
}

// FlashLEDs -
func (lc *LEDController) FlashLEDs(t *TrafficLight) {
	l := phaseString[t.program.Phases[t.program.state].Light]
	for i := 0; i < len(l); i += 2 {
		if l[i] == byte('.') {
			lc.Off(t.leds.Pin[i/2])
		} else {
			lc.On(t.leds.Pin[i/2])
		}
	}
}

// NewLEDController -
func NewLEDController() *LEDController {
	l := &LEDController{enabled: globals.EnableLEDs, gpio2bcm: defaultgpio2bcm}
	l.Open()
	l.ClearAll()
	return l
}

// ClearAll -
func (lc *LEDController) ClearAll() {
	if lc.enabled {
		for _, p := range lc.gpio2bcm {
			pin := rpio.Pin(p)
			pin.Output()
			pin.Low()
		}
	}
}

// Open -
func (lc *LEDController) Open() {
	if lc.enabled {
		err := rpio.Open()
		if err != nil {
			fmt.Printf("Error accessing GPIO: %s\n", err.Error())
			os.Exit(42)
		}

		defs, err := lc.ReadGPIOConfig()
		if err == nil {
			if globals.Monitor {
				fmt.Printf("GPIO mapping to pins 0-7: %v\n", defs)
			}
			lc.gpio2bcm = *defs
		} else {
			fmt.Printf("Using default GPIO mapping to pins 0-7: %v\n", defs)
		}
	}
}

// Close -
func (lc *LEDController) Close() {
	if lc.enabled {
		lc.ClearAll()
		rpio.Close()
	}
}

// ReadGPIOConfig -
func (lc *LEDController) ReadGPIOConfig() (*[8]int, error) {

	file, e := ioutil.ReadFile(gPIOFile)
	if e != nil { // error
		fmt.Printf("%s read error: %v\n", gPIOFile, e)
		return nil, e
	}

	var GPIODefs [8]int
	e = json.Unmarshal(file, &GPIODefs)
	if e != nil {
		fmt.Printf("JSON parse error: %v\n", e)
		return nil, e
	}
	return &GPIODefs, nil
}

// On -
func (lc *LEDController) On(gpio int) {
	pin := rpio.Pin(lc.gpio2bcm[gpio])
	pin.Output()
	pin.High()
}

// Off -
func (lc *LEDController) Off(gpio int) {
	pin := rpio.Pin(lc.gpio2bcm[gpio])
	pin.Output()
	pin.Low()
}
