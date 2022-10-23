package entity2

import "sync"

type state struct {
	persistent  Persistent
	commitIndex uint32
	lastApplied uint32
	nextIndex   map[string]uint32
	matchIndex  map[string]uint32
	mu          sync.RWMutex
}

func NewState() *state {
	return &state{}
}

type Persistent struct {
	CurrentTerm uint32       `json:"current_term"`
	VotedFor    string       `json:"voted_for"`
	MachineLogs []MachineLog `json:"machine_logs"`
}

type MachineLog struct {
	Term  uint32 `json:"term"`
	Value string `json:"value"`
}

func (s *state) DeleteLogsFrom(index uint32) {
	s.mu.Lock()
	defer s.mu.Unlock()
	lastLogIndex := int32(len(s.persistent.MachineLogs))
	if index < uint32(lastLogIndex) {
		s.persistent.MachineLogs = s.persistent.MachineLogs[:index]
	}
}
func (s *state) AppendLogs(entries []string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	size := len(entries) + len(s.persistent.MachineLogs)
	currentTerm := s.persistent.CurrentTerm
	logs := make([]MachineLog, size)
	for _, entry := range entries {
		logs = append(
			logs,
			MachineLog{Term: currentTerm, Value: entry},
		)
	}
	s.persistent.MachineLogs = logs
}

type State struct {
	Persistent
	CommitIndex uint32
	LastApplied uint32
	NextIndex   map[string]uint32
	MatchIndex  map[string]uint32
}

func (s *state) GetState() *State {
	s.mu.RUnlock()
	defer s.mu.RUnlock()

	machineLogs := make([]MachineLog, len(s.persistent.MachineLogs))
	copy(machineLogs, s.persistent.MachineLogs)

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
			CurrentTerm: s.persistent.CurrentTerm,
			VotedFor:    s.persistent.VotedFor,
			MachineLogs: machineLogs,
		},
		CommitIndex: s.commitIndex,
		LastApplied: s.lastApplied,
		NextIndex:   nextIndex,
		MatchIndex:  matchIndex,
	}
}

func (s *State) GetLog(index int) MachineLog {
	if index < int(s.LastLogIndex()) {
		log := MachineLog{
			Term:  s.MachineLogs[index].Term,
			Value: s.MachineLogs[index].Value,
		}
		return log
	}
	return MachineLog{Term: 0}
}
func (s *State) LastLogIndex() uint32 {
	return uint32(len(s.MachineLogs))
}
func (s *State) LastLogTerm() uint32 {
	if lastLogIndex := s.LastLogIndex(); lastLogIndex != 0 {
		return s.MachineLogs[lastLogIndex-1].Term
	}
	return 0
}
