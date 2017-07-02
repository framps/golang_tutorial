// Samples used in a small go tutorial
//
// Copyright (C) 2017 framp at linux-tips-and-tricks dot de
//
// Samples for go syntax & packages - calculate fibonacci number via recursion

package main

import (
	"fmt"

	function "github.com/framps/golang_tutorial/functions/fibonacci"
)

const max = 11

func main() {

	// calculate fibonacci numbers, include some invalid arguments for call
	for i := -1; i <= function.DefaultMax+1; i++ {
		if value, err := function.Fibonacci(i, max); err != nil {
			fmt.Printf("Error calculating Fibonacci number of %d: %v\n", i, err)
		} else {
			fmt.Printf("Fibonacci number of %d is %d\n", i, value)
		}
	}

}
