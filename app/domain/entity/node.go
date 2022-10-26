package entity

import (
	"fmt"
	"log"
	"math"
	"sync"
)

// Write here every function that does not need to send external signal

// DDD aggregate over role, peers, fsmState(persistent)
type Node struct {
	id        string
	leaderId  string
	peers     Peers
	role      Role
	*fsmState // ensure that we can call public method of private fsmState
	// WARNING: might cause issue with concurrent access since pointer and not
	// value copy (in GrantVote for instance)
}

func NewNode(id string, peers Peers) *Node {
	fmt.Println("peers", peers)
	return &Node{
		id:       id,
		peers:    peers,
		role:     Follower,
		fsmState: NewFsmState(),
	}
}

func (n Node) GetId() string {
	return n.id
}

func (n *Node) SetClusterLeader(leaderId string) {
	if n.leaderId != leaderId {
		log.Printf("FOLLOWING CLUSTER LEADER: %s\n", leaderId)
		n.leaderId = leaderId
	}
}

func (n Node) GetState() FsmState {
	return n.fsmState.getStateCopy()
}

func (n Node) GetQuorum() int {
	totalClusterNodes := len(n.peers) + 1 // add self
	return int(math.Ceil(float64(totalClusterNodes) / 2.0))
}

func (n *Node) SetRole(newRole Role) {
	n.role = newRole
}
func (n Node) isRole(role Role) bool {
	return n.role == role
}
func (n Node) IsFollower() bool {
	return n.isRole(Follower)
}
func (n Node) IsCandidate() bool {
	return n.isRole(Candidate)
}
func (n Node) IsLeader() bool {
	return n.isRole(Leader)
}

func (n Node) VoteForSelf() {
	n.SetVotedFor(n.id)
}

func (n Node) CanGrantVote(id string, lastLogIndex uint32, lastLogTerm uint32) bool {
	currentLogIndex := n.GetLastLogIndex()
	currentLogTerm := n.GetLastLogTerm()
	votedFor := n.votedFor

	voteAvailable := votedFor == "" || votedFor == id
	candidateUpToDate := currentLogTerm <= lastLogTerm && currentLogIndex <= lastLogIndex
	if voteAvailable && candidateUpToDate {
		return true
	}
	return false
}

func (n Node) Broadcast(fn func(p Peer)) {
	var wg sync.WaitGroup
	for _, peer := range n.peers {
		wg.Add(1)
		go func(p Peer, w *sync.WaitGroup) {
			defer w.Done()
			fn(p)
		}(peer, &wg)
	}
	wg.Wait()
}

func (n Node) GetRequestVoteInput() RequestVoteInput {
	return RequestVoteInput{
		CandidateId:  n.id,
		Term:         n.currentTerm,
		LastLogIndex: n.GetLastLogIndex(),
		LastLogTerm:  n.GetLastLogTerm(),
	}
}

func (n Node) GetAppendEntriesInput() AppendEntriesInput {
	return AppendEntriesInput{
		LeaderId:     n.id,
		Term:         n.currentTerm,
		PrevLogIndex: n.GetLastLogIndex(),
		PrevLogTerm:  n.GetLastLogTerm(),
		LeaderCommit: n.commitIndex,
	}
}
