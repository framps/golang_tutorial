// Samples used in a small go tutorial
//
// Copyright (C) 2017 framp at linux-tips-and-tricks dot de
//
// interfaces in go

package main

import (
	"fmt"

	"github.com/framps/golang_tutorial/interfaces/simple"
)

func main() {

	jeff := simple.Human{"Jeff"}
	dumbo := simple.Bird{"Dumbo"}
	charles := simple.Programmer{"Charles"}

	fmt.Println("\n--- Speaker --- ")
	speaker := [...]simple.Speaker{jeff, dumbo, charles}
	for _, s := range speaker {
		fmt.Printf("%s\n", s.Say())
	}

	fmt.Println("\n--- Eater --- ")
	eater := [...]simple.Eater{jeff, dumbo, charles}
	for _, e := range eater {
		fmt.Printf("%s\n", e.Eat())
	}

	fmt.Println("\n--- EaterAndSpeaker --- ")
	eaterAndSpeaker := [...]simple.EaterAndSpeaker{jeff, dumbo, charles}
	for _, e := range eaterAndSpeaker {
		fmt.Printf("%s. It's lunch time: %s\n", e.Say(), e.Eat())
	}

}
