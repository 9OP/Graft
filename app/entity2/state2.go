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
func (s *state) CommitIndex() uint32 {
	return s.commitIndex
}
func (s *state) LastApplied() uint32 {
	return s.lastApplied
}
func (s *state) NextIndex() map[string]uint32 {
	cp := make(map[string]uint32, len(s.nextIndex))
	for k, v := range s.nextIndex {
		cp[k] = v
	}
	return cp
}
func (s *state) MatchIndex() map[string]uint32 {
	cp := make(map[string]uint32, len(s.matchIndex))
	for k, v := range s.matchIndex {
		cp[k] = v
	}
	return cp
}
func (s *state) CurrentTerm() uint32 {
	return s.persistent.CurrentTerm
}
func (s *state) VotedFor() string {
	return s.persistent.VotedFor
}
func (s *state) MachineLogs() []machineLog {
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
func DeleteLogFrom(s state, n uint32) state {
	machineLogs := s.MachineLogs()
	lastLogIndex := s.LastLogIndex()

	if n < lastLogIndex {
		machineLogs = machineLogs[:n]
	}

	return state{
		persistent: persistentState{
			CurrentTerm: s.CurrentTerm(),
			VotedFor:    s.VotedFor(),
			MachineLogs: machineLogs,
		},
		commitIndex: s.CommitIndex(),
		lastApplied: s.LastApplied(),
		nextIndex:   s.NextIndex(),
		matchIndex:  s.MatchIndex(),
	}
}
func AppendLogs(s state, entries []string) state {
	machineLogs := s.MachineLogs()
	currentTerm := s.CurrentTerm()
	for _, entry := range entries {
		machineLogs = append(machineLogs, machineLog{Term: currentTerm, Value: entry})
	}

	return state{
		persistent: persistentState{
			CurrentTerm: s.CurrentTerm(),
			VotedFor:    s.VotedFor(),
			MachineLogs: machineLogs,
		},
		commitIndex: s.CommitIndex(),
		lastApplied: s.LastApplied(),
		nextIndex:   s.NextIndex(),
		matchIndex:  s.MatchIndex(),
	}

}
