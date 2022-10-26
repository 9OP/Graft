package servers

import (
	"graft/app/domain/entity"
	"graft/app/rpc"
	"graft/app/usecase/persister"
	"graft/app/usecase/runner"
	"log"
	"math"
	"sync"
)

// Move to entity/server.go or server_state.go
type Runner struct {
	id            string
	role          *entity.Role
	peers         []entity.Peer
	clusterLeader string
	state         *entity.Istate
	timeout       *entity.Timeout
	ticker        *entity.Ticker
	persister     *persister.Service
	mu            sync.Mutex
}

func NewRunner(id string, peers []entity.Peer, persister *persister.Service, timeout *entity.Timeout, ticker *entity.Ticker) *Runner {
	ps, _ := persister.LoadState()
	state := entity.NewState(ps)
	role := entity.NewRole()
	return &Runner{
		id:        id,
		peers:     peers,
		state:     state,
		timeout:   timeout,
		ticker:    ticker,
		role:      role,
		persister: persister,
	}
}

func (s *Runner) GetState() *entity.State {
	return s.state.GetState()
}

func (s *Runner) GetQuorum() int {
	totalNodes := len(s.peers) + 1 // add self
	return int(math.Ceil(float64(totalNodes) / 2.0))
}

func (s *Runner) SetClusterLeader(leaderId string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.clusterLeader != leaderId {
		log.Printf("FOLLOWING CLUSTER LEADER: %s\n", leaderId)
		s.clusterLeader = leaderId
	}
}

func (s *Runner) saveState() {
	state := s.GetState()
	s.persister.SaveState(&state.Persistent)
}

func (s *Runner) resetTimeout() {
	s.timeout.RReset()
	s.ticker.Stop()
}

func (s *Runner) Heartbeat() {
	s.resetTimeout()
}

func (s *Runner) IsLeader() bool {
	return s.role.Is(entity.Leader)
}

func (s *Runner) DeleteLogsFrom(index uint32) {
	defer s.saveState()
	s.state.DeleteLogsFrom(index)
}

func (s *Runner) AppendLogs(entries []string) {
	defer s.saveState()
	s.state.AppendLogs(entries)
}

func (s *Runner) SetCommitIndex(index uint32) {
	s.state.SetCommitIndex(index)
}

func (s *Runner) DowngradeFollower(term uint32, leaderId string) {
	defer s.saveState()
	log.Printf("DOWNGRADE TO FOLLOWER TERM: %d\n", term)
	s.resetTimeout()
	s.SetClusterLeader(leaderId)
	s.state.SetCurrentTerm(term)
	s.state.SetVotedFor("")
	s.role.Shift(entity.Follower)
}

func (s *Runner) IncrementTerm() {
	if s.role.Is(entity.Candidate) {
		defer s.saveState()
		state := s.GetState()
		newTerm := state.CurrentTerm + 1
		log.Printf("INCREMENT CANDIDATE TERM: %d\n", newTerm)
		s.resetTimeout()
		s.state.SetCurrentTerm(newTerm)
		s.state.SetVotedFor(s.id)
	}
}

func (s *Runner) UpgradeCandidate() {
	if s.role.Is(entity.Follower) {
		log.Println("UPGRADE TO CANDIDATE")
		s.role.Shift(entity.Candidate)
	}
}

func (s *Runner) UpgradeLeader() {
	if s.role.Is(entity.Candidate) {
		state := s.GetState()
		log.Printf("UPGRADE TO LEADER TERM: %d\n", state.CurrentTerm)
		s.ticker.Start()
		s.role.Shift(entity.Leader)
	}
}

func (s *Runner) GrantVote(id string, lastLogIndex uint32, lastLogTerm uint32) bool {
	defer s.saveState()
	state := s.GetState()

	currentLogIndex := state.LastLogIndex()
	currentLogTerm := state.LastLogTerm()
	votedFor := state.VotedFor

	voteAvailable := votedFor == "" || votedFor == id
	candidateUpToDate := currentLogTerm <= lastLogTerm && currentLogIndex <= lastLogIndex

	if voteAvailable && candidateUpToDate {
		s.state.SetVotedFor(id)
		return true
	}
	return false
}

func (s *Runner) RequestVoteInput() *rpc.RequestVoteInput {
	state := s.GetState()
	return &rpc.RequestVoteInput{
		CandidateId:  s.id,
		Term:         state.CurrentTerm,
		LastLogIndex: state.LastLogIndex(),
		LastLogTerm:  state.LastLogTerm(),
	}
}

func (s *Runner) AppendEntriesInput() *rpc.AppendEntriesInput {
	state := s.GetState()
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

	for range s.role.Signal() {
		switch {
		case s.role.Is(entity.Follower):
			go service.RunFollower(s)
		case s.role.Is(entity.Candidate):
			go service.RunCandidate(s)
		case s.role.Is(entity.Leader):
			go service.RunLeader(s)
		}

	}
}
