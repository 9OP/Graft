package entity

type Peer struct {
	Id   string
	Host string
	Port string
}

type State struct {
	PersistentState
	CommitIndex uint32
	LastApplied uint32
	NextIndex   []string // leader only
	MatchIndex  []string // leader only
}

func NewState(ps *PersistentState) *State {
	return &State{
		PersistentState: *ps,
		CommitIndex:     0,
		LastApplied:     0,
		NextIndex:       []string{},
		MatchIndex:      []string{},
	}
}

type PersistentState struct {
	CurrentTerm uint32       `json:"current_term"`
	VotedFor    string       `json:"voted_for"`
	MachineLogs []MachineLog `json:"machine_logs"`
}

type MachineLog struct {
	Term  uint32 `json:"term"`
	Value string `json:"value"`
}

func (state *PersistentState) LastLogIndex() uint32 {
	return uint32(len(state.MachineLogs))
}

func (state *PersistentState) GetLogIndex(n int) MachineLog {
	if n < int(state.LastLogIndex()) {
		return state.MachineLogs[n]
	}
	return MachineLog{Term: 0}
}

func (state *PersistentState) LastLogTerm() uint32 {
	if lastLogIndex := state.LastLogIndex(); lastLogIndex != 0 {
		return uint32((state.MachineLogs)[lastLogIndex-1].Term)
	}
	return 0
}
