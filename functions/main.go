// Samples used in a small go tutorial
//
// Copyright (C) 2017 framp at linux-tips-and-tricks dot de
//
// Samples for go syntax - calculate fibonacci number via recursion

package main

import (
	"fmt"

	"github.com/framps/golang_tutorial/functions/fibonacci"
)

func main() {

	for i := -2; i < fibonacci.MaxFibonacci+3; i++ {
		if value, err := fibonacci.Fibonacci(i); err != nil {
			fmt.Printf("Fibonacci number of %d cannot be calculated: %v\n", i, err)
		} else {
			fmt.Printf("Fibonacci number of %d is %d\n", i, value)
		}
	}

}
