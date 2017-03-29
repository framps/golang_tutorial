package single

// Samples used in a small go tutorial
//
// Copyright (C) 2017 framp at linux-tips-and-tricks dot de
//
// methods in go

import "fmt"

type motorState int // internal motorstate

const ( // motorstate enum
	motorStateOff motorState = iota
	motorStateOn
)

type motorStateText string

// Textual representation of state enum
var motorStateAsText = map[motorState]string{
	motorStateOff: "Off",
	motorStateOn:  "On",
}

// Car - hiding internals because all var names start with lower case
type Car struct {
	name         string
	color        string
	maxSpeed     int
	currentSpeed int
	state        motorState
}

// NewCar - kind of 'constructor' for car
func NewCar(name, color string, maxSpeed int) *Car {
	return &Car{name, color, maxSpeed, 0, motorStateOff}
}

// String - Formatterof Car for fmt
func (c *Car) String() string {
	return fmt.Sprintf("%s-> color: %s maxSpeed: %d speed: %d state: %s", c.name, c.color, c.maxSpeed, c.currentSpeed, motorStateAsText[c.state])
}

// Accelerate - may be positive or negative
func (c *Car) Accelerate(increment int) {
	fmt.Println(c)
	if increment > 0 {
		if c.currentSpeed+increment < c.maxSpeed {
			c.currentSpeed += increment
			fmt.Printf("%s: Accelerated by %d to %d\n", c.name, increment, c.currentSpeed)
		} else {
			c.currentSpeed = c.maxSpeed
			fmt.Printf("%s: Reached max speed %d\n", c.name, c.currentSpeed)
		}
	} else {
		if c.currentSpeed+increment > 0 {
			c.currentSpeed += increment
			fmt.Printf("%s: Accelerated by %d to %d\n", c.name, increment, c.currentSpeed)
		} else {
			c.currentSpeed = 0
			fmt.Printf("%s: Reached min speed %d\n", c.name, c.currentSpeed)
		}
	}
}

// TurnOn -
func (c *Car) TurnOn() {
	fmt.Printf("%s-> Turning on\n", c.name)
	c.state = motorStateOn
	c.currentSpeed = 0
	fmt.Println(c)
}

// TurnOff -
func (c *Car) TurnOff() {
	fmt.Printf("%s-> Turning off\n", c.name)
	c.state = motorStateOff
	c.currentSpeed = 0
	fmt.Println(c)
}
