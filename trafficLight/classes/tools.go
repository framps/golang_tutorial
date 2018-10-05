package classes

import (
	"fmt"

	"github.com/framps/golang_tutorial/trafficLight/globals"
)

// Debugging helper
func debugMessage(f string, p ...interface{}) {
	if globals.Debug {
		fmt.Printf(f, p...)
	}
}
