package main

// Samples used in a small go tutorial
//
// Turns on and off all GPIOs in sequence
//
// Modified sample blink program from https://github.com/mrmorphic/hwio/tree/master/examples
// For more samples see https://github.com/framps/golang_tutorial
//
// Copyright (C) 2017 framp at linux-tips-and-tricks dot depackage main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/mrmorphic/hwio"
)

const (
	sleepTime  = 250 // time to sleep between blink on and off
	blinkCount = 3   // number of blinks per run
)

func main() {

	// catch ctrlc and cleanup
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		hwio.CloseAll()
		os.Exit(1)
	}()

	// list of GPIOs to use
	pins := []string{"gpio4", "gpio17", "gpio18", "gpio22", "gpio23", "gpio24", "gpio25", "gpio27"}

	var (
		ledPin hwio.Pin
		err    error
	)

	// endless loop
	for {

		// loop through all GPIOs
		for _, pin := range pins {

			// get a pin by name. You could also just use the logical pin number, but this is
			// more readable. On BeagleBone, USR0 is an on-board LED.
			if ledPin, err = hwio.GetPin(pin); err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			// Set the mode of the pin to output. This will return an error if, for example,
			// we were trying to set an analog input to an output.
			if err = hwio.PinMode(ledPin, hwio.OUTPUT); err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			for j := 0; j < blinkCount; j++ {
				hwio.DigitalWrite(ledPin, hwio.HIGH)
				hwio.Delay(sleepTime)
				hwio.DigitalWrite(ledPin, hwio.LOW)
				hwio.Delay(sleepTime)
			}
			hwio.ClosePin(ledPin)
		}
	}
}
