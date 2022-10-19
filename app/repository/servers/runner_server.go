package servers

import (
	"graft/app/entity"
	"graft/app/rpc"
	"graft/app/usecase/persister"
	"graft/app/usecase/runner"
	"log"
	"math"
	"sync"
)

const ELECTION_TIMEOUT = 350 // ms
type Runner struct {
	id            string
	role          Role
	peers         []entity.Peer
	clusterLeader string
	state         entity.State
	timeout       *entity.Timeout
	ticker        *entity.Ticker
	persister     *persister.Service
	mu            sync.Mutex
}

type Role struct {
	value  entity.Role
	signal chan struct{}
}

func NewRunner(id string, peers []entity.Peer, persister *persister.Service, timeout *entity.Timeout, ticker *entity.Ticker) *Runner {

	ps, _ := persister.LoadState()

	srv := &Runner{
		id:        id,
		peers:     peers,
		state:     *entity.NewState(ps),
		timeout:   timeout,
		ticker:    ticker,
		role:      Role{value: entity.Follower, signal: make(chan struct{}, 1)},
		persister: persister,
	}

	return srv
}

func (s *Runner) GetState() entity.State {
	return s.state
}

func (s *Runner) GetQuorum() int {
	totalNodes := len(s.peers) + 1 // add self
	return int(math.Ceil(float64(totalNodes) / 2.0))
}

func (s *Runner) setClusterLeader(leaderId string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.clusterLeader != leaderId {
		log.Printf("FOLLOWING CLUSTER LEADER: %s\n", leaderId)
		s.clusterLeader = leaderId
	}
}

func (s *Runner) saveState() {
	s.persister.SaveState(&s.state.PersistentState)
}

func (s *Runner) resetTimeout() {
	s.timeout.RReset()
	s.ticker.Stop()
}

func (s *Runner) Heartbeat() {
	s.resetTimeout()
}

func (s *Runner) DowngradeFollower(term uint32, leaderId string) {
	s.setClusterLeader(leaderId)
	log.Printf("DOWNGRADE TO FOLLOWER TERM: %d\n", term)
	s.resetTimeout()

	s.mu.Lock()
	s.role.value = entity.Follower
	s.state.CurrentTerm = term
	s.state.VotedFor = ""
	s.mu.Unlock()

	s.saveState()
	s.role.signal <- struct{}{}
}

func (s *Runner) IncrementTerm() {
	if s.role.value == entity.Candidate {
		log.Printf("INCREMENT CANDIDATE TERM: %d\n", s.state.CurrentTerm+1)
		s.resetTimeout()

		s.mu.Lock()
		s.state.CurrentTerm += 1
		s.state.VotedFor = s.id
		s.mu.Unlock()

		s.saveState()
	}
}

func (s *Runner) UpgradeCandidate() {
	if s.role.value == entity.Follower {
		log.Println("UPGRADE TO CANDIDATE")

		s.mu.Lock()
		s.role.value = entity.Candidate
		s.mu.Unlock()

		s.saveState()
		s.role.signal <- struct{}{}
	}
}

func (s *Runner) UpgradeLeader() {
	if s.role.value == entity.Candidate {
		log.Printf("UPGRADE TO LEADER TERM: %d\n", s.state.CurrentTerm)
		s.ticker.Start()
		s.role.value = entity.Leader
		s.role.signal <- struct{}{}
	}
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

func (s *Runner) RequestVoteInput() *rpc.RequestVoteInput {
	state := s.state
	return &rpc.RequestVoteInput{
		CandidateId:  s.id,
		Term:         state.CurrentTerm,
		LastLogIndex: state.LastLogIndex(),
		LastLogTerm:  state.LastLogTerm(),
	}
}

func (s *Runner) AppendEntriesInput() *rpc.AppendEntriesInput {
	state := s.state
	return &rpc.AppendEntriesInput{
		LeaderId: s.id,
		Term:     state.CurrentTerm,
	}
}

func (s *Runner) Broadcast(fn func(peer entity.Peer)) {
	peers := s.peers

	var wg sync.WaitGroup
	for _, peer := range peers {
		wg.Add(1)
		go func(p entity.Peer, w *sync.WaitGroup) {
			defer w.Done()
			fn(p)
		}(peer, &wg)
	}
	wg.Wait()
}

func (s *Runner) Start(service *runner.Service) {
	log.Println("START RUNNER SERVER")

	s.role.signal <- struct{}{}

	for range s.role.signal {
		switch s.role.value {
		case entity.Follower:
			go service.RunFollower(s)
		case entity.Candidate:
			go service.RunCandidate(s)
		case entity.Leader:
			go service.RunLeader(s)
		}

	}
}
