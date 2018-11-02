package classes

// Samples used in a small go tutorial
//
// Copyright (C) 2017,2018 framp at linux-tips-and-tricks dot de
//
// Samples for go - simple trafficlight simulation using go channels and go routines
//
// See github.com/framps/golang_tutorial for latest code

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

// SaveJSON -
func SaveJSON(object interface{}, f string) error {
	b, e := json.MarshalIndent(object, "", "   ")
	if e != nil {
		fmt.Printf("JSON marshal error: %v\n", e)
		return e
	}
	e = ioutil.WriteFile(f, b, 0644)
	if e != nil { // error
		fmt.Printf("%s write error: %v\n", f, e)
		return e
	}
	fmt.Printf("Saving %s\n", f)
	return nil
}
