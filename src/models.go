package main

import (
	"log"
	"math"
	"sync"
)

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
	id   string
	host string
	port string
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

func NewServerState(name string, nodes *[]Node) *ServerState {
	state := ServerState{
		Role:            Follower,
		PersistentState: PersistentState{},
		FollowerState:   FollowerState{CommitIndex: 0, LastApplied: 0},
		Name:            name,
		Nodes:           *nodes,
		// heartbeat need to be buffered 1 otherwise it is blocking
		heartbeat: make(chan bool, 1),
	}
	state.loadState()
	return &state
}

func (state *ServerState) persistenceLocation() string {
	location := "orchestrator/state_" + state.Name + ".json"
	return location
}
func (state *ServerState) loadState() {
	state.PersistentState.loadState(state.persistenceLocation())
}
func (state *ServerState) saveState() {
	state.PersistentState.saveState(state.persistenceLocation())
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
	state.saveState()
}

func (state *ServerState) CanVote(candidateId string, candidateLastLogIndex int, candidateLastLogTerm int) bool {
	state.mu.Lock()
	defer state.mu.Unlock()

	currentLogIndex := int(state.LastLogIndex())
	currentLogTerm := state.LastLogTerm()

	voteAvailable := state.VotedFor == "" || state.VotedFor == candidateId
	candidateUpToDate := int(currentLogTerm) <= candidateLastLogTerm && currentLogIndex <= candidateLastLogIndex

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
	state.saveState()
}

func (state *ServerState) RaiseToCandidate() {
	log.Printf("RAISE TO CANDIDATE TERM: %d\n", state.CurrentTerm+1)
	state.mu.Lock()
	defer state.mu.Unlock()

	state.Role = Candidate
	state.CurrentTerm += 1
	state.VotedFor = state.Name
	state.saveState()
}

func (state *ServerState) PromoteToLeader() {
	log.Printf("PROMOTE TO LEADER TERM: %d\n", state.CurrentTerm)
	state.mu.Lock()
	defer state.mu.Unlock()

	state.Role = Leader
}