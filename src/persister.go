package main

import (
	"encoding/json"
	"os"
)

type PersistentState struct {
	CurrentTerm uint16       `json:"current_term"`
	VotedFor    string       `json:"voted_for"`
	MachineLogs []MachineLog `json:"machine_logs"`
}

type MachineLog struct {
	Term  uint16 `json:"term"`
	Value string `json:"value"`
}

var DEFAULT_STATE PersistentState = PersistentState{
	CurrentTerm: 0,
	VotedFor:    "",
	MachineLogs: []MachineLog{{Term: 0, Value: ""}},
}

func (state *PersistentState) saveState(location string) error {
	data, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(location, data, 0644)
}

func (state *PersistentState) loadState(location string) error {
	data, err := os.ReadFile(location)
	if err != nil {
		state.CurrentTerm = DEFAULT_STATE.CurrentTerm
		state.VotedFor = DEFAULT_STATE.VotedFor
		state.MachineLogs = DEFAULT_STATE.MachineLogs
		return err
	}
	return json.Unmarshal(data, state)
}

func (state *PersistentState) LastLogIndex() uint16 {
	return uint16(len(state.MachineLogs))
}

func (state *PersistentState) LastLogTerm() uint16 {
	if lastLogIndex := state.LastLogIndex(); lastLogIndex != 0 {
		return uint16((state.MachineLogs)[lastLogIndex-1].Term)
	}
	return 0
}
