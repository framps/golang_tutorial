// Samples used in a small go tutorial
//
// Copyright (C) 2017 framp at linux-tips-and-tricks dot de
//
// amples for go syntax & packages - calculate fibonacci number via recursion

package fibonacci

import "fmt"

// MaxFibonacci - maximum value for Fibonacci calculation
const MaxFibonacci = 10 // just calculate fibo numbers until this limit

// Fibonacci - f(n) = f(n-1) + f(n-2))
// returns two values and the second is an indication for any errors
func Fibonacci(n int) (int, error) {

	// use switch instead of if then elif then
	switch {
	case n < 0:
		return 0, fmt.Errorf("ERROR: number too small: %v - Minimum number: 0", n) // returns an error object
	case n < 2:
		return n, nil
	default:
		if n > MaxFibonacci {
			return 0, fmt.Errorf("ERROR: number too big: %v - Maximum number: %v", n, MaxFibonacci) // returns an error object
		}

		n1, _ := Fibonacci(n - 1) // ignore error (n is <= max && >= 0)
		n2, _ := Fibonacci(n - 2)
		return n1 + n2, nil

	}

}
