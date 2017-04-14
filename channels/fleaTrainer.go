// Samples used in a small go tutorial
//
// Copyright (C) 2017 framp at linux-tips-and-tricks dot de
//
// Samples for go channels, functions and mute
//
// See github.com/framps/golang_tutorial for latest code

// @@@@@@@@@@@@@@@@@@@@@@@@@@
// @@@ under construction @@@
// @@@@@@@@@@@@@@@@@@@@@@@@@@

package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"time"

	"github.com/framps/golang_tutorial/channels/objects"
)

var (
	fleaNames = [...]string{"Hugo", "Karl", "Emil"}
)

func startStatusReporter(area *objects.Area) {
	timer := time.NewTimer(time.Second)
	go func(t *time.Timer) {
		for {
			<-t.C
			fmt.Println(area)
			t = time.NewTimer(time.Second)
		}
	}(timer)
}

func main() {

	rand.Seed(time.Now().UTC().UnixNano())
	area := objects.NewArea()

	numberOfFleas := len(fleaNames)

	commandChan := make(chan objects.FleaCommand, numberOfFleas)
	fleas := make([]*objects.Flea, numberOfFleas)

	// Add some fleas

	for i := 0; i < numberOfFleas; i++ {
		fleas[i] = objects.NewFlea(fleaNames[i], &area.Locations[0], commandChan)
		area.Locations[0].AddFlea(fleas[i])
		fmt.Printf("Adding flea %s to location %d\n", fleas[i].Name, i)
	}

	startStatusReporter(area)

	go func() {
		for {
			for i := 0; i < numberOfFleas; i++ {
				fleas[i].YoureFree()
				time.Sleep(time.Millisecond * 500)
			}
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	for {
		<-c
		os.Exit(0)
	}

}
