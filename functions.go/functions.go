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

	if n == 0 {
		return 0, nil
	} else if n == 1 {
		return 1, nil
	} else if n > 10 {
		return 0, fmt.Errorf("number too big: %v - Limit: %v", n, maxFibonacci) // returns an error object
	}

	var ( // scoping - n1 and n1 are locat to the if statement otherwise if not defined beforehand
		err error
		n1  int
		n2  int
	)

	if n2, err = fibonacci(n - 2); err != nil {
		return 0, err
	}

	if n1, err = fibonacci(n - 1); err != nil { // should not occur
		return 0, err
	}
	return n1 + n2, nil
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
