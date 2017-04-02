package main

import (
	"fmt"
	"html/template"
	"math/rand"
	"net/http"
	"os"
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
	var done bool
	var g int
	var err error
	r.ParseForm()
	guess := r.Form.Get("guess")
	if g, err = strconv.Atoi(guess); err == nil {
		done, err = highLow.Guess(g)
		if done {
			os.Exit(0)
		}
	}
	t, _ := template.ParseFiles("highlow.html")
	t.Execute(w, &htmlParms{&highLow.CurrentScore, err})
}

func highlowHandler(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("highlow.html")
	err := fmt.Errorf("")
	t.Execute(w, &htmlParms{&highLow.CurrentScore, err})
}

func main() {

	rand.Seed(time.Now().UTC().UnixNano())
	highLow = game.NewHighLow()

	server := http.Server{
		Addr: "127.0.0.1:8080",
	}

	http.HandleFunc("/process", processHandler)
	http.HandleFunc("/", highlowHandler)
	server.ListenAndServe()
}
