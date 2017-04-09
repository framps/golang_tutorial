// Samples used in a small go tutorial
//
// Copyright (C) 2017 framp at linux-tips-and-tricks dot de
//
// Samples for go funcitons, channels and mutex

// Based on and enhanced with additional functions and comments

/**
* Toy Go implementation of Dining Philosophers problem
* http://en.wikipedia.org/wiki/Dining_philosophers_problem
* Author: Doug Fritz
* Date: 2011-01-05
**/

package main

import (
	"fmt"
	"math/rand"
	"strings"
	"sync"
	"time"
)

// mutex to synchronize output
var m sync.Mutex

// map true/false to a character to display the availability of a chopstick
var chopStickAvailability = map[bool]string{
	false: "_",
	true:  "F",
}

// the philospher
type philosopher struct {
	name         string       // has a name
	chopstick    chan bool    // a channel to wait for right chopstick
	hasChopstick bool         // flag which tells whether the right chopstick was aquired
	neighbor     *philosopher // the right neighbor
}

// all the philospher who want to dine
var philosophers []*philosopher

// philosopher just tells what happend and what he's doing right now
// use mutex so the he's not interupted
func (phil *philosopher) say(text string, args ...interface{}) {
	m.Lock()
	fmt.Printf("%s", strings.Repeat("-", 60))
	fmt.Printf(text, args...)
	m.Unlock()
}

// displays the current status of chopsticks and philosophers
// use mutex so the status creation will not be interupted
// Format (keep in mind it's a table):
// F Kant F Heid _ Witt _ Lock F Desc F Newt F Hume F Leib
func status() {
	m.Lock()
	for i, p := range philosophers {
		if i != len(philosophers) {
			fmt.Printf("%s %s ", chopStickAvailability[p.hasChopstick], p.name)
			continue
		}
		if i == len(philosophers) {
			fmt.Printf("%s %s", p.name, chopStickAvailability[p.neighbor.hasChopstick])
			continue
		}

	}
	fmt.Printf("\n")
	m.Unlock()
}

// constructor for a new philosopher
func newPhilosopher(name string, neighbor *philosopher) *philosopher {
	phil := &philosopher{name, make(chan bool, 1), false, neighbor}
	phil.chopstick <- true // chopstick is not used, signal to be free
	return phil
}

// philosopher thinks and doesn't eat
func (phil *philosopher) think() {
	phil.say("%v is thinking.\n", phil.name)
	status()
	time.Sleep(time.Duration(rand.Int63n(1e9)))
}

// philosopher aquired two chopsticks and now can eat
func (phil *philosopher) eat() {
	phil.say("%v is eating.\n", phil.name)
	status()
	time.Sleep(time.Duration(rand.Int63n(1e9)))
}

// philosopher tries to aquired two chopsticks left and right to him
func (phil *philosopher) getChopsticks() {

	timeout := make(chan bool, 1) // channle required to get notified about a timeout

	go func() { time.Sleep(1e9); timeout <- true }() // set time out just so the left chopstick will be released if no right chopstick could be aquired
	<-phil.chopstick                                 // wait for left chopstick to be available
	phil.hasChopstick = true                         // aquired left chopstick
	phil.say("%v got his chopstick.\n", phil.name)
	status()

	select {
	case <-phil.neighbor.chopstick: // wait for right chopstick
		phil.neighbor.hasChopstick = true
		phil.say("%v got %v's chopstick.\n", phil.name, phil.neighbor.name)
		phil.say("%v has two chopsticks.\n", phil.name)
		status()
		return
	case <-timeout: // bad luck , release chopstick and try again to eat later
		phil.hasChopstick = false
		phil.chopstick <- true
		phil.think()         // but first think again (hungry)
		phil.getChopsticks() // and try again to get both chopsticks
	}
}

// philosopher finished eating so reaturn both chopsticks
func (phil *philosopher) returnChopsticks() {
	phil.hasChopstick = false
	phil.chopstick <- true // signal left chopstick to be free
	phil.neighbor.hasChopstick = false
	phil.neighbor.chopstick <- true // signal right chopstick to be free
}

// philosopher got two chopsticks and now eats
func (phil *philosopher) dine(announce chan *philosopher) {
	phil.think()            // think first before dining
	phil.getChopsticks()    // aquire chopsticks
	phil.eat()              // now eat
	phil.returnChopsticks() // and return both chopsticks
	announce <- phil        // announce I'm done with dining
}

func main() {
	/*
		names := []string{"Kant", "Heidegger", "Wittgenstein",
			"Locke", "Descartes", "Newton", "Hume", "Leibniz"}
	*/
	names := []string{"Kant", "Heid", "Witt", "Lock", "Desc", "Newt", "Hume", "Leib"}
	philosophers = make([]*philosopher, len(names))

	var phil *philosopher
	for i := len(names) - 1; i >= 0; i-- { // create instances of the philosophers and add them to the gloabl list
		phil = newPhilosopher(names[i], phil)
		philosophers[i] = phil
	}
	philosophers[len(names)-1].neighbor = phil

	fmt.Printf("There are %v philosophers sitting at a table.\n", len(philosophers))
	fmt.Printf("They each have one chopstick, and must borrow from their neighbor to eat.\n")

	announce := make(chan *philosopher) // channel of philosopher to announce he finished eating

	for _, phil := range philosophers { // every philosopher now starts to dine
		go phil.dine(announce)
	}

	for i := 0; i < len(names); i++ { // wait for any philosopher to finish eating
		phil := <-announce
		phil.say("%v is done dining.\n", phil.name) // and report he's done
		status()
	}
}
