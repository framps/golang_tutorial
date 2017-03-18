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
	input    int  // function input
	expected int  // expected result
	err      bool // should create error
}{
	{-2, 0, false}, // should be true, used to force test error to show up
	{-1, 0, true},
	{1, 1, false},
	{2, 1, false},
	{3, 2, false},
	{4, 3, false},
	{5, 5, false},
	{6, 8, false},
	{7, 13, false},
	{8, 20, false}, // should be 21 instead of 20, used to force test error to show up
	{11, 0, false}, // should be true, used to force test error to show up
}

// TestFibonacci - run it with go test ./fibonacci
func TestFibonacci(t *testing.T) {
	for _, tt := range fibonacciTests {
		t.Logf("Calculating Fibonacci number for %d", tt.input)
		actual, err := Fibonacci(tt.input)
		if tt.err && err == nil {
			t.Errorf("Fibonacci(%d): expected error", tt.input)
			continue
		}
		if !tt.err && err != nil {
			t.Errorf("Fibonacci(%d): expected no error. Received '%v'", tt.input, err)
			continue
		}
		if actual != tt.expected {
			t.Errorf("Fibonacci(%d): expected %d, actual %d", tt.input, tt.expected, actual)
			continue
		}
	}
}
