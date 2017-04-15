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
	fleaNames         = [...]string{"Hugo", "Karl", "Emil"}
	numberOfFleas     = len(fleaNames)
	numberOfLocations = 9
)

func startFleas(fleas []*objects.Flea) {
	for i := 0; i < numberOfFleas; i++ {
		fleas[i].YoureFree()
	}
}

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
	area := objects.NewArea(numberOfLocations, numberOfFleas)

	fleas := make([]*objects.Flea, numberOfFleas)
	pubSubCtrl := objects.PubSubController{}

	// Add some fleas

	for i := 0; i < numberOfFleas; i++ {
		fleas[i] = objects.NewFlea(fleaNames[i], &area.Locations[0])
		area.Locations[0].AddFlea(fleas[i])
		pubSubCtrl.Register("Listen", fleas[i].CommandChannel)
		fmt.Printf("Adding flea %s to location %d\n", fleas[i].Name, i)
	}

	startStatusReporter(area)

	startFleas(fleas)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	for {
		<-c
		os.Exit(0)
	}

}
