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
	Phases     []Phase
	state      int
	clockSpeed time.Duration
}

// ProgramTest - Turn every lamp on
var ProgramTest = &Program{
	Phases: []Phase{
		Phase{red, 1},
		Phase{yellow, 1},
		Phase{green, 1},
		Phase{yellow, 1},
	},
	clockSpeed: time.Millisecond * 100,
}

// ProgramWarning - Traffic light is not working, just blink
var ProgramWarning = &Program{
	Phases: []Phase{
		Phase{yellow, 1},
		Phase{off, 1},
	},
	clockSpeed: time.Millisecond * 500,
}

// ProgramNormal - Just the common traffic light
var ProgramNormal = &Program{
	Phases: []Phase{
		Phase{green, 3},
		Phase{yellow, 1},
		Phase{red, 3},
		Phase{redyellow, 1},
	},
	clockSpeed: time.Second * 1,
}
