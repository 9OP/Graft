package entity

import (
	"encoding/json"
	"os"
)

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

type PersistentState struct {
	CurrentTerm uint32       `json:"current_term"`
	VotedFor    string       `json:"voted_for"`
	MachineLogs []MachineLog `json:"machine_logs"`
}

type MachineLog struct {
	Term  uint32 `json:"term"`
	Value string `json:"value"`
}

var DEFAULT_STATE PersistentState = PersistentState{
	CurrentTerm: 0,
	VotedFor:    "",
	MachineLogs: []MachineLog{{Term: 0, Value: ""}},
}

func (state *PersistentState) SaveState(location string) error {
	data, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(location, data, 0644)
}

func (state *PersistentState) LoadState(location string) error {
	data, err := os.ReadFile(location)
	if err != nil {
		state.CurrentTerm = DEFAULT_STATE.CurrentTerm
		state.VotedFor = DEFAULT_STATE.VotedFor
		state.MachineLogs = DEFAULT_STATE.MachineLogs
		return err
	}
	return json.Unmarshal(data, state)
}

func (state *PersistentState) LastLogIndex() uint32 {
	return uint32(len(state.MachineLogs))
}

func (state *PersistentState) LastLogTerm() uint32 {
	if lastLogIndex := state.LastLogIndex(); lastLogIndex != 0 {
		return uint32((state.MachineLogs)[lastLogIndex-1].Term)
	}
	return 0
}
