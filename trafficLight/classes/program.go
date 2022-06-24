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
	"time"
)

// lamp colors
const (
	off = iota
	green
	yellow
	red
	redyellow
)

const (
	programFile        = "programs.json"
	ProgramWarningName = "Warning"
	ProgramTestName    = "Test"
)

// ProgramRepository -
type ProgramRepository map[string]*Program

// Phase consists of light and number of ticks to flash the lights
type Phase struct {
	Lights int `json:"lights"`
	Ticks  int `json:"ticks"`
}

// Program - Has phases and a state (active phase)
type Program struct {
	Name       string  `json:"name"`
	Phases     []Phase `json:"phases"`
	state      int
	ClockSpeed time.Duration `json:"clock_speed"`
}

// ProgramTest - Turn every lamp on
var ProgramTest = &Program{
	Name: ProgramTestName,
	Phases: []Phase{
		Phase{red, 1},
		Phase{yellow, 1},

		Phase{green, 1},
		Phase{yellow, 1},
	},
	ClockSpeed: time.Millisecond * 50,
}

// ProgramWarning - Traffic light is not working, just blink
var ProgramWarning = &Program{
	Name: ProgramWarningName,
	Phases: []Phase{
		Phase{yellow, 1},
		Phase{off, 1},
	},
	ClockSpeed: time.Millisecond * 500,
}

var ProgramNormal1 = &Program{
	Name: "Normal1",
	Phases: []Phase{
		Phase{green, 3},
		Phase{yellow, 1},
		Phase{red, 1},
		Phase{red, 1},
		Phase{red, 3},

		Phase{red, 3},
		Phase{red, 1},
		Phase{red, 1},
		Phase{redyellow, 1},
		Phase{green, 3},
	},
	ClockSpeed: time.Second * 1,
}

// ProgramNormal2 - Just the common traffic light
var ProgramNormal2 = &Program{
	Name: "Normal2",
	Phases: []Phase{
		Phase{green, 3},
		Phase{yellow, 1},
		Phase{red, 1},
		Phase{red, 3},

		Phase{red, 3},
		Phase{red, 1},
		Phase{redyellow, 1},
		Phase{green, 3},
	},
	ClockSpeed: time.Second * 1,
}

// ProgramNormal3 - Just the common traffic light
var ProgramNormal3 = &Program{
	Name: "Normal3",
	Phases: []Phase{
		Phase{green, 3},
		Phase{yellow, 1},
		Phase{yellow, 1},

		Phase{red, 3},
		Phase{red, 1},
		Phase{redyellow, 1},
	},
	ClockSpeed: time.Second * 1,
}

func NewProgramRepository() ProgramRepository {

	// Default programs
	prd := ProgramRepository{
		ProgramNormal1.Name: ProgramNormal1,
		ProgramNormal2.Name: ProgramNormal2,
		ProgramNormal3.Name: ProgramNormal3,
	}

	prc, err := prd.Load()
	if err == nil {
		// Add Test and Warning programs
		fmt.Printf("Using custom traffic programs defined in %s\n", programFile)
		prc[ProgramTest.Name] = ProgramTest
		prc[ProgramWarning.Name] = ProgramWarning
		return prc
	} else {
		fmt.Printf("Using default traffic programs\n")
		return prd
	}
}

// Load -
func (pr *ProgramRepository) Load() (ProgramRepository, error) {
	file, e := ioutil.ReadFile(programFile)
	if e != nil { // error
		return nil, e
	}
	var repository ProgramRepository
	e = json.Unmarshal(file, &repository)
	if e != nil {
		fmt.Printf("JSON parse error: %v\n", e)
		return nil, e
	}
	return repository, nil
}

// Save -
func (pr *ProgramRepository) Save() error {
	return SaveJSON(pr, programFile)
}
