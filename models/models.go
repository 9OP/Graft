package models

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

const PERSISTENT_STATE_FILE = "state.json"

type PersistentState struct {
	CurrentTerm uint16
	VotedFor    string
	Log         []string
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

type ServerState struct {
	PersistentState
	FollowerState
	Role
	Heartbeat chan bool
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

type Server struct {
	Timeout *time.Ticker
}

func (s *Server) Start(state *ServerState) {
	for {
		select {
		case <-s.Timeout.C:
			fmt.Println("Server timeout, start election")
		case <-state.Heartbeat:
			fmt.Printf("Received heartbeat %v\n", <-state.Heartbeat)
			s.Timeout.Stop()
			s.Timeout = time.NewTicker(3000 * time.Millisecond)
		}
	}
}
