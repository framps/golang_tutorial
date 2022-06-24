package main

// Samples used in a small go tutorial
//
// Copyright (C) 2017 framp at linux-tips-and-tricks dot de
//
// Use composition to build more complex objects
//
// See github.com/framps/golang_tutorial for latest code

import "fmt"

type Headlights struct {
	On bool
}

func (h *Headlights) TurnOn() {
	fmt.Println("Turning lights on")
	h.On = true
}
func (h *Headlights) TurnOff() {
	fmt.Println("Turning lights off")
	h.On = false
}

type Motor struct {
	Horsepower int
	On         bool
}

func (m *Motor) TurnOn() {
	fmt.Println("Starting motor")
	m.On = true
}
func (m *Motor) TurnOff() {
	fmt.Println("Stopping motor")
	m.On = false
}

type Car struct {
	Color     string
	Direction int
	Speed     int
	Motor
	Headlights
}

func (c *Car) Turn(degrees int) {
	fmt.Printf("Turning car %d°\n", degrees)
	c.Direction += degrees
}

func (c *Car) Accelerate(speed int) {
	c.Speed += speed
	fmt.Printf("Accelerating to %d mph\n", c.Speed)
}

func (c *Car) Stop() {
	c.Speed = 0
	fmt.Printf("Stopping\n")
}

func NewCar(color string, power int) *Car {
	// either full blown ctor
	// c:=&Car{color, 0, 0, Motor{power, false}, Headlights{false}}
	// or short form by using defaults
	c := &Car{Color: color}
	fmt.Printf("Created car %s\n", c)
	return c
}

func (c *Car) String() string {
	return fmt.Sprintf("%#v", c)
}

func main() {

	volvo := NewCar("blue", 210)
	// volvo.TurnOn() ambigious
	volvo.Motor.TurnOn()
	fmt.Println(volvo)
	volvo.Headlights.TurnOn()
	fmt.Println(volvo)
	volvo.Accelerate(50)
	fmt.Println(volvo)

	// drive
	for i := 0; i < 360; i += 45 {
		volvo.Turn(45)
		fmt.Println(volvo)
	}

	volvo.Stop()
	fmt.Println(volvo)
	volvo.Headlights.TurnOff()
	fmt.Println(volvo)
	volvo.Motor.TurnOff()
	fmt.Println(volvo)
}
