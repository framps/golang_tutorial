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
	"html/template"
	"net/http"
	"strconv"
)

// structure passed to template
type htmlParms struct {
	Game          *CurrentScore
	OnlinePlayers int
	Error         error
}

// HTTPGame -
type HTTPGame struct {
	id         string
	highLow    *HighLow           // the game itself
	myTemplate *template.Template // the template used
}

// NewHTTPGame -
func NewHTTPGame(id string) *HTTPGame {
	return &HTTPGame{id: id}
}

// Play -
func (gme *HTTPGame) Play(w http.ResponseWriter, r *http.Request, onlinePlayers int) {

	if gme.highLow == nil {
		gme.highLow = NewHighLow()
		Log.Infof("Creating new Highlow game %+v", gme.highLow)
	}

	var done bool // game finished
	var g int     // guessed number
	var err error // any errors

	Log.Infof("Highlow game %+v", gme.highLow)

	if gme.highLow.GetState() != GameInitialized {

		Log.Infof("Game not initialized")

		r.ParseForm() // retrieve from form the guess value
		guess := r.Form.Get("guess")

		if g, err = strconv.Atoi(guess); err == nil && len(guess) > 0 { // convert to int
			done, err = gme.highLow.Guess(g) // execute game
			if done {
				err = fmt.Errorf(fmt.Sprintf("Congratulations ! You solved the game with %d guesses. Try again.", gme.highLow.Guesses))
				gme.highLow = NewHighLow()
				gme.highLow.State = GameRunning
			}
		} else {
			err = fmt.Errorf("Invalid number")
		}
	} else {
		Log.Infof("Game initialized ... setting to running")
		gme.highLow.State = GameRunning
	}

	if gme.myTemplate == nil { // don't parse the template all the time
		Log.Infof("Creating template")
		gme.myTemplate, err = template.ParseFiles("highlow.html")
		if err != nil {
			panic(err)
		}
	}
	Log.Infof("Displaying template")
	gme.myTemplate.Execute(w, &htmlParms{&gme.highLow.CurrentScore, onlinePlayers, err}) // display template

}
