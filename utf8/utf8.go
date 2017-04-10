// Samples used in a small go tutorial
//
// Copyright (C) 2017 framp at linux-tips-and-tricks dot de
//
// Samples for go utf-8 handling
//
// See github.com/framps/golang_tutorial for latest code

package main

import "fmt"

func main() {

	// string constants are utf-8 all the time
	// Note: string variables are most of the time utf-8

	// see https://blog.golang.org/strings

	s := "\xbd\xb2\x3d\xbc\x20\xe2\x8c\x98"

	fmt.Printf("%%s: %s\n", s) // print string
	// %s: ��=� ⌘

	fmt.Printf("%%x: % x\n", s) // print bytes
	// %x: bd b2 3d bc 20 e2 8c 98

	fmt.Printf("%%q: %q\n", s) // print bytes in hex escaped format
	// %q: "\xbd\xb2=\xbc ⌘"

	fmt.Printf("Runes:\n")
	sAsRune := []rune(s)             // convert string into runes (utf-8)
	fmt.Printf("%%s: %q\n", sAsRune) // every rune as character
	// %s: ['�' '�' '=' '�' ' ' '⌘']
	fmt.Printf("%%x: % x\n", sAsRune) // every rune in hex
	// %x: [ fffd  fffd  3d  fffd  20  2318]

}
