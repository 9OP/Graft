package entity

import (
	"log"
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

type Peer struct {
	id   string
	host string
	port string
}

type State struct {
	commitIndex uint32
	lastApplied uint32
	nextIndex   []string // leader only
	matchIndex  []string // leader only
}

func (s *Server) Heartbeat() {
	ELECTION_TIMEOUT := 350 // ms
	s.mu.Lock()
	defer s.mu.Unlock()

	rand.Seed(time.Now().UnixNano())
	timeout := (rand.Intn(ELECTION_TIMEOUT/2) + ELECTION_TIMEOUT/2)
	s.timeout.Reset(time.Duration(timeout) * time.Millisecond)
}

func NewServer() *Server {
	srv := &Server{
		timeout: *time.NewTimer(350 * time.Millisecond),
		role:    make(chan Role, 1),
	}
	srv.role <- Follower
	return srv
}

func (s *Server) Start() {
	log.Println("START SERVER")

	service := NewService()

	for {
		switch <-s.role {
		case Leader:
			service.RunLeader(s)
		case Candidate:
			service.RunCandidate(s)
		case Follower:
			service.RunFollower(s)
		}
	}
}

type UseCase interface {
	RunFollower(follower IFollower)
	RunCandidate(candidate ICandidate)
	RunLeader(leader ILeader)
}

type Service struct{}

func NewService() *Service {
	return &Service{}
}

func (s *Service) RunFollower(follower IFollower)    {}
func (s *Service) RunCandidate(candidate ICandidate) {}
func (s *Service) RunLeader(leader ILeader)          {}
