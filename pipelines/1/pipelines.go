// usage of pipelines in golang
//
// See https://blog.golang.org/pipelines for details

package main

import "fmt"

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

func main() {
	// Set up the pipeline.
	c := gen(2, 3)
	out := sq(c)

	// Consume the output.
	fmt.Println(<-out) // 4
	fmt.Println(<-out) // 9
}
