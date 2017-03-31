package simple

import "fmt"

// Samples used in a small go tutorial
//
// Copyright (C) 2017 framp at linux-tips-and-tricks dot de
//
// interfaces in go

// everyone can speak

type Speaker interface {
	Say() string
}

// everyone eats

type Eater interface {
	Eat() string
}

// this one can speak and eat (multiple interfaces)

type EaterAndSpeaker interface {
	Speaker
	Eater
}

// just a human

type Human struct {
	Name string
}

func (h Human) String() string { // implement stringer interface of fmt to allow to use struct in fmt arguments
	return fmt.Sprintf("%s", h.Name)
}

func (h Human) Say() string {
	return fmt.Sprintf("Hi. My name is %s", h.Name)
}

func (h Human) Eat() string {
	return fmt.Sprintf("Schmatz, knurpsl (by %s)", h.Name)
}

// just a programmer

type Programmer struct {
	Name string
}

func (p Programmer) String() string { // implement stringer interface of fmt to allow to use struct in fmt arguments
	return fmt.Sprintf("%s", p.Name)
}

func (p Programmer) Say() string {
	return fmt.Sprintf("Hi. I'm %s and I like go programming", p)
}

func (p Programmer) Eat() string {
	return fmt.Sprintf("%s (That's me) likes pizza", p)
}

// just a bird

type Bird struct {
	Name string
}

func (b Bird) String() string { // implement stringer interface of fmt to allow to use struct in fmt arguments
	return fmt.Sprintf("%s", b.Name)
}

func (b Bird) Say() string {
	return fmt.Sprintf("Queek - %s - Queek", b)
}

func (b Bird) Eat() string {
	return fmt.Sprintf("Pick by %s", b)
}
