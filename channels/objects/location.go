package objects

import "sync"

var locationIDCounter int

// Lock -
var Lock sync.Mutex

// Location -
type Location struct {
	LocationID int
	Area       *Area
	Visitors   map[int]*Flea
}

// NewLocation -
func NewLocation(area *Area) *Location {
	r := &Location{
		locationIDCounter,
		area,
		make(map[int]*Flea),
	}
	locationIDCounter++
	return r
}

// AddFlea -
func (l *Location) AddFlea(f *Flea) {
	l.Visitors[f.ID] = f
}

// JumpRight -
func (l *Location) JumpRight(f *Flea) {
	Lock.Lock()
	defer Lock.Unlock()
	n := l.Area.NextLocation(l)
	n.Visitors[f.ID] = f
	delete(l.Visitors, f.ID)
	f.Location = n
}

// JumpLeft -
func (l *Location) JumpLeft(f *Flea) {
	Lock.Lock()
	defer Lock.Unlock()
	p := l.Area.PreviousLocation(l)
	p.Visitors[f.ID] = f
	delete(l.Visitors, f.ID)
	f.Location = p
}
