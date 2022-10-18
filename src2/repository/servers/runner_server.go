package servers

import (
	"graft/src2/entity"
	"graft/src2/usecase/persister"
	"graft/src2/usecase/runner"
	"log"
	"math/rand"
	"sync"
	"time"
)

// Create /server/interface
// Create /server/entity
type Runner struct {
	id        string
	role      chan entity.Role
	peers     []entity.Peer
	state     entity.State
	timeout   time.Timer
	mu        sync.Mutex
	persister *persister.Service
	runner    *runner.Service
}

func NewRunner(id string, peers []entity.Peer, persister *persister.Service, runner *runner.Service) *Runner {
	ps, _ := persister.LoadState()

	srv := &Runner{
		id:        id,
		peers:     peers,
		state:     *entity.NewState(ps),
		timeout:   *time.NewTimer(350 * time.Millisecond),
		role:      make(chan entity.Role, 1),
		persister: persister,
		runner:    runner,
	}

	srv.role <- entity.Follower
	return srv
}

func (s *Runner) GetState() entity.State {
	return s.state
}

func (s *Runner) Role() chan entity.Role {
	return s.role
}

func (s *Runner) Timeout() <-chan time.Time {
	return s.timeout.C
}

func (s *Runner) saveState() {
	s.persister.SaveState(&s.state.PersistentState)
}

func (s *Runner) Heartbeat() {
	ELECTION_TIMEOUT := 350 // ms
	s.mu.Lock()
	defer s.mu.Unlock()

	rand.Seed(time.Now().UnixNano())
	timeout := (rand.Intn(ELECTION_TIMEOUT/2) + ELECTION_TIMEOUT/2)
	s.timeout.Reset(time.Duration(timeout) * time.Millisecond)
}

func (s *Runner) DowngradeFollower(term uint32) {
	log.Printf("DOWNGRADE TO FOLLOWER TERM: %d\n", s.state.CurrentTerm)
	s.mu.Lock()
	defer s.mu.Unlock()
	s.role <- entity.Follower

	s.state.CurrentTerm = term
	s.state.VotedFor = ""

	s.saveState()
}

func (s *Runner) UpgradeCandidate() {
	log.Printf("UPGRADE TO CANDIDATE TERM: %d\n", s.state.CurrentTerm)
	s.mu.Lock()
	defer s.mu.Unlock()
	s.role <- entity.Candidate

	s.state.CurrentTerm += 1
	s.state.VotedFor = s.id

	s.saveState()
}

func (s *Runner) UpgradeLeader() {
	log.Printf("UPGRADE TO LEADER TERM: %d\n", s.state.CurrentTerm)
	s.role <- entity.Leader
}

func (s *Runner) GrantVote(id string, lastLogIndex uint32, lastLogTerm uint32) bool {
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

func (s *Runner) Broadcast(fn func(peer entity.Peer, rn *Runner)) {
	// state := s.state
	peers := s.peers

	var wg sync.WaitGroup
	for _, peer := range peers {
		wg.Add(1)
		go func(p entity.Peer, w *sync.WaitGroup) {
			defer w.Done()
			fn(p, s)
		}(peer, &wg)
	}
	wg.Wait()
}

func (s *Runner) Start() {
	for {
		switch <-s.role {
		case entity.Follower:
			s.runner.RunFollower(s)
		case entity.Candidate:
			s.runner.RunCandidate(s)
		case entity.Leader:
			s.runner.RunLeader(s)
		}
	}
}
