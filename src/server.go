package main

import (
	"math/rand"
	"sync"
	"time"
)

// Create /server/interface
// Create /server/entity
type Server struct {
	id      string
	role    chan Role
	peers   []Peer
	state   State
	timeout time.Timer
	mu      sync.Mutex
	// heartbeatTicker *time.Ticker // define in state or in runLeader()
}

func (s *Server) Heartbeat() {
	ELECTION_TIMEOUT := 350 // ms
	s.mu.Lock()
	defer s.mu.Unlock()

	rand.Seed(time.Now().UnixNano())
	timeout := (rand.Intn(ELECTION_TIMEOUT/2) + ELECTION_TIMEOUT/2)
	s.timeout.Reset(time.Duration(timeout) * time.Millisecond)
}

func (s *Server) GetState() State {
	return s.state
}

func NewServer() *Server {
	srv := &Server{
		timeout: *time.NewTimer(350 * time.Millisecond),
		role:    make(chan Role, 1),
	}
	srv.role <- Follower
	return srv
}

// Move elswhere
// func (s *Server) Start() {
// 	log.Println("START SERVER")

// 	service := NewService()

// 	for {
// 		switch <-s.role {
// 		case Leader:
// 			service.RunLeader(s)
// 		case Candidate:
// 			service.RunCandidate(s)
// 		case Follower:
// 			service.RunFollower(s)
// 		}
// 	}
// }
