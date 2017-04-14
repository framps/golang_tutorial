package objects

import (
	"bytes"
	"fmt"
	"log"
	"sync"
)

// Area -
type Area struct {
	m         sync.Mutex
	Locations [9]Location
}

// NewArea -
func NewArea() *Area {

	log.Println("Creating new area")

	newArea := &Area{}

	for i := range newArea.Locations {
		newArea.Locations[i] = *NewLocation(newArea)
	}
	return newArea
}

// Return some usefule textual area representation - just display the number of fleas in the area
func (a *Area) String() string {
	var result bytes.Buffer
	result.WriteString("(")
	for _, l := range a.Locations {

		for _, v := range l.Visitors {
			result.WriteString(fmt.Sprintf("%s", string(v.Name[0])))
		}
		result.WriteString(" - ")
	}
	result.WriteString(")")
	return result.String()
}

// PreviousLocation -
func (a *Area) PreviousLocation(l *Location) *Location {
	return &a.Locations[((l.LocationID-1)+len(a.Locations))%len(a.Locations)]
}

// NextLocation -
func (a *Area) NextLocation(l *Location) *Location {
	return &a.Locations[(l.LocationID+1)%len(a.Locations)]
}

// Add - Add a flea to a location in area
func (a *Area) Add(flea *Flea, locationID int) {
	a.Locations[locationID%len(a.Locations)].AddFlea(flea)
}

func (a *Area) move(flea *Flea, from, to *Location) {
	a.m.Lock()
	defer a.m.Unlock()
	to.Visitors[flea.ID] = flea
	delete(from.Visitors, flea.ID)
	flea.Location = to.Visitors[flea.ID].Location
}
