// Samples used in a small go tutorial
//
// Turns on and off all GPIOs in sequence
//
// Modified sample blink program from https://github.com/mrmorphic/hwio/tree/master/examples
// For more samples see https://github.com/framps/golang_tutorial
//
// Copyright (C) 2017 framp at linux-tips-and-tricks dot de

package main

import (
	"fmt"
	"os"

	"github.com/mrmorphic/hwio"
)

func main() {

	pinNames := []string{"gpio4", "gpio17", "gpio18", "gpio22", "gpio23", "gpio24", "gpio25", "gpio27"}

	var (
		pin hwio.Pin
		err error
	)

	pins := make([]hwio.Pin, 0)

	for _, pinName := range pinNames {

		// get a pin by name. You could also just use the logical pin number, but this is
		// more readable. On BeagleBone, USR0 is an on-board LED.
		if pin, err = hwio.GetPin(pinName); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Set the mode of the pin to output. This will return an error if, for example,
		// we were trying to set an analog input to an output.
		if err = hwio.PinMode(pin, hwio.OUTPUT); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		pins = append(pins, pin)
	}

	for _, pin := range pins {

		for {
			for i := 0; i < 3; i++ {
				fmt.Printf("Blinking %v pin %v \n", pin)
				hwio.DigitalWrite(pin, hwio.HIGH)
				hwio.Delay(500)
				hwio.DigitalWrite(pin, hwio.LOW)
				hwio.Delay(500)
			}
		}
	}
}
