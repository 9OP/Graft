package main

import (
	"graft/src2/entity"
	"graft/src2/usecase/state"
	"log"
	"math/rand"
	"sync"
	"time"
)

// Create /server/interface
// Create /server/entity
type Server struct {
	id           string
	role         chan entity.Role
	peers        []entity.Peer
	state        entity.State
	timeout      time.Timer
	mu           sync.Mutex
	stateService *state.Service
}

func NewServer(id string, peers []entity.Peer, stateService *state.Service) *Server {
	ps, _ := stateService.LoadState()

	srv := &Server{
		id:           id,
		peers:        peers,
		state:        *entity.NewState(ps),
		timeout:      *time.NewTimer(350 * time.Millisecond),
		role:         make(chan entity.Role, 1),
		stateService: stateService,
	}

	srv.role <- entity.Follower
	return srv
}

func (s *Server) GetState() entity.State {
	return s.state
}

func (s *Server) saveState() {
	s.stateService.SaveState(&s.state.PersistentState)
}

func (s *Server) Heartbeat() {
	ELECTION_TIMEOUT := 350 // ms
	s.mu.Lock()
	defer s.mu.Unlock()

	rand.Seed(time.Now().UnixNano())
	timeout := (rand.Intn(ELECTION_TIMEOUT/2) + ELECTION_TIMEOUT/2)
	s.timeout.Reset(time.Duration(timeout) * time.Millisecond)
}

func (s *Server) DowngradeFollower(term uint32) {
	log.Printf("DOWNGRADE TO FOLLOWER TERM: %d\n", s.state.CurrentTerm)
	s.mu.Lock()
	defer s.mu.Unlock()
	s.role <- entity.Follower

	s.state.CurrentTerm = term
	s.state.VotedFor = ""

	s.saveState()
}

func (s *Server) UpgradeCandidate() {
	log.Printf("UPGRADE TO CANDIDATE TERM: %d\n", s.state.CurrentTerm)
	s.mu.Lock()
	defer s.mu.Unlock()
	s.role <- entity.Candidate

	s.state.CurrentTerm += 1
	s.state.VotedFor = s.id

	s.saveState()
}

func (s *Server) UpgradeLeader() {
	log.Printf("UPGRADE TO LEADER TERM: %d\n", s.state.CurrentTerm)
	s.role <- entity.Leader
}

func (s *Server) GrantVote(id string, lastLogIndex uint32, lastLogTerm uint32) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	currentLogIndex := s.state.LastLogIndex()
	currentLogTerm := s.state.LastLogTerm()

	voteAvailable := s.state.VotedFor == "" || s.state.VotedFor == id
	candidateUpToDate := currentLogTerm <= lastLogTerm && currentLogIndex <= lastLogIndex

	if voteAvailable && candidateUpToDate {
		s.state.VotedFor = id
		return true
	}

	return false
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
