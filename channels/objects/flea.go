package objects

import (
	"math/rand"
	"time"
)

var fleaIDCounter int

const (
	fleaSuspendTime = time.Millisecond * 2000
)

// Flea -
type Flea struct {
	ID             int    // id of flea
	Name           string // readable name
	Location       *Location
	CommandChannel FleaCommandChannel // command channel the flea listens on
}

// NewFlea -
func NewFlea(name string, location *Location) *Flea {
	id := fleaIDCounter
	fleaIDCounter++
	return &Flea{ID: id, Name: name, Location: location}
}

// FleaCommand -
type FleaCommand string

// FleaCommandChannel -
type FleaCommandChannel chan FleaCommand

const (
	// Quiet -
	Quiet FleaCommand = "QUIET"
	// Home -
	Home FleaCommand = "HOME"
)

// YoureFree -
func (f *Flea) YoureFree() {
	go func() {
		for {
			f.Walk()
			time.Sleep(fleaSuspendTime)
		}
	}()
}

// Walk -
func (f *Flea) Walk() {
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
