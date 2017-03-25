package main

import "fmt"

// Samples used in a small go tutorial
//
// Copyright (C) 2017 framp at linux-tips-and-tricks dot de
//
// Samples for go syntax of types

func basicConstTypes() {

	// basic types
	// see https://tour.golang.org/basics/1

	// All constants and variables are typed, either implicit or explicit

	const (

		// ### integeres ###

		// following types are also available as unsigned, i.e.	uint uint8 uint16 uint32 uint64 uintptr
		cint   int   = 42
		cint8  int8  = 0xd
		cint16 int16 = 0377
		cint32 int32 = 0xbeef
		cint64 int64 = 0xaffe

		// ### characters/strings ###

		// everything is UTF-8 in go
		cstring      = "Hello world" // internal string representation is UTF-8
		cbyte   byte = 0xff          // except just a simple binary byte
		crune        = "日本語"         // UTF-8 character/code point

		// ### floats ###

		cfloat32 float32 = 47.11
		cfloat64 float64 = 47.1111e22

		// ### complex ###

		ccomplex64  complex64 = 42i
		ccomplex128 complex64 = 4242i

		// ### bool ###

		cbool bool = true
	)

	// following statement uses constants to eliminate compiler complaing about unused constants
	use(cbool, cbyte, ccomplex128, ccomplex64, cfloat32, cfloat64, cint, cint16, cint32, cfloat64)

	// define two constants in one line
	const boolVariable1, boolVariable2 = true, false

}

// see https://tour.golang.org/basics/1

func basicVariableTypes() {

	// declare and initialze variables. Type is implicit

	var (
		v1 = "Hello world" // string
		v2 = 4711          // int
		v3 = false         // bool
	)

	// Note: Explicit variable declaration has variable name first followed by the type

	var v11, v22, v33 int // all variables are int
	var (
		v4 int
		v5 bool
		v6 string
	)

	// implicit variable declaration and assignment

	var4s := 4711 // defines variable var4 as int and assigns 4711
	// short for
	var var4l int
	var4l = 4711

	// following statement uses variables to eliminate compiler complaing about unused constants
	use(v1, v2, v3, v4, v5, v6, var4l, var4s, v11, v22, v33)

}

