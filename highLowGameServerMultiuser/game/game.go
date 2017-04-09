// Samples used in a small go tutorial
//
// Copyright (C) 2017 framp at linux-tips-and-tricks dot de
//
// Samples for go templates
//
// See github.com/framps/golang_tutorial for latest code

package game

import (
	"fmt"
	"math/rand"
)

// States -
type States int

const (
	// GameInitialized -
	GameInitialized States = iota
	// GameRunning -
	GameRunning
	// GameFinished -
	GameFinished
)

// CurrentScore -
type CurrentScore struct {
	Low     int
	High    int
	Guesses int
}

// HighLow -
type HighLow struct {
	ActualValue int
	State       States
	CurrentScore
}

const ( // upper and lower bounds
	lowLimit  = 1
	highLimit = 99
)

// NewHighLow - Create new game
func NewHighLow() *HighLow {
	return &HighLow{
		rand.Intn(highLimit-lowLimit+1) + lowLimit,
		GameInitialized,
		CurrentScore{
			Low:     lowLimit - 1,
			High:    highLimit + 1,
			Guesses: 0,
		},
	}
}

// GetState -
func (h *HighLow) GetState() (s States) {
	return h.State
}

// Guess - execute a guess
func (h *HighLow) Guess(guess int) (hit bool, err error) {

	Log.Infof("Guessing")

	if guess <= h.Low || guess >= h.High { // check guess
		err := fmt.Errorf("Number %d is out of bounds", guess)
		return false, err
	}

	h.Guesses++ // another round

	if h.ActualValue == guess { // hit
		h.State = GameFinished
		return true, nil
	}

	switch h.ActualValue > guess { // set lower/upper value accordingly
	case true:
		h.Low = guess
	case false:
		h.High = guess
	}

	h.State = GameRunning

	return false, nil
}
