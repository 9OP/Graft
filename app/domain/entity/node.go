package entity

import (
	"log"
	"math"
)

// Write here every function that does not need to send external signal

// DDD aggregate over role, peers, fsmState(persistent)
type node struct {
	id       string
	leaderId string
	role     Role
	peers    Peers
	fsmState *fsmState
}

func NewNode(id string, peers Peers) *node {

	return &node{
		id:       id,
		peers:    peers,
		fsmState: NewFsmState(),
	}
}

func (n node) GetState() FsmState {
	return n.fsmState.getStateCopy()
}

func (n node) GetQuorum() int {
	totalClusterNodes := len(n.peers) + 1 // add self
	return int(math.Ceil(float64(totalClusterNodes) / 2.0))
}

func (n node) GetPeers() Peers {
	return n.peers
}

func (n *node) SetClusterLeader(leaderId string) {
	if n.leaderId != leaderId {
		log.Printf("FOLLOWING CLUSTER LEADER: %s\n", leaderId)
		n.leaderId = leaderId
	}
}

func (n node) isRole(role Role) bool {
	return n.role == role
}
func (n node) IsFollower() bool {
	return n.isRole(Follower)
}
func (n node) IsCandidate() bool {
	return n.isRole(Candidate)
}
func (n node) IsLeader() bool {
	return n.isRole(Leader)
}

func (n node) GrantVote(id string, lastLogIndex uint32, lastLogTerm uint32) bool {
	state := n.GetState()
	currentLogIndex := state.GetLastLogIndex()
	currentLogTerm := state.GetLastLogTerm()
	votedFor := state.VotedFor

	voteAvailable := votedFor == "" || votedFor == id
	candidateUpToDate := currentLogTerm <= lastLogTerm && currentLogIndex <= lastLogIndex
	if voteAvailable && candidateUpToDate {
		return true
	}
	return false
}
