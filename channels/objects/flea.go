package objects

import (
	"math/rand"
)

var fleaIDCounter int

// Flea -
type Flea struct {
	ID             int    // id of flea
	Name           string // readable name
	Location       *Location
	CommandChannel chan FleaCommand // command channel the flea listens on
}

// NewFlea -
func NewFlea(name string, location *Location, commandChannel chan FleaCommand) *Flea {
	id := fleaIDCounter
	fleaIDCounter++
	return &Flea{id, name, location, commandChannel}
}

// FleaCommand -
type FleaCommand string

const (
	// Quiet -
	Quiet FleaCommand = "QUIET"
	// Home -
	Home FleaCommand = "HOME"
)

// YoureFree -
func (f *Flea) YoureFree() {
	r := rand.Intn(2)
	if r == 0 {
		f.Right()
		return
	}
	f.Left()
}

// Right -
func (f *Flea) Right() {
	f.Location.JumpRight(f)
}

// Left -
func (f *Flea) Left() {
	f.Location.JumpLeft(f)
}
