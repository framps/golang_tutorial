package classes

// Samples used in a small go tutorial
//
// Copyright (C) 2017,2018 framp at linux-tips-and-tricks dot de
//
// Samples for go - simple trafficlight simulation using go channels and go routines
//
// See github.com/framps/golang_tutorial for latest code

import (
	"fmt"

	"github.com/framps/golang_tutorial/trafficLight/globals"
)

// Debugging helper
func debugMessage(f string, p ...interface{}) {
	if globals.Debug {
		fmt.Printf(f+"\n", p...)
	}
}
