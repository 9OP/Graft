package models

import (
	"encoding/json"
	"os"
	"sync"
)

const PERSISTENT_STATE_FILE = "orchestrator/state.json"
const HEARTBEAT = 3000 // ms

type PersistentState struct {
	CurrentTerm uint16 `json:"current_term"`
	VotedFor    string `json:"voted_for"`
	Logs        []Log  `json:"logs"`
}

type Log struct {
	Term  uint16 `json:"term"`
	Value string `json:"value"`
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
	mu        sync.Mutex // required for safe mutation accross go routines
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

func (state *ServerState) IsRole(role Role) bool {
	state.mu.Lock()
	defer state.mu.Unlock()
	return state.Role == role
}

func (state *ServerState) SwitchRole(role Role) {
	state.mu.Lock()
	defer state.mu.Unlock()
	state.Role = role
}

func (state *ServerState) FallbackToFollower(term uint16) {
	state.mu.Lock()
	defer state.mu.Unlock()
	state.SwitchRole(Follower)
	state.CurrentTerm = term
	state.VotedFor = ""
}
