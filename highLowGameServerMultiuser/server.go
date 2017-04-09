// Samples used in a small go tutorial
//
// Copyright (C) 2017 framp at linux-tips-and-tricks dot de
//
// Samples for go templates, http handler, gofunctions and others
//
// See github.com/framps/golang_tutorial for latest code

package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/framps/golang_tutorial/highLowGameServerMultiuser/game"
)

var sessionManager *game.SessionManager

// request handler
func onlineHandler(w http.ResponseWriter, r *http.Request) {

	sessions := sessionManager.OnlineSessions()

	for _, s := range *sessions {
		w.Write([]byte(fmt.Sprintf("Id: %s - timeout: %v\n", (*s).ID, (*s).TimeoutIn())))
	}
}

// request handler
func processHandler(w http.ResponseWriter, r *http.Request) {

	session, err := sessionManager.GetOrCreateSession(w, r)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`No more game sessions available`))
		return
	}

	game.Log.Infof("Playing game for %+v", session)

	session.PlayGame(w, r, len(*sessionManager.OnlineSessions()))

}

func main() {

	game.NewLog()
	sessionManager = game.NewSessionManager()

	go sessionManager.CleanupSessions()

	rand.Seed(time.Now().UTC().UnixNano())

	server := http.Server{
		Addr: ":8080",
	}

	game.Log.Infof("Starting highlow game server on %s", server.Addr)
	http.HandleFunc("/", processHandler)
	// http.HandleFunc("/online", onlineHandler)
	server.ListenAndServe()
}
