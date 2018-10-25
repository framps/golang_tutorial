package main

import (
	"flag"
	"time"

	"github.com/framps/golang_tutorial/trafficLight/classes"
)

func main() {

	duration := flag.Duration("sleep", time.Millisecond*100, "LED blink sleep between LEDs")
	flag.Parse()

	lc := classes.NewLEDController()
	lc.Enable()
	lc.Open()
	lc.ClearAll()

	const max = 8
	for i := 0; i < max; i++ {
		ip1 := (i + 1) % max
		lc.On(i)
		time.Sleep(*duration)
		lc.Off(i)
		lc.On(ip1)
		time.Sleep(*duration)
	}
	lc.Close()
}
