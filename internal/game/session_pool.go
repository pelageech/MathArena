package game

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

var ErrSessionExists = errors.New("session exists")

type ActiveSessionsPool struct {
	mu       sync.Mutex
	sessions map[SessionID]*Session
}

func NewActiveSessionsPool() *ActiveSessionsPool {
	return &ActiveSessionsPool{
		sessions: make(map[SessionID]*Session),
	}
}

func (ap *ActiveSessionsPool) Get(sessionID SessionID, timeNow time.Time) (*Session, error) {
	ap.mu.Lock()
	defer ap.mu.Unlock()
	s, ok := ap.sessions[sessionID]
	if !ok {
		return nil, fmt.Errorf("session %v not found", sessionID)
	}
	if !s.CheckTime(timeNow) {
		delete(ap.sessions, sessionID)
		return nil, fmt.Errorf("%v: %w", sessionID, ErrTimeIsLeft)
	}

	return ap.sessions[sessionID], nil
}

func (ap *ActiveSessionsPool) Put(session *Session) error {
	ap.mu.Lock()
	defer ap.mu.Unlock()
	if _, ok := ap.sessions[session.sessionID]; ok {
		return fmt.Errorf("%v: %w", session.sessionID, ErrSessionExists)
	}
	ap.sessions[session.sessionID] = session
	return nil
}

func (ap *ActiveSessionsPool) Delete(sessionID SessionID) {
	ap.mu.Lock()
	defer ap.mu.Unlock()
	delete(ap.sessions, sessionID)
}
