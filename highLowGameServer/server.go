// Samples used in a small go tutorial
//
// Copyright (C) 2017 framp at linux-tips-and-tricks dot de
//
// Samples for go templates
//
// See github.com/framps/golang_tutorial for latest code

package main

import (
	"fmt"
	"html/template"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/framps/golang_tutorial/highLowGameServer/game"
)

var highLow *game.HighLow // the game itself

// structure passed to template
type htmlParms struct {
	Game  *game.CurrentScore
	Error error
}

// request handler
func processHandler(w http.ResponseWriter, r *http.Request) {

	if highLow == nil {
		highLow = game.NewHighLow()
	}

	var done bool // game finished
	var g int     // guessed number
	var err error // any errors

	if highLow.GetState() != game.GameInitialized {

		r.ParseForm() // retrieve form guess value
		guess := r.Form.Get("guess")

		if g, err = strconv.Atoi(guess); err == nil && len(guess) > 0 { // convert to int
			done, err = highLow.Guess(g)
			if done {
				err = fmt.Errorf(fmt.Sprintf("Congratulations ! You solved the game with %d guesses. Try again.", highLow.Guesses))
				highLow = game.NewHighLow()
				highLow.State = game.GameRunning
			}
		} else {
			err = fmt.Errorf("Invalid number")
		}
	} else {
		highLow.State = game.GameRunning
	}

	t, _ := template.ParseFiles("highlow.html")
	t.Execute(w, &htmlParms{&highLow.CurrentScore, err})
}

func main() {

	rand.Seed(time.Now().UTC().UnixNano())

	server := http.Server{
		Addr: "127.0.0.1:8080",
	}

	fmt.Printf("Starting highlow game server on %s\n", server.Addr)
	http.HandleFunc("/", processHandler)
	server.ListenAndServe()
}
