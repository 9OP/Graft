package entity2

func NewState() *state {
	return &state{
		persistent: persistentState{
			CurrentTerm: 0,
			VotedFor:    "",
			MachineLogs: []machineLog{},
		},
		commitIndex: 0,
		lastApplied: 0,
		nextIndex:   map[string]uint32{},
		matchIndex:  map[string]uint32{},
	}
}

/*
Goal:
- GetState returns a state immutable that does not contain methods to update srv.state
- srv contains a state and can mute the state with methods

*/

// Public
type State struct {
	Persistent
	CommitIndex uint32
	LastApplied uint32
	NextIndex   map[string]uint32
	MatchIndex  map[string]uint32
}

// Public
type Persistent struct {
	persistentState
}

type state struct {
	persistent  persistentState
	commitIndex uint32
	lastApplied uint32
	nextIndex   map[string]uint32
	matchIndex  map[string]uint32
}

type persistentState struct {
	CurrentTerm uint32       `json:"current_term"`
	VotedFor    string       `json:"voted_for"`
	MachineLogs []machineLog `json:"machine_logs"`
}

type machineLog struct {
	Term  uint32 `json:"term"`
	Value string `json:"value"`
}

// Getters
func (s *state) GetImmutable() *State {
	return &State{
		Persistent:  Persistent{},
		CommitIndex: s.commitIndex,
		LastApplied: s.lastApplied,
		NextIndex:   s.copyNextIndex(),
		MatchIndex:  s.copyMatchIndex(),
	}
}
func (s *state) copyNextIndex() map[string]uint32 {
	cp := make(map[string]uint32, len(s.nextIndex))
	for k, v := range s.nextIndex {
		cp[k] = v
	}
	return cp
}
func (s *state) copyMatchIndex() map[string]uint32 {
	cp := make(map[string]uint32, len(s.matchIndex))
	for k, v := range s.matchIndex {
		cp[k] = v
	}
	return cp
}
func (s *state) copyMachineLogs() []machineLog {
	cp := make([]machineLog, len(s.persistent.MachineLogs))
	copy(cp, s.persistent.MachineLogs)
	return cp
}
func (s *state) LastLogIndex() uint32 {
	return uint32(len(s.persistent.MachineLogs))
}
func (s *state) GetLogIndex(n int) machineLog {
	if n < int(s.LastLogIndex()) {
		return s.persistent.MachineLogs[n]
	}
	return machineLog{Term: 0}
}
func (s *state) LastLogTerm() uint32 {
	if lastLogIndex := s.LastLogIndex(); lastLogIndex != 0 {
		return s.persistent.MachineLogs[lastLogIndex-1].Term
	}
	return 0
}

// Transforms
// func DeleteLogFrom(s state, n uint32) state {
// 	machineLogs := s.copyMachineLogs()
// 	lastLogIndex := s.LastLogIndex()

// 	if n < lastLogIndex {
// 		machineLogs = machineLogs[:n]
// 	}

// 	return state{
// 		persistent: persistentState{
// 			CurrentTerm: s.CurrentTerm(),
// 			VotedFor:    s.VotedFor(),
// 			MachineLogs: machineLogs,
// 		},
// 		commitIndex: s.CommitIndex(),
// 		lastApplied: s.LastApplied(),
// 		nextIndex:   s.NextIndex(),
// 		matchIndex:  s.MatchIndex(),
// 	}
// }
// func AppendLogs(s state, entries []string) state {
// 	machineLogs := s.MachineLogs()
// 	currentTerm := s.CurrentTerm()
// 	for _, entry := range entries {
// 		machineLogs = append(machineLogs, machineLog{Term: currentTerm, Value: entry})
// 	}

// 	return state{
// 		persistent: persistentState{
// 			CurrentTerm: s.CurrentTerm(),
// 			VotedFor:    s.VotedFor(),
// 			MachineLogs: machineLogs,
// 		},
// 		commitIndex: s.CommitIndex(),
// 		lastApplied: s.LastApplied(),
// 		nextIndex:   s.NextIndex(),
// 		matchIndex:  s.MatchIndex(),
// 	}

// }
