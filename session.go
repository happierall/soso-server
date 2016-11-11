package soso

import (
	"fmt"
	"sync"
)

var Sessions = NewSessionRepository()

type Session interface {
	// Id returns a session id
	ID() string
	// Recv reads one text frame from session
	Recv() ([]byte, error)
	// Send sends one text frame to session
	Send(string) error
	// Close closes the session with provided code and reason.
	Close(status uint32, reason string) error
	IsClosed() bool
}

type SessionRepository interface {
	//Push adds session to collection
	Push(session Session, uid string) int
	//Get retries all active sessions for the user
	Get(uid string) []Session
	//Pull removes session object from collection
	Pull(session Session) bool
	//Size returns count of active sessions
	Size(uid string) int
}

func NewSessionRepository() SessionRepository {
	return &SessionRepositoryImpl{
		sessions: make(map[string]string),
		users:    make(map[string][]Session),
	}
}

type SessionRepositoryImpl struct {
	sync.Mutex
	sessions map[string]string
	users    map[string][]Session
}

func (s *SessionRepositoryImpl) Push(session Session, uid string) int {
	s.Lock()
	defer s.Unlock()
	fmt.Println(fmt.Sprintf("Push session %s for user %v", session.ID(), uid))
	sessions, ok := s.users[uid]
	if !ok {
		sessions = make([]Session, 0)
	}
	if _, ok := s.sessions[session.ID()]; !ok {
		s.users[uid] = append(sessions, session)
		s.sessions[session.ID()] = uid
	}
	fmt.Println(fmt.Sprintf("Session %s for user %v pushed, total %v", session.ID(), uid, len(s.users[uid])))
	return len(s.users[uid])
}

func (s *SessionRepositoryImpl) Get(uid string) []Session {
	s.Lock()
	defer s.Unlock()
	sessions, ok := s.users[uid]
	if !ok {
		sessions = make([]Session, 0)
	}
	return sessions
}

func (s *SessionRepositoryImpl) Pull(session Session) bool {
	s.Lock()
	defer s.Unlock()
	uid, ok := s.sessions[session.ID()]
	if !ok {
		return false
	}
	fmt.Println(fmt.Sprintf("Pull session %s for user %v", session.ID(), uid))
	var found int
	for key, value := range s.users[uid] {
		if value.ID() == session.ID() {
			found = key
		}
	}
	s.users[uid] = append(s.users[uid][:found], s.users[uid][found+1:]...)
	delete(s.sessions, session.ID())
	fmt.Println(fmt.Sprintf("Session %s for user %v pulled, total %v", session.ID(), uid, len(s.users[uid])))
	return true
}

func (s *SessionRepositoryImpl) Size(uid string) int {
	s.Lock()
	defer s.Unlock()
	sessions, ok := s.users[uid]
	if !ok {
		return 0
	}
	return len(sessions)
}
