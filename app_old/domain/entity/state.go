package entity

import "sync"

// Rename logstate or fsmstate
// Separate Persistent

type Istate struct {
	Persistent
	commitIndex uint32
	lastApplied uint32
	nextIndex   map[string]uint32
	matchIndex  map[string]uint32
	mu          sync.RWMutex
}

func NewState(ps *Persistent) *Istate {
	return &Istate{}
}

func (s *Istate) SetCommitIndex(index uint32) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.commitIndex = index
}

type State struct {
	Persistent
	CommitIndex uint32
	LastApplied uint32
	NextIndex   map[string]uint32
	MatchIndex  map[string]uint32
}

func (s *Istate) GetState() *State {
	s.mu.RLock()
	defer s.mu.RUnlock()

	machineLogs := make([]MachineLog, len(s.MachineLogs))
	copy(machineLogs, s.MachineLogs)

	nextIndex := map[string]uint32{}
	for k, v := range s.nextIndex {
		nextIndex[k] = v
	}

	matchIndex := map[string]uint32{}
	for k, v := range s.matchIndex {
		matchIndex[k] = v
	}

	return &State{
		Persistent: Persistent{
			CurrentTerm: s.CurrentTerm,
			VotedFor:    s.VotedFor,
			MachineLogs: machineLogs,
		},
		CommitIndex: s.commitIndex,
		LastApplied: s.lastApplied,
		NextIndex:   nextIndex,
		MatchIndex:  matchIndex,
	}
}
func (s *State) LastLogIndex() uint32 {
	return uint32(len(s.MachineLogs))
}
func (s *State) GetLog(index uint32) MachineLog {
	if index < s.LastLogIndex() {
		log := MachineLog{
			Term:  s.MachineLogs[index].Term,
			Value: s.MachineLogs[index].Value,
		}
		return log
	}
	return MachineLog{Term: 0}
}
func (s *State) LastLogTerm() uint32 {
	if lastLogIndex := s.LastLogIndex(); lastLogIndex != 0 {
		return s.MachineLogs[lastLogIndex-1].Term
	}
	return 0
}
