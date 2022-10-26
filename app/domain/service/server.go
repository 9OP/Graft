package service

import (
	"graft/app/domain/entity"
	"log"
	"sync"
)

// Create server
// Inject persistent
// Inject timeout

// inject server as repository

// Domain logic here

// Busines logic in use case

// func test() {
// 	n := entity.NewNode("", nil)
// 	n.
// }

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
	// Hide all node methods from outside
	node entity.Node
	mu   sync.RWMutex
}

// func NewServer(id string, peers []Peer, ps *Persistent) *Server {
// 	state := NewState(ps)
// 	role := entitynew.Follower

// 	srv := &Server{
// 		Id:    id,
// 		Peers: peers,
// 		role:  role,
// 		state: state,
// 		Signals: Signals{
// 			SaveState:          make(chan struct{}, 1),
// 			ShiftRole:          make(chan struct{}, 1),
// 			ResetElectionTimer: make(chan struct{}, 1),
// 			ResetLeaderTicker:  make(chan struct{}, 1),
// 		},
// 	}
// 	srv.ShiftRole <- struct{}{}
// 	return srv
// }

func (s *Server) Heartbeat() {
	// Dispatch application of FSM / commitIndex / lastApplied ?
	s.resetTimeout()
}

func (s *Server) GetState() entity.FsmState {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.node.GetState()
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
		// s.node.SetClusterLeader(s.node.Id)
		s.node.SetRole(entity.Leader)
		s.resetLeaderTicker()
		s.shiftRole()
	}
}

// // func (s *Server) GrantVote(id string, lastLogIndex uint32, lastLogTerm uint32) bool {
// // 	defer s.saveState()
// // 	state := s.GetState()

// // 	currentLogIndex := state.LastLogIndex()
// // 	currentLogTerm := state.LastLogTerm()
// // 	votedFor := state.VotedFor

// // 	voteAvailable := votedFor == "" || votedFor == id
// // 	candidateUpToDate := currentLogTerm <= lastLogTerm && currentLogIndex <= lastLogIndex

// // 	if voteAvailable && candidateUpToDate {
// // 		s.state.SetVotedFor(id)
// // 		return true
// // 	}
// // 	return false
// // }

func (s *Server) Broadcast(fn func(peer entity.Peer)) {
	peers := s.node.Peers

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
