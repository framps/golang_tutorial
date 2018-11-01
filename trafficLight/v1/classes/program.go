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
	repositoryFile     = "programs.json"
	ProgramWarningName = "Warning"
	ProgramTestName    = "Test"
)

// ProgramRepository -
type ProgramRepository map[string]*Program

// Phase consists of light and number of ticks to flash the lights
type Phase struct {
	Light int `json:"light"`
	Ticks int `json:"ticks"`
}

// Program - Has phases and a state (active phase)
type Program struct {
	Name       string        `json:"name"`
	Phases     []Phase       `json:"phases"`
	state      int           `json:"state"`
	clockSpeed time.Duration `json:"clock_speed"`
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
	clockSpeed: time.Millisecond * 50,
}

// ProgramWarning - Traffic light is not working, just blink
var ProgramWarning = &Program{
	Name: ProgramWarningName,
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
		fmt.Printf("Using custom programs defined in %s\n", repositoryFile)
		prc[ProgramTest.Name] = ProgramTest
		prc[ProgramWarning.Name] = ProgramWarning
		return prc
	}
	fmt.Printf("%s\n", err.Error())
	fmt.Printf("Using default programs\n")
	return prd
}

// Load -
func (pr *ProgramRepository) Load() (ProgramRepository, error) {
	file, e := ioutil.ReadFile(repositoryFile)
	if e != nil { // error
		fmt.Printf("%s read error: %v\n", repositoryFile, e)
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

	b, e := json.MarshalIndent(pr, "", "   ")
	if e != nil {
		fmt.Printf("JSON marshal error: %v\n", e)
		return e
	}

	e = ioutil.WriteFile(repositoryFile, b, 0644)
	if e != nil { // error
		fmt.Printf("%s write error: %v\n", repositoryFile, e)
		return e
	}
	fmt.Printf("Saving %s\n", repositoryFile)
	return nil
}
