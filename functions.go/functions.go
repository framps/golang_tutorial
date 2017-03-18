// Samples used in a small go tutorial
//
// Copyright (C) 2017 framp at linux-tips-and-tricks dot de
//
// Samples for go syntax - calculate fibonacci number via recursion

package main

import "fmt"

const maxFibonacci = 10 // just calculate fibo numbers until this limit

// f(n) = f(n-1) + f(n-2))
// returns two values and the second is an indication for any errors
func fibonacci(n int) (int, error) {

	switch {
	case n <= 2:
		return n, nil
	default:
		if n > 10 {
			return 0, fmt.Errorf("number too big: %v - Limit: %v", n, maxFibonacci) // returns an error object
		}

		n1, _ := fibonacci(n - 1)
		n2, _ := fibonacci(n - 2)
		return n1 + n2, nil

	}

}

func main() {

	for i := 0; i < maxFibonacci+3; i++ {
		value, err := fibonacci(i)
		if err != nil {
			fmt.Printf("Fibonacci number of %d cannot be calculated: %v\n", i, err)
		} else {
			fmt.Printf("Fibonacci number of %d is %d\n", i, value)
		}
	}

}
