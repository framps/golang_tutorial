package classes

// Samples used in a small go tutorial
//
// Copyright (C) 2017,2018 framp at linux-tips-and-tricks dot de
//
// Samples for go - simple trafficlight simulation using go channels and go routines
//
// See github.com/framps/golang_tutorial for latest code

import "time"

// lamp colors
const (
	off = iota
	green
	yellow
	red
	redyellow
)

// Phase consists of light and number of ticks to flash the lights
type Phase struct {
	Lights int
	Ticks  int
}

// Program - Has phases and a state (active phase)
type Program struct {
	Name       string
	Phases     []Phase
	state      int
	clockSpeed time.Duration
}

// ProgramTest - Turn every lamp on
var ProgramTest = &Program{
	Name: "Test",
	Phases: []Phase{
		Phase{red, 1},
		Phase{yellow, 1},
		Phase{green, 1},
		Phase{yellow, 1},
	},
	clockSpeed: time.Millisecond * 50,
}

// ProgramWarning - Traffic light is not working, just blink
var ProgramWarning = &Program{
	Name: "Warning",
	Phases: []Phase{
		Phase{yellow, 1},
		Phase{off, 1},
	},
	clockSpeed: time.Millisecond * 500,
}

// ProgramNormal1 - Just the common traffic light
var ProgramNormal1 = &Program{
	Name: "Normal1",
	Phases: []Phase{
		Phase{green, 3},
		Phase{yellow, 2},
		Phase{red, 4},
		Phase{redyellow, 1},
	},
	clockSpeed: time.Second * 1,
}

// ProgramNormal2 - Just the common traffic light
var ProgramNormal2 = &Program{
	Name: "Normal2",
	Phases: []Phase{
		Phase{green, 4},
		Phase{yellow, 1},
		Phase{red, 4},
		Phase{redyellow, 1},
	},
	clockSpeed: time.Second * 1,
}

// ProgramNormal3 - Just the common traffic light
var ProgramNormal3 = &Program{
	Name: "Normal3",
	Phases: []Phase{
		Phase{green, 4},
		Phase{yellow, 1},
		Phase{yellow, 1},

		Phase{red, 4},
		Phase{red, 1},
		Phase{redyellow, 1},
	},
	clockSpeed: time.Second * 1,
}

// ProgramNormal4 - Just the common traffic light
var ProgramNormal4 = &Program{
	Name: "Normal3",
	Phases: []Phase{
		Phase{green, 4},
		Phase{yellow, 1},
		Phase{red, 1},

		Phase{red, 4},
		Phase{red, 1},
		Phase{redyellow, 1},
	},
	clockSpeed: time.Second * 1,
}
