package service

import (
	"graft/app/domain/entity"
	"log"
	"sync"
)

type signals struct {
	SaveState          chan struct{}
	ShiftRole          chan entity.Role
	ResetElectionTimer chan struct{}
	ResetLeaderTicker  chan struct{}
	Commit             chan struct{}
}

func newSignals() signals {
	return signals{
		SaveState:          make(chan struct{}, 1),
		ShiftRole:          make(chan entity.Role, 1),
		ResetElectionTimer: make(chan struct{}, 1),
		ResetLeaderTicker:  make(chan struct{}, 1),
		Commit:             make(chan struct{}),
	}
}

func (s *signals) saveState() {
	s.SaveState <- struct{}{}
}
func (s *signals) shiftRole(role entity.Role) {
	s.ShiftRole <- role
}
func (s *signals) resetTimeout() {
	s.ResetElectionTimer <- struct{}{}
}
func (s *signals) resetLeaderTicker() {
	s.ResetLeaderTicker <- struct{}{}
}
func (s *signals) commit() {
	s.Commit <- struct{}{}
}

type Server struct {
	signals
	entity.Node
}

func NewServer(id string, peers entity.Peers, persistent *entity.Persistent) *Server {
	srv := &Server{
		signals: newSignals(),
		Node:    *entity.NewNode(id, peers, persistent),
	}
	srv.shiftRole(entity.Follower)
	srv.resetTimeout()
	return srv
}

func (s *Server) GetState() *entity.FsmState {
	// TODO: verify if we need to send a copy instead
	return s.Node.FsmState
}

func (s *Server) Heartbeat() {
	s.resetTimeout()
}

func (s *Server) Broadcast(fn func(p entity.Peer)) {
	var wg sync.WaitGroup
	for _, peer := range s.Node.Peers {
		wg.Add(1)
		go func(p entity.Peer, w *sync.WaitGroup) {
			defer w.Done()
			fn(p)
		}(peer, &wg)
	}
	wg.Wait()
}

func (s *Server) DeleteLogsFrom(index uint32) {
	s.Node.DeleteLogsFrom(index)
	s.saveState()
}

func (s *Server) AppendLogs(entries []entity.LogEntry, prevLogIndex uint32) {
	if s.Node.AppendLogs(entries, prevLogIndex) {
		s.saveState()
	}
}

func (s *Server) SetCommitIndex(index uint32) {
	if s.Node.SetCommitIndex(index) {
		s.commit()
	}
}

// If server can grant vote, it will set votedFor and return success
func (s *Server) GrantVote(id string, lastLogIndex uint32, lastLogTerm uint32) bool {
	defer s.saveState()
	if s.Node.CanGrantVote(id, lastLogIndex, lastLogTerm) {
		s.Node.SetVotedFor(id)
		return true
	}
	return false
}

func (s *Server) SetRole(role entity.Role) {
	s.Node.SetRole(role)
	s.signals.shiftRole(role)
}

func (s *Server) DowngradeFollower(term uint32) {
	log.Printf("DOWNGRADE TO FOLLOWER TERM: %d\n", term)
	s.SetCurrentTerm(term)
	s.SetVotedFor("")
	s.SetRole(entity.Follower)
	s.saveState()
	s.resetTimeout()
}

func (s *Server) IncrementTerm() {
	if s.IsRole(entity.Candidate) {
		log.Printf("INCREMENT CANDIDATE TERM: %d\n", s.CurrentTerm+1)
		s.SetCurrentTerm(s.CurrentTerm + 1)
		s.SetVotedFor(s.Id)
		s.saveState()
		s.resetTimeout()
	}
}

func (s *Server) UpgradeCandidate() {
	if s.IsRole(entity.Follower) {
		log.Println("UPGRADE TO CANDIDATE")
		s.SetRole(entity.Candidate)
	}
}

func (s *Server) UpgradeLeader() {
	if s.IsRole(entity.Candidate) {
		log.Printf("UPGRADE TO LEADER TERM: %d\n", s.CurrentTerm)
		s.InitializeLeader(s.Peers)
		s.SetClusterLeader(s.Id)
		s.SetRole(entity.Leader)
		s.resetLeaderTicker()
	}
}

func (s *Server) Execute(entry string) chan interface{} {
	state := s.GetState()
	result := make(chan interface{}, 1)
	newLog := entity.LogEntry{Value: entry, Term: state.CurrentTerm, C: result}
	s.AppendLogs([]entity.LogEntry{newLog}, state.GetLastLogIndex())
	return result
}
