package main

import (
	"fmt"
	"html/template"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"./game"
)

var highLow *game.HighLow

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

	r.ParseForm() // retrieve form guess value
	guess := r.Form.Get("guess")

	if g, err = strconv.Atoi(guess); err == nil { // convert to int
		done, err = highLow.Guess(g)
		if done {
			err = fmt.Errorf(fmt.Sprintf("Congratulations ! You solved the previous game with %d guesses. Try again.", highLow.Guesses))
			highLow = game.NewHighLow()
		}
	} else {
		err = fmt.Errorf("Invalid number")
	}

	t, _ := template.ParseFiles("highlow.html")
	t.Execute(w, &htmlParms{&highLow.CurrentScore, err})
}

func main() {

	rand.Seed(time.Now().UTC().UnixNano())

	server := http.Server{
		Addr: "127.0.0.1:8080",
	}

	http.HandleFunc("/", processHandler)
	server.ListenAndServe()
}
