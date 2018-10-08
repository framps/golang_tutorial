package main

import (
	"flag"
	"fmt"
	"time"

	rpio "github.com/stianeikeland/go-rpio"
)

func cleaAll(gpio2bcm []int) {
	// Clear all GPIOs
	for _, p := range gpio2bcm {
		pin := rpio.Pin(p)
		pin.Output()
		pin.Low()
	}
}

func main() {

	// map GPIO numbers to BCM GPIO numbers
	//                 0   1   2   3   4   5   6   7
	gpio2bcm := []int{17, 18, 27, 22, 23, 24, 25, 4}

	err := rpio.Open()
	if err != nil {
		panic(err)
	}
	defer func() {
		cleaAll(gpio2bcm)
		rpio.Close()
	}()

	cleaAll(gpio2bcm)

	bcmpin := flag.Int("gpio", 21, "GPIO to use")
	flag.Parse()

	// bcm2835 pin

	fmt.Printf("Using pin %d\n", *bcmpin)
	pin := rpio.Pin(gpio2bcm[*bcmpin])
	pin.Output() // Output mode

	for {
		for _, p := range gpio2bcm {
			pin := rpio.Pin(p)
			pin.Toggle() // Toggle pin (Low -> High -> Low)
			time.Sleep(time.Second)
		}
	}
}
