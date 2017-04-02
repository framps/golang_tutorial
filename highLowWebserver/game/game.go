package game

import (
	"fmt"
	"math/rand"
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
	CurrentScore
}

// constant definitions
const ( // upper and lower bounds
	highLimit = 100 // integer constant
	lowLimit  = 1
)

// NewHighLow - Create new game
func NewHighLow() *HighLow {
	return &HighLow{
		rand.Intn(highLimit-lowLimit+1) + lowLimit,
		CurrentScore{
			Low:     lowLimit,
			High:    highLimit,
			Guesses: 0,
		},
	}
}

// Score - retun current score
func (h *HighLow) Score() (s CurrentScore) {
	return h.Score()
}

// Guess - execute a guess
func (h *HighLow) Guess(guess int) (hit bool, err error) {

	fmt.Printf("Guess: %d\nState: %+v\n", guess, *h)

	if guess < h.Low || guess > h.High { // logical comparisons
		err := fmt.Errorf("Number %d is out of bounds\n", guess)
		fmt.Printf("Error: %v", err)
		return false, err
	}

	h.Guesses++

	if h.ActualValue == guess {
		return true, nil
	}

	switch h.ActualValue > guess {
	case true:
		h.Low = guess + 1
		fmt.Printf("Lower\n")
	case false:
		h.High = guess - 1
		fmt.Printf("Higher\n")
	}

	fmt.Printf("State: %+v\n", *h)

	return false, nil
}
