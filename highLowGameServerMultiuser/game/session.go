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
	"net/http"
	"sync"
	"time"

	uuid "github.com/satori/go.uuid"
)

const (
	maxPlayers      = 10                         // maximum number of sessions
	maxPlayTime     = time.Minute * 5            // session cleanup
	cookieName      = "framps-highLowGameServer" // name of cookie
	cleanUpInterval = time.Minute * 30           // session cleanup interval
)

// Session -
type Session struct {
	ID       string
	game     *HTTPGame
	lastSeen time.Time
}

// NewSession -
func NewSession(id string) *Session {
	return &Session{
		ID:       id,
		game:     NewHTTPGame(id),
		lastSeen: time.Now(),
	}
}

// TimeoutIn -
func (s *Session) TimeoutIn() time.Duration {
	return s.lastSeen.Add(maxPlayTime).Sub(time.Now())
}

// SessionManager -
type SessionManager struct {
	sessions       map[string]*Session
	onlineSessions int
	lock           sync.Mutex // protects session
}

// NewSessionManager -
func NewSessionManager() *SessionManager {
	return &SessionManager{
		sessions:       make(map[string]*Session),
		onlineSessions: 0}
}

// retrieve cookie from client or create a new one on client
// return id of client and flag whether client is new
func retrieveOrCreateCookie(w http.ResponseWriter, r *http.Request) (id string, newUser bool) {

	header, err := r.Cookie(cookieName)

	var user string
	newUser = false

	if err != nil {
		user = fmt.Sprintf("%v", uuid.NewV4())

		http.SetCookie(w, &http.Cookie{
			Name:  cookieName,
			Value: user,
		},
		)

		Log.Infof("New user %s", user)
		newUser = true
	} else {
		user = header.Value
		Log.Infof("Found user %s", user)
	}

	return user, newUser
}

// PlayGame -
func (s *Session) PlayGame(w http.ResponseWriter, r *http.Request, onlinePlayers int) {
	s.game.Play(w, r, onlinePlayers)
}

// OnlineSessions -
func (m *SessionManager) OnlineSessions() *[]*Session {

	result := make([]*Session, 0)

	for _, s := range m.sessions {
		if s != nil {
			result = append(result, s)
		}
	}
	return &result
}

// GetOrCreateSession -
// retrieve existing session or create a new one for request
func (m *SessionManager) GetOrCreateSession(w http.ResponseWriter, r *http.Request) (*Session, error) {

	m.cleanupSessionsImpl()

	user, newUser := retrieveOrCreateCookie(w, r)

	m.lock.Lock()
	defer m.lock.Unlock()

	if newUser || m.sessions[user] == nil {

		if m.onlineSessions >= maxPlayers {
			Log.Warn("Maximum users reached %d", m.onlineSessions)
			err := fmt.Errorf("All game slots used right now :-(")
			return nil, err
		}

		m.sessions[user] = NewSession(user)
		m.onlineSessions++
		Log.Infof("Created game for %s - Online: %d", user, m.onlineSessions)
	}

	m.sessions[user].lastSeen = time.Now()

	return m.sessions[user], nil
}

// CleanupSessions -
// remove sessions
func (m *SessionManager) CleanupSessions() {
	ticker := time.NewTicker(cleanUpInterval)
	go func() {
		for _ = range ticker.C {
			m.cleanupSessionsImpl()
		}
	}()
}

// do the actual sessioncleanup
func (m *SessionManager) cleanupSessionsImpl() {

	m.lock.Lock()
	defer m.lock.Unlock()

	for i, u := range m.sessions {
		if u != nil {
			Log.Infof("%s: %v - %v", (*u).ID, time.Now(), (*u).lastSeen.Add(maxPlayTime))
			if time.Now().After((*u).lastSeen.Add(maxPlayTime)) {
				Log.Infof("Id %s timed out", (*u).ID)
				m.sessions[i] = nil
				m.onlineSessions--
			}
		}
	}
	Log.Infof("Current active sessions: %d", m.onlineSessions)

}
