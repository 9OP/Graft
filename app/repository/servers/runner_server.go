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

type Runner struct {
	id            string
	role          *entity.Role
	peers         []entity.Peer
	clusterLeader string
	state         *entity.State
	timeout       *entity.Timeout
	ticker        *entity.Ticker
	persister     *persister.Service
	mu            sync.Mutex
}

func NewRunner(id string, peers []entity.Peer, persister *persister.Service, timeout *entity.Timeout, ticker *entity.Ticker) *Runner {

	ps, _ := persister.LoadState()
	state := entity.NewState(ps)
	role := entity.NewRole()

	srv := &Runner{
		id:        id,
		peers:     peers,
		state:     state,
		timeout:   timeout,
		ticker:    ticker,
		role:      role,
		persister: persister,
	}

	return srv
}

func (s *Runner) GetState() entity.ImmerState {
	// Returns a new copy of the server state

	var logs []entity.MachineLog
	copy(logs, s.state.MachineLogs)

	return entity.ImmerState{
		PersistentState: entity.PersistentState{
			CurrentTerm: s.state.CurrentTerm,
			VotedFor:    s.state.VotedFor,
			MachineLogs: logs,
		},
		CommitIndex: s.state.CommitIndex,
		LastApplied: s.state.LastApplied,
		NextIndex:   s.state.NextIndex,  // do deep copy
		MatchIndex:  s.state.MatchIndex, // do deep copy
	}
}

func (s *Runner) GetQuorum() int {
	totalNodes := len(s.peers) + 1 // add self
	return int(math.Ceil(float64(totalNodes) / 2.0))
}

func (s *Runner) SetClusterLeader(leaderId string) {
	if s.clusterLeader != leaderId {
		log.Printf("FOLLOWING CLUSTER LEADER: %s\n", leaderId)
		s.mu.Lock()
		s.clusterLeader = leaderId
		s.mu.Unlock()
	}
}

func (s *Runner) saveState() {
	s.mu.Lock()
	s.persister.SaveState(&s.state.PersistentState)
	s.mu.Unlock()
}

func (s *Runner) resetTimeout() {
	s.timeout.RReset()
	s.ticker.Stop()
}

func (s *Runner) Heartbeat() {
	s.resetTimeout()
}

func (s *Runner) DeleteLogsFrom(n int) {
	s.state.PersistentState.DeleteLogFrom(n)
	s.saveState()
}

func (s *Runner) AppendLogs(entries []string) {
	s.state.PersistentState.AppendLogs(entries)
	s.saveState()
}

func (s *Runner) SetCommitIndex(ind uint32) {
	s.state.CommitIndex = ind
}

func (s *Runner) DowngradeFollower(term uint32, leaderId string) {
	log.Printf("DOWNGRADE TO FOLLOWER TERM: %d\n", term)
	s.resetTimeout()
	s.SetClusterLeader(leaderId)

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