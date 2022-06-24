// Samples used in a small go tutorial
//
// Copyright (C) 2019 framp at linux-tips-and-tricks dot de
//
// slices used as parameters
//
// See github.com/framps/golang_tutorial for latest code

package main

import "fmt"

func sliceStatusPtr(sliceName string, slice *[]int) {
	if slice != nil {
		fmt.Printf("SlicePtr %s: isNil: %t len: %d cap: %d\n", sliceName, slice == nil, len(*slice), cap(*slice))
	} else {
		fmt.Printf("SlicePtr %s: isNil: %t\n", sliceName, slice == nil)
	}
}

func sliceStatus(sliceName string, slice []int) {
	fmt.Printf("Slice    %s: isNil: %t len: %d cap: %d\n", sliceName, slice == nil, len(slice), cap(slice))
}

func main() {

	var nullSlice []int
	emptySlice := []int{}
	usedSlice := []int{1, 2, 3, 4}

	sliceStatusPtr("nilSlice", nil)
	sliceStatus("nilSlice", nil)
	sliceStatusPtr("nullSlice", &nullSlice)
	sliceStatus("nullSlice", nullSlice)
	sliceStatusPtr("emptySlice", &emptySlice)
	sliceStatus("emptySlice", emptySlice)
	sliceStatusPtr("usedSlice", &usedSlice)
	sliceStatus("usedSlice", usedSlice)

}
