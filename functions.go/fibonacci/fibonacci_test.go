//
// Copyright (C) 2017 framp at linux-tips-and-tricks dot de
//
// Samples for go syntax - calculate fibonacci number via recursion

package fibonacci

import (
	"testing"
)

// table drivern test - usual test pattern in go
var fibonacciTests = []struct {
	input    int // function input
	expected int // expected result
}{
	{1, 1},
	{2, 1},
	{3, 2},
	{4, 3},
	{5, 5},
	{6, 8},
	{7, 13},
	{8, 20}, // should be 21 instead of 20, used to force e test error to show up
}

// TestFibonacci - run it with go test ./fibonacci
func TestFibonacci(t *testing.T) {
	for _, tt := range fibonacciTests {
		t.Logf("Calculating Fibonacci number for %d", tt.input)
		actual, _ := Fibonacci(tt.input)
		if actual != tt.expected {
			t.Errorf("Fibonacci(%d): expected %d, actual %d", tt.input, tt.expected, actual)
		}
	}
}
