package entity

import "sync"

type FsmState struct {
	*Persistent
	CommitIndex uint32
	LastApplied uint32
	NextIndex   map[string]uint32
	MatchIndex  map[string]uint32
	mu          sync.RWMutex
}

func NewFsmState(persistent *Persistent) *FsmState {
	return &FsmState{
		Persistent:  persistent,
		CommitIndex: 0,
		LastApplied: 0,
		NextIndex:   map[string]uint32{},
		MatchIndex:  map[string]uint32{},
	}
}

func (s *FsmState) GetCopy() *FsmState {
	s.mu.RLock()
	defer s.mu.RUnlock()

	nextIndex := make(map[string]uint32, len(s.NextIndex))
	for id, index := range s.NextIndex {
		nextIndex[id] = index
	}

	matchIndex := make(map[string]uint32, len(s.MatchIndex))
	for id, index := range s.MatchIndex {
		matchIndex[id] = index
	}

	return &FsmState{
		Persistent:  s.Persistent.GetCopy(),
		CommitIndex: s.CommitIndex,
		LastApplied: s.LastApplied,
		NextIndex:   nextIndex,
		MatchIndex:  matchIndex,
	}
}

func (s *FsmState) InitializeLeader(peers Peers) {
	s.mu.Lock()
	defer s.mu.Unlock()

	defaultNextIndex := s.GetLastLogIndex() + 1
	nextIndex := make(map[string]uint32, len(peers))
	matchIndex := make(map[string]uint32, len(peers))

	for _, peer := range peers {
		nextIndex[peer.Id] = defaultNextIndex
		matchIndex[peer.Id] = 0
	}

	s.NextIndex = nextIndex
	s.MatchIndex = matchIndex
}

// func (s *FsmState) IncrementLastApplied() {
// 	s.mu.Lock()
// 	defer s.mu.Unlock()
// 	s.LastApplied += 1
// }

func (s *FsmState) SetLastApplied(lastApplied uint32) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.LastApplied = lastApplied
}

func (s *FsmState) DecrementNextIndex(pId string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.NextIndex[pId] > 0 {
		s.NextIndex[pId] -= 1
	}
}

// Rename SetPeerNextIndex
func (s *FsmState) SetNextIndex(pId string, index uint32) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.NextIndex[pId] = index
}

// Rename SetPeerMatchIndex
func (s *FsmState) SetMatchIndex(pId string, index uint32) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.MatchIndex[pId] = index
}

// Return true when commit index was updated
func (s *FsmState) SetCommitIndex(index uint32) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	update := !(s.CommitIndex == index)
	s.CommitIndex = index
	return update
}
