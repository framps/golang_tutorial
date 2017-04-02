package main

import (
	"fmt"
	"html/template"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/framps/golang_tutorial/highLowWebserver/game"
)

var highLow *game.HighLow

type htmlParms struct {
	Game  *game.CurrentScore
	Error error
}

func processHandler(w http.ResponseWriter, r *http.Request) {

	if highLow == nil {
		highLow = game.NewHighLow()
	}

	var done bool
	var g int
	var err error
	r.ParseForm()
	guess := r.Form.Get("guess")
	if g, err = strconv.Atoi(guess); err == nil {
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

	http.HandleFunc("/process", processHandler)
	server.ListenAndServe()
}
