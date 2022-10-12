package models

import (
	"encoding/json"
	"log"
	"math"
	"os"
	"sync"
)

const PERSISTENT_STATE_FILE = "orchestrator/state.json"

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

func (state *PersistentState) LastLogIndex() uint16 {
	return uint16(len(state.Logs))
}

func (state *PersistentState) LastLogTerm() uint16 {
	if lastLogIndex := state.LastLogIndex(); lastLogIndex != 0 {
		return uint16(state.Logs[lastLogIndex-1].Term)
	}
	return 0
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
	Name  string
	Nodes []Node

	heartbeat chan bool
	mu        sync.Mutex
}

func NewServerState() *ServerState {
	state := ServerState{
		Role:            Follower,
		PersistentState: PersistentState{},
		FollowerState:   FollowerState{CommitIndex: 0, LastApplied: 0},
		// heartbeat need to be buffered 1 otherwise it is blocking
		heartbeat: make(chan bool, 1),
	}
	state.PersistentState.loadState(PERSISTENT_STATE_FILE)
	return &state
}

func (state *ServerState) Ping() {
	state.mu.Lock()
	defer state.mu.Unlock()
	state.heartbeat <- true
}

func (state *ServerState) Heartbeat() chan bool {
	return state.heartbeat
}

func (state *ServerState) Quorum() int {
	totalNodes := len(state.Nodes) + 1 // add self
	return int(math.Ceil(float64(totalNodes) / 2.0))
}

func (state *ServerState) Vote(candidateId string) {
	state.mu.Lock()
	defer state.mu.Unlock()

	state.VotedFor = candidateId
}

func (state *ServerState) CanVote(candidateId string, candidateLastLogIndex int, candidateLastLogTerm int) bool {
	state.mu.Lock()
	defer state.mu.Unlock()

	currentLogIndex := state.LastLogIndex()
	currentLogTerm := state.LastLogTerm()

	voteAvailable := state.VotedFor == "" || state.VotedFor == candidateId
	candidateUpToDate := int(currentLogTerm) <= candidateLastLogTerm && currentLogIndex <= currentLogIndex

	return voteAvailable && candidateUpToDate
}

func (state *ServerState) isRole(role Role) bool {
	state.mu.Lock()
	defer state.mu.Unlock()
	return state.Role == role
}
func (state *ServerState) IsFollower() bool {
	return state.isRole(Follower)
}
func (state *ServerState) IsCandidate() bool {
	return state.isRole(Candidate)
}
func (state *ServerState) IsLeader() bool {
	return state.isRole(Leader)
}

func (state *ServerState) DowngradeToFollower(term uint16) {
	log.Printf("DOWNGRADE TO FOLLOWER TERM: %d\n", term)
	state.mu.Lock()
	defer state.mu.Unlock()

	state.Role = Follower
	state.CurrentTerm = term
	state.VotedFor = ""
	// state.saveState(PERSISTENT_STATE_FILE)
}

func (state *ServerState) RaiseToCandidate() {
	log.Printf("RAISE TO CANDIDATE TERM: %d\n", state.CurrentTerm+1)
	state.mu.Lock()
	defer state.mu.Unlock()

	state.Role = Candidate
	state.CurrentTerm += 1
	state.VotedFor = state.Name
	// state.saveState(PERSISTENT_STATE_FILE)
}

func (state *ServerState) PromoteToLeader() {
	log.Printf("PROMOTE TO LEADER TERM: %d\n", state.CurrentTerm)
	state.mu.Lock()
	defer state.mu.Unlock()

	state.Role = Leader
}
