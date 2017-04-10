package multiple

// Samples used in a small go tutorial
//
// Copyright (C) 2017 framp at linux-tips-and-tricks dot de
//
// 'multiple inheritance' in go
//
// See github.com/framps/golang_tutorial for latest code

import "fmt"

// Animal -
type Animal struct { // 'base class'
	Name string
}

// Eat -
func (a Animal) Eat() {
	fmt.Printf("%s eating...\n", a.Name)
}

// Mammal -
type Mammal struct { // is an Animal
	GestationPeriod int
	Animal
}

// Breathe -
func (m Mammal) Breathe() {
	fmt.Printf("%s beathing...\n", m.Name)
}

// WingedAnimal -
type WingedAnimal struct { // is a mammal
	Span int
	Animal
}

// Flap -
func (w WingedAnimal) Flap() {
	fmt.Printf("%s flapping...\n", w.Name)
}

// Bat -
type Bat struct { // is a mammal and winged animal
	Mammal
	WingedAnimal
}

// Ultrasonnic -
func (b Bat) Ultrasonnic() {
	fmt.Printf("%s beeping...\n", b.Mammal.Name)
}
