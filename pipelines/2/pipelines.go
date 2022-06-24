// usage of pipelines in golang
//
// See https://blog.golang.org/pipelines for details

package main

import (
	"fmt"
	"sync"
)

func gen(nums ...int) <-chan int {
	out := make(chan int)
	go func() {
		for _, n := range nums {
			fmt.Printf("--- Generator: Generate %d\n", n)
			out <- n
		}
		fmt.Println("--- Generator: Closing channel")
		close(out)
	}()
	return out
}

func sq(in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		for n := range in {
			fmt.Printf("--- Square: Calculate square of %d. result: %d\n", n, n*n)
			out <- n * n
		}
		fmt.Println("--- Square: Closing channel")
		close(out)
	}()
	return out
}

func merge(cs ...<-chan int) <-chan int {
	var wg sync.WaitGroup
	out := make(chan int)

	// Start an output goroutine for each input channel in cs.  output
	// copies values from c to out until c is closed, then calls wg.Done.
	output := func(c <-chan int) {
		for n := range c {
			fmt.Printf("--- Merge: Reading %d\n", n)
			out <- n
		}
		fmt.Println("--- Merge: Done")
		wg.Done()
	}
	fmt.Printf("--- Merge: Add %d to waitgroup\n", len(cs))
	wg.Add(len(cs))
	for _, c := range cs {
		go output(c)
	}

	// Start a goroutine to close out once all the output goroutines are
	// done.  This must start after the wg.Add call.
	go func() {
		fmt.Println("--- Merge: Waiting")
		wg.Wait()
		fmt.Println("--- Merge: Closing")
		close(out)
	}()
	return out
}

func main() {
	in := gen(2, 3)

	// Distribute the sq work across two goroutines that both read from in.
	c1 := sq(in)
	c2 := sq(in)

	// Consume the merged output from c1 and c2.
	for n := range merge(c1, c2) {
		fmt.Println(n) // 4 then 9, or 9 then 4
	}
}
