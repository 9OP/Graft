package entity

import (
	"math"
	"sync"
)

type Node struct {
	Id       string
	LeaderId string
	Peers    Peers
	Role     Role
	mu       sync.RWMutex
	*FsmState
}

func NewNode(id string, peers Peers, persistent *Persistent) *Node {
	return &Node{
		Id:       id,
		Peers:    peers,
		Role:     Follower,
		FsmState: NewFsmState(persistent),
	}
}

func (n *Node) SetRole(role Role) {
	n.mu.Lock()
	defer n.mu.Unlock()
	n.Role = role
}

func (n *Node) IsRole(role Role) bool {
	n.mu.RLock()
	defer n.mu.RUnlock()
	return n.Role == role
}

func (n *Node) SetClusterLeader(leaderId string) {
	n.mu.Lock()
	defer n.mu.Unlock()
	n.LeaderId = leaderId
}

func (n *Node) GetQuorum() int {
	n.mu.RLock()
	defer n.mu.RUnlock()
	totalPeers := len(n.Peers) + 1 // add self
	return int(math.Ceil(float64(totalPeers) / 2.0))
}

func (n *Node) IsLeader() bool {
	n.mu.RLock()
	defer n.mu.RUnlock()
	return n.LeaderId == n.Id
}

func (n *Node) CanGrantVote(peerId string, lastLogIndex uint32, lastLogTerm uint32) bool {
	// thread-safe access
	persistentCopy := n.FsmState.Persistent.GetCopy()

	currentLogIndex := persistentCopy.GetLastLogIndex()
	currentLogTerm := persistentCopy.GetLastLogTerm()
	votedFor := persistentCopy.VotedFor

	voteAvailable := votedFor == "" || votedFor == peerId
	candidateUpToDate := currentLogTerm <= lastLogTerm && currentLogIndex <= lastLogIndex

	if voteAvailable && candidateUpToDate {
		return true
	}
	return false
}

func (n *Node) GetAppendEntriesInput(peerId string) *AppendEntriesInput {
	// thread-safe access
	stateCopy := n.FsmState.GetCopy()

	prevLogIndex := stateCopy.NextIndex[peerId]
	prevLogTerm := stateCopy.GetLogByIndex(prevLogIndex).Term
	if prevLogTerm == 0 {
		prevLogTerm = stateCopy.GetLastLogTerm()
	}

	logs := stateCopy.GetLogsFromIndex(prevLogIndex)
	entries := make([]string, len(logs))
	for _, log := range logs {
		entries = append(entries, log.Value)
	}

	return &AppendEntriesInput{
		LeaderId:     n.Id,
		Term:         stateCopy.CurrentTerm,
		PrevLogIndex: prevLogIndex,
		PrevLogTerm:  prevLogTerm,
		Entries:      entries,
		LeaderCommit: stateCopy.CommitIndex,
	}
}
func (n *Node) GetRequestVoteInput() *RequestVoteInput {
	// thread-safe access
	stateCopy := n.FsmState.GetCopy()

	return &RequestVoteInput{
		CandidateId:  n.Id,
		Term:         stateCopy.CurrentTerm,
		LastLogIndex: stateCopy.GetLastLogIndex(),
		LastLogTerm:  stateCopy.GetLastLogTerm(),
	}
}
