package service

import (
	"graft/app/domain/entity"
	"sync"
)

// Create server
// Inject persistent
// Inject timeout

// inject server as repository

// Domain logic here

// Busines logic in use case

func test() {
	n := entity.NewNode("", nil)
}

type Signals struct {
	SaveState          chan struct{}
	ShiftRole          chan struct{}
	ResetElectionTimer chan struct{}
	ResetLeaderTicker  chan struct{}
}

type Server struct {
	Signals
	mu sync.RWMutex
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

// func (s *Server) saveState() {
// 	s.SaveState <- struct{}{}
// }

// func (s *Server) resetTimeout() {
// 	s.ResetElectionTimer <- struct{}{}
// }

// func (s *Server) Heartbeat() {
// 	// Dispatch application of FSM / commitIndex / lastApplied ?
// 	s.resetTimeout()
// }

// func (s *Server) shiftRole(r entitynew.Role) {
// 	s.mu.RLock()
// 	defer s.mu.RUnlock()
// 	s.role = r
// 	s.ShiftRole <- struct{}{}
// }

// func (s *Server) DeleteLogsFrom(index uint32) {
// 	defer s.saveState()
// 	s.state.DeleteLogsFrom(index)
// }

// func (s *Server) AppendLogs(entries []string) {
// 	defer s.saveState()
// 	s.state.AppendLogs(entries)
// }

// func (s *Server) SetCommitIndex(index uint32) {
// 	s.state.SetCommitIndex(index)
// }

// func (s *Server) DowngradeFollower(term uint32, leaderId string) {
// 	defer s.saveState()
// 	log.Printf("DOWNGRADE TO FOLLOWER TERM: %d\n", term)
// 	s.resetTimeout()
// 	s.SetClusterLeader(leaderId)
// 	s.state.SetCurrentTerm(term)
// 	s.state.SetVotedFor("")
// 	s.shiftRole(entitynew.Follower)
// }

// func (s *Server) IncrementTerm() {
// 	if s.IsCandidate() {
// 		defer s.saveState()
// 		state := s.GetState()
// 		newTerm := state.CurrentTerm + 1
// 		log.Printf("INCREMENT CANDIDATE TERM: %d\n", newTerm)
// 		s.resetTimeout()
// 		s.state.SetCurrentTerm(newTerm)
// 		s.state.SetVotedFor(s.Id)
// 	}
// }

// func (s *Server) UpgradeCandidate() {
// 	if s.IsFollower() {
// 		log.Println("UPGRADE TO CANDIDATE")
// 		s.shiftRole(entitynew.Candidate)
// 	}
// }

// func (s *Server) UpgradeLeader() {
// 	if s.IsCandidate() {
// 		state := s.GetState()
// 		log.Printf("UPGRADE TO LEADER TERM: %d\n", state.CurrentTerm)
// 		s.SetClusterLeader(s.Id)
// 		s.ResetLeaderTicker <- struct{}{}
// 		s.shiftRole(entitynew.Leader)
// 	}
// }

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

// func (s *Server) RequestVoteInput() *rpc.RequestVoteInput {
// 	state := s.GetState()
// 	return &rpc.RequestVoteInput{
// 		CandidateId:  s.Id,
// 		Term:         state.CurrentTerm,
// 		LastLogIndex: state.LastLogIndex(),
// 		LastLogTerm:  state.LastLogTerm(),
// 	}
// }

// func (s *Server) AppendEntriesInput() *rpc.AppendEntriesInput {
// 	state := s.GetState()
// 	return &rpc.AppendEntriesInput{
// 		LeaderId: s.Id,
// 		Term:     state.CurrentTerm,
// 	}
// }

// func (s *Server) Broadcast(fn func(peer Peer)) {
// 	var wg sync.WaitGroup
// 	for _, peer := range s.Peers {
// 		wg.Add(1)
// 		go func(p Peer, w *sync.WaitGroup) {
// 			defer w.Done()
// 			fn(p)
// 		}(peer, &wg)
// 	}
// 	wg.Wait()
// }
