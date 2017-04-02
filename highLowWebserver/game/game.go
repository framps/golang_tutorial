package game

import (
	"fmt"
	"math/rand"
)

// States -
type States int

const (
	gameStarted States = iota
	gameRunning
	gameFinished
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

// constant definitions
const ( // upper and lower bounds
	highLimit = 99 // integer constant
	lowLimit  = 1
)

// NewHighLow - Create new game
func NewHighLow() *HighLow {
	return &HighLow{
		rand.Intn(highLimit-lowLimit+1) + lowLimit,
		gameStarted,
		CurrentScore{
			Low:     lowLimit - 1,
			High:    highLimit + 1,
			Guesses: 0,
		},
	}
}

// States -
func (h *HighLow) getState() (s States) {
	return h.State
}

// Guess - execute a guess
func (h *HighLow) Guess(guess int) (hit bool, err error) {

	if guess <= h.Low || guess >= h.High { // logical comparisons
		err := fmt.Errorf("Number %d is out of bounds\n", guess)
		return false, err
	}

	h.Guesses++

	if h.ActualValue == guess {
		h.State = gameFinished
		return true, nil
	}

	switch h.ActualValue > guess {
	case true:
		h.Low = guess
	case false:
		h.High = guess
	}

	h.State = gameRunning

	return false, nil
}
