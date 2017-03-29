// Samples used in a small go tutorial
//
// Copyright (C) 2017 framp at linux-tips-and-tricks dot de
//
// methods in go

package main

import "github.com/framps/golang_tutorial/methods/single"

const (
	turnOn = iota
	turnOff
	accelerate
)

// helper which executes actions on all carrs
func execute(cars []*single.Car, function int, increment ...int) {

	// excute command for all cars
	for _, car := range cars {
		switch function {
		case turnOn:
			car.TurnOn()
		case turnOff:
			car.TurnOff()
		case accelerate:
			car.Accelerate(increment[0])
		}
	}
}

func main() {
	porsche := single.NewCar("Porsche", "Black", 220)
	käfer := single.NewCar("Käfer", "Blue", 80)
	ferrari := single.NewCar("Ferrari", "Red", 350)

	cars := []*single.Car{porsche, käfer, ferrari}

	execute(cars, turnOn)         // turn on car
	execute(cars, accelerate, 50) // accelerate them
	execute(cars, accelerate, 100)
	execute(cars, accelerate, 100)
	execute(cars, accelerate, 100)
	execute(cars, accelerate, -100) // get slower
	execute(cars, accelerate, -100)
	execute(cars, accelerate, -100)
	execute(cars, accelerate, -100)
	execute(cars, turnOff) // turn off car

}
