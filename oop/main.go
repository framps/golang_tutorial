// Samples used in a small go tutorial
//
// Copyright (C) 2017 framp at linux-tips-and-tricks dot de
//
// Samples for go syntax & packages - calculate fibonacci number via recursion

package main

import (
	"fmt"

	"github.com/framps/golang_tutorial/methods/multiple"
)

func useMultiple() {

	// Nasty things with multiple inheritance in go

	elephant := multiple.Mammal{
		GestationPeriod: 10,
		Animal:          multiple.Animal{Name: "Elephant"},
	}

	batman := multiple.Bat{
		Mammal: multiple.Mammal{GestationPeriod: 3,
			Animal: multiple.Animal{Name: "Batman"},
		},
		WingedAnimal: multiple.WingedAnimal{Span: 5,
			Animal: multiple.Animal{Name: "Batman"}, // has to be initialized twice ... otherwise flap prints no name
		},
	}

	animals := make([]multiple.Animal, 0)
	animals = append(animals, elephant.Animal)
	animals = append(animals, batman.Mammal.Animal)

	fmt.Println("--- Everybody eats")
	for _, animal := range animals {
		animal.Eat()
	}

	fmt.Println("--- Methods of bat")
	batman.Mammal.Eat()  // from Animal
	batman.Breathe()     // from Mammal
	batman.Flap()        // from WidgedAnimal
	batman.Ultrasonnic() // from Bat

}

func main() {
	useMultiple()
}
