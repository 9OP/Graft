package service

import (
	"graft/app/domain/entity"
	"log"
	"sync"
)

type Signals struct {
	SaveState          chan struct{}
	ShiftRole          chan struct{}
	ResetElectionTimer chan struct{}
	ResetLeaderTicker  chan struct{}
}

func (s *Signals) saveState() {
	s.SaveState <- struct{}{}
}
func (s *Signals) resetTimeout() {
	s.ResetElectionTimer <- struct{}{}
}
func (s *Signals) resetLeaderTicker() {
	s.ResetLeaderTicker <- struct{}{}
}
func (s *Signals) shiftRole() {
	s.ShiftRole <- struct{}{}
}

type Server struct {
	Signals
	node entity.Node
	mu   sync.RWMutex
}

func NewServer(id string, peers entity.Peers) *Server {
	srv := &Server{
		Signals: Signals{
			SaveState:          make(chan struct{}, 1),
			ShiftRole:          make(chan struct{}, 1),
			ResetElectionTimer: make(chan struct{}, 1),
			ResetLeaderTicker:  make(chan struct{}, 1),
		},
		node: *entity.NewNode(id, peers),
	}
	srv.shiftRole()
	srv.resetTimeout()
	return srv
}

func (s *Server) GetState() entity.FsmState {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.node.GetState()
}

func (s *Server) Heartbeat() {
	// Dispatch application of FSM / commitIndex / lastApplied ?
	s.resetTimeout()
}

func (s *Server) Broadcast(fn func(p entity.Peer)) {
	// Check if it raise execption on concurrency
	// if yes, define broadcast here instead of node
	s.node.Broadcast(fn)
}

func (s *Server) SetClusterLeader(leaderId string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.node.SetClusterLeader(leaderId)
}
func (s *Server) GetQuorum() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.node.GetQuorum()
}
func (s *Server) GetAppendEntriesInput() entity.AppendEntriesInput {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.node.GetAppendEntriesInput()
}
func (s *Server) GetRequestVoteInput() entity.RequestVoteInput {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.node.GetRequestVoteInput()
}
func (s *Server) IsFollower() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.node.IsFollower()
}
func (s *Server) IsCandidate() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.node.IsCandidate()
}
func (s *Server) IsLeader() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.node.IsLeader()
}

func (s *Server) DeleteLogsFrom(index uint32) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.node.DeleteLogsFrom(index)
	s.saveState()
}

func (s *Server) AppendLogs(entries []string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.node.AppendLogs(entries)
	s.saveState()
}

func (s *Server) SetCommitIndex(index uint32) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.node.SetCommitIndex(index)
}

// If server can grant vote, it will set votedFor and return success
func (s *Server) GrantVote(id string, lastLogIndex uint32, lastLogTerm uint32) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.node.CanGrantVote(id, lastLogIndex, lastLogTerm) {
		s.node.SetVotedFor(id)
		s.saveState()
		return true
	}
	return false
}

func (s *Server) DowngradeFollower(term uint32, leaderId string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	log.Printf("DOWNGRADE TO FOLLOWER TERM: %d\n", term)
	s.node.SetClusterLeader(leaderId)
	s.node.SetCurrentTerm(term)
	s.node.SetVotedFor("")
	s.node.SetRole(entity.Follower)
	s.saveState()
	s.resetTimeout()
	s.shiftRole()
}

func (s *Server) IncrementTerm() {
	state := s.GetState()
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.node.IsCandidate() {
		log.Printf("INCREMENT CANDIDATE TERM: %d\n", state.CurrentTerm+1)
		s.node.SetCurrentTerm(state.CurrentTerm + 1)
		s.node.VoteForSelf()
		s.saveState()
		s.resetTimeout()
	}
}

func (s *Server) UpgradeCandidate() {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.node.IsFollower() {
		log.Println("UPGRADE TO CANDIDATE")
		s.node.SetRole(entity.Candidate)
		s.shiftRole()
	}
}

func (s *Server) UpgradeLeader() {
	state := s.GetState()
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.node.IsCandidate() {
		log.Printf("UPGRADE TO LEADER TERM: %d\n", state.CurrentTerm)
		s.node.SetClusterLeader(s.node.GetId())
		s.node.SetRole(entity.Leader)
		s.resetLeaderTicker()
		s.shiftRole()
	}
}
