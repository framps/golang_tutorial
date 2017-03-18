// Samples used in a small go tutorial
//
// Copyright (C) 2017 framp at linux-tips-and-tricks dot de
//
// Samples for go syntax - calculate fibonacci number via recursion

package fibonacci

import "fmt"

// MaxFibonacci - maximum value for Fibonacci calculation
const MaxFibonacci = 10 // just calculate fibo numbers until this limit

// Fibonacci - f(n) = f(n-1) + f(n-2))
// returns two values and the second is an indication for any errors
func Fibonacci(n int) (int, error) {

	switch {
	case n < 2:
		return n, nil
	default:
		if n > MaxFibonacci {
			return 0, fmt.Errorf("number too big: %v - Limit: %v", n, MaxFibonacci) // returns an error object
		}

		n1, _ := Fibonacci(n - 1)
		n2, _ := Fibonacci(n - 2)
		return n1 + n2, nil

	}

}