func basicDataStructures() {

	// ### arrays ###
	// https://tour.golang.org/moretypes/6

	// in contrast to C arrays are NOT pointers

	var arrayInt10 [10]int // array of 10 ints
	var arrayInt5 [5]int   // array of 5 ints
	// Note: array1Int10 is not compatible with arrayInt5 because of different array size
	// arrayInt10 = arrayInt5 // assignment not possible

	// ### structs ###
	// https://tour.golang.org/moretypes/2

	var cstruct = struct {
		name    string
		address string
	}{"Ronald Grump", "Washington"}
	printStruct("cstruct", cstruct)
	// struct: {name:Ronald Grump address:Washington}

	// stuct type definition
	type Car struct {
		color  string
		wheels int
	}

	var redCar = Car{color: "red", wheels: 4} // long constant form
	printStruct("redCar", redCar)
	// redCar: {color:red wheels:4}

	var blackCar = Car{color: "black"} // long constant form, wheels initialized with 0 (default initialization)
	printStruct("blackCar", blackCar)
	// blackCar: {color:black wheels:0}

	var greenCar = Car{"green", 4} // short constant form, elements have to be in sequence
	printStruct("greenCar", greenCar)
	// greenCar: {color:green wheels:4}

	// ### slices ###
	// dynamic arrays
	// See https://tour.golang.org/moretypes/7
	// see https://blog.golang.org/slices

	// arrays are fixed. Variables arrays are called slices
	// slices have a length and a capacity
	slice1 := []int{1, 2, 3, 4}
	printSlice("slice1", slice1)
	// slice1: [1 2 3 4] - len: 4, cap: 4

	slice2 := make([]int, 3, 10) // create slice of size 3 and capacity 10
	printSlice("slice2", slice2)
	// slice2: [0 0 0] - len: 3, cap: 10

	slice3 := append(slice2, 100, 200, 300) // append 3 elements
	printSlice("slice3: slice2 appended 100,200,300", slice3)
	// slice3: slice2 appended 100,200,300: [0 0 0 100 200 300] - len: 6, cap: 10

	slice4 := append(slice2, slice1...) // slice... converts slice_1 into a variadic parameter for append (variable list of parameters)
	printSlice("slice4: slice1 appended to slice2", slice4)
	// slice4: slice1 appended to slice2: [0 0 0 1 2 3 4] - len: 7, cap: 10

	slice5 := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10} // create slice with 10 elements
	printSlice("slice5", slice5)
	// slice5: [1 2 3 4 5 6 7 8 9 10] - len: 10, cap: 10

	slice6 := slice5[:2] // new slice which has the first two elements and shares elements with slice5
	printSlice("slice6: slice5[:2]", slice6)
	// slice6: slice5[:2]: [1 2] - len: 2, cap: 10

	slice6[0] = 42
	printSlice("slice6: slice6[0] was set to 42", slice6)
	// slice6: slice6[0] was set to 42: [42 2] - len: 2, cap: 10
	printSlice("slice5: slice6[0] was set to 42", slice5)
	// slice5: slice6[0] was set to 42: [42 2 3 4 5 6 7 8 9 10] - len: 10, cap: 10

	slice7 := slice5[5:8] // new slice which starts at 5 and ends at 7! (8-1)
	printSlice("slice7: slice5[5:8]", slice7)
	// slice7: slice5[5:8]: [6 7 8] - len: 3, cap: 5

	slice8 := slice5[3:] // new slice which starts at 3 in slice5
	printSlice("slice8: slice5[3:]", slice8)
	// slice8: slice5[3:]: [4 5 6 7 8 9 10] - len: 7, cap: 7

	// double capacity of a slice
	newSlice8 := make([]int, len(slice8), 2*cap(slice8)) // create new slice with doule capacity
	copy(newSlice8, slice8)                              // copy slice8 into bigger slice
	printSlice("newSlice8: slice8 with double capacity", newSlice8)
	// newSlice8: slice8 with double capacity: [4 5 6 7 8 9 10] - len: 7, cap: 14

	// ### maps ###
	// https://tour.golang.org/moretypes/19

	// keys: string, values: int
	stringMap := map[string]int{
		"eins": 1,
		"zwei": 2,
		"drei": 3}

	printStruct("stringMap", stringMap)
	// stringMap: map[eins:1 zwei:2 drei:3]

	type Address struct {
		street string
		city   string
	}

	// map of users with their address
	userMap := map[string]Address{
		"Chris":    Address{"Bahnhofstrasse", "Stuttgart"},
		"Charly":   Address{"Place de l'etoile", "Paris"},
		"Clarence": Address{"Broadway", "New York"},
	}
	printStruct("userMap", userMap)
	// userMap: map[Clarence:{street:Broadway city:New York} Chris:{street:Bahnhofstrasse city:Stuttgart} Charly:{street:Place de l'etoile city:Paris}]

	// ### pointers ###

	var p1 *int         // pointer to int
	var p2 *[5]int      // pointer to array of 5 ints
	var p3 [8]*int      // array of 8 pointers to int
	var p4 *[10]*[8]int // pointer to 10 arrays of 8 pointers to ints

	var p5 *[10]int // pointer to array of 10 ints
	// p5 = p2					// not possible - pointers have same type but differnet array size

	// following statement uses variables to eliminate compiler complaing about unused constants
	use(arrayInt10, arrayInt5, slice1, slice2, slice3, slice4, slice5, slice6, slice7, cstruct, userMap, p1, p2, p3, p4, p5)

}

// print helpers
func printSlice(name string, slice []int) {
	fmt.Printf("%s: %+v - len: %d, cap: %d\n", name, slice, len(slice), cap(slice))
}

func printStruct(name string, s interface{}) {
	fmt.Printf("%s: %+v\n", name, s)
}

// helper to just use the passed elements to get rid of compiler warning the element is not used
func use(elements ...interface{}) {
	goto nop // yes, go supports go statements also ...
nop:
}

func main() {
	basicConstTypes()
	basicVariableTypes()
	basicDataStructures()
}
