package models

import (
	"encoding/json"
	"os"
)

const PERSISTENT_STATE_FILE = "state.json"
const HEARTBEAT = 3000 // ms

type PersistentState struct {
	CurrentTerm uint16
	VotedFor    string
	Log         []Log
}

type Log struct {
	Term  uint16
	Value string
}

func (state *PersistentState) saveState(location string) error {
	data, err := json.Marshal(*state)
	if err != nil {
		return err
	}
	return os.WriteFile(location, data, 0644)
}

func (state *PersistentState) loadState(location string) error {
	data, err := os.ReadFile(location)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, state)
}

type FollowerState struct {
	CommitIndex uint16
	LastApplied uint16
}

type LeaderState struct {
	FollowerState
	NextIndex  []string
	MatchIndex []string
}

type Role struct {
	slug string
}

func (r Role) String() string {
	return r.slug
}

var (
	Unknown   = Role{""}
	Follower  = Role{"Follower"}
	Candidate = Role{"Candidate"}
	Leader    = Role{"Leader"}
)

type Node struct {
	Name string
	Host string
}
type ServerState struct {
	PersistentState
	FollowerState
	Role
	Heartbeat chan bool
	Name      string
	Nodes     []Node
}

func NewServerState() *ServerState {
	state := ServerState{
		Role:            Follower,
		PersistentState: PersistentState{},
		FollowerState:   FollowerState{CommitIndex: 0, LastApplied: 0},
		Heartbeat:       make(chan bool),
	}
	state.PersistentState.loadState(PERSISTENT_STATE_FILE)
	return &state
}
