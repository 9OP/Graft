package entity

import (
	"bytes"
	"math"
	"os/exec"
	"sync"

	log "github.com/sirupsen/logrus"
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
	if n.LeaderId != leaderId {
		n.LeaderId = leaderId
		if leaderId == n.Id {
			log.Info("LEADER")
		} else {
			log.Infof("FOLLOW %s", leaderId)
		}
	}
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

	// index of log entry immediately preceding new ones
	prevLogIndex := stateCopy.NextIndex[peerId]
	prevLogTerm := stateCopy.GetLogByIndex(prevLogIndex).Term
	if prevLogIndex > stateCopy.GetLastLogIndex() {
		prevLogTerm = stateCopy.CurrentTerm
	}

	newLogs := stateCopy.GetLogsFromIndex(prevLogIndex + 1)
	entries := make([]LogEntry, len(newLogs))
	copy(entries, newLogs)

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

func (n *Node) ComputeNewCommitIndex() uint32 {
	/*
		Compute new commitIndex N such that:
			- N > commitIndex,
			- a majority of matchIndex[i] ≥ N
			- log[N].term == currentTerm:

	*/
	n.mu.RLock()
	defer n.mu.RUnlock()

	lastLogIndex := n.GetLastLogIndex() // Upper value of N
	commitIndex := n.CommitIndex        // Lower value of N
	quorum := n.GetQuorum()

	for N := lastLogIndex; N > commitIndex; N-- {
		// Get a majority for which matchIndex >= n
		count := 1 // count self
		for _, matchIndex := range n.MatchIndex {
			if matchIndex >= N {
				count += 1
			}
		}

		if count >= quorum {
			log := n.GetLogByIndex(N)
			if log.Term == n.CurrentTerm {
				return N
			}
		}
	}
	return commitIndex
}

func (n *Node) ApplyLogs() {
	// thread-safe access
	stateCopy := n.FsmState.GetCopy()

	commitIndex := stateCopy.CommitIndex
	lastApplied := stateCopy.LastApplied

	for commitIndex > lastApplied {
		lastApplied += 1 // increment local
		log := stateCopy.GetLogByIndex(lastApplied)
		res := n.Exec(log.Value)
		if log.C != nil {
			log.C <- res
		}
	}

	// Commit last applied
	n.SetLastApplied(commitIndex)
}

func (n *Node) Exec(entry string) interface{} {
	log.Debug("EXECUTE: ", entry)

	cmd := exec.Command(entry)
	var outb, errb bytes.Buffer
	cmd.Stdout = &outb
	cmd.Stderr = &errb

	err := cmd.Run()

	if err != nil {
		return errb.Bytes()
	}

	return outb.Bytes()
}
