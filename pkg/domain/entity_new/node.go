package entitynew

import (
	"graft/pkg/domain/entity"
	"math"
)

type Node struct {
	id       string
	leaderId string
	peers    entity.Peers
	role     entity.Role
	FsmState
}

func NewNode(id string, peers entity.Peers, persistent Persistent) Node {
	return Node{
		id:       id,
		peers:    peers,
		role:     entity.Follower,
		FsmState: NewFsmState(persistent),
	}
}

func (n Node) Id() string {
	return n.id
}

func (n Node) LeaderId() string {
	return n.leaderId
}

func (n Node) Role() entity.Role {
	return n.role
}

func (n Node) WithInitializeLeader() Node {
	defaultNextIndex := n.LastLogIndex()
	nextIndex := make(peerIndex, len(n.peers))
	matchIndex := make(peerIndex, len(n.peers))
	for _, peer := range n.peers {
		nextIndex[peer.Id] = defaultNextIndex
		matchIndex[peer.Id] = 0
	}
	n.nextIndex = nextIndex
	n.matchIndex = matchIndex
	return n
}

func (n Node) WithRole(role entity.Role) Node {
	n.role = role
	return n
}

func (n Node) WithClusterLeader(leaderId string) Node {
	n.leaderId = leaderId
	return n
}

func (n Node) Quorum() int {
	totalPeers := len(n.peers) + 1
	return int(math.Ceil(float64(totalPeers) / 2.0))
}

func (n Node) IsLeader() bool {
	return n.Id() == n.LeaderId()
}

func (n Node) CanGrantVote(peerId string, lastLogIndex uint32, lastLogTerm uint32) bool {
	currentLogIndex := n.LastLogIndex()
	currentLogTerm := n.LastLog().Term
	votedFor := n.VotedFor()

	voteAvailable := votedFor == "" || votedFor == peerId
	candidateUpToDate := currentLogTerm <= lastLogTerm && currentLogIndex <= lastLogIndex

	if voteAvailable && candidateUpToDate {
		return true
	}
	return false
}

func (n Node) AppendEntriesInput(peerId string) entity.AppendEntriesInput {
	// index of log entry immediately preceding new ones
	prevLogIndex := n.NextIndexForPeer(peerId)
	prevLog, err := n.MachineLog(prevLogIndex)
	prevLogTerm := prevLog.Term

	if err != nil {
		prevLogTerm = n.CurrentTerm()
	}

	newEntries := n.MachineLogsFrom(prevLogIndex + 1)

	return entity.AppendEntriesInput{
		LeaderId:     n.Id(),
		Term:         n.CurrentTerm(),
		PrevLogIndex: prevLogIndex,
		PrevLogTerm:  prevLogTerm,
		Entries:      newEntries,
		LeaderCommit: n.CommitIndex(),
	}
}

func (n Node) RequestVoteInput() entity.RequestVoteInput {
	return entity.RequestVoteInput{
		CandidateId:  n.Id(),
		Term:         n.CurrentTerm(),
		LastLogIndex: n.LastLogIndex(),
		LastLogTerm:  n.LastLog().Term,
	}
}

func (n Node) ComputeNewCommitIndex() uint32 {
	/*
		Compute new commitIndex N such that:
			- N > commitIndex,
			- a majority of matchIndex[i] ≥ N
			- log[N].term == currentTerm:
	*/
	quorum := n.Quorum()
	lastLogIndex := n.LastLogIndex() // Upper value of N
	commitIndex := n.CommitIndex()   // Lower value of N
	matchIndex := n.MatchIndex()

	for N := lastLogIndex; N > commitIndex; N-- {
		// Get a majority for which matchIndex >= n
		count := 1 // count self
		for _, matchIndex := range matchIndex {
			if matchIndex >= N {
				count += 1
			}
		}
		if count >= quorum {
			if log, err := n.MachineLog(N); err == nil {
				if log.Term == n.CurrentTerm() {
					return N
				}
			}
		}
	}

	return commitIndex
}

// Move to service
func (n Node) ApplyLogs() {
	commitIndex := n.CommitIndex()
	for lastApplied := n.LastApplied(); commitIndex > lastApplied; {
		n.WithIncrementLastApplied()
		if _, err := n.MachineLog(lastApplied); err == nil {
			// res := n.Exec(log.Value)
			// if log.C != nil {
			// 	log.C <- res
			// }
		}
		// pointer switch
	}
}

// Move to service
// func (n *Node) Exec(entry string) interface{} {
// 	log.Debug("EXECUTE: ", entry)

// 	cmd := exec.Command(entry)
// 	var outb, errb bytes.Buffer
// 	cmd.Stdout = &outb
// 	cmd.Stderr = &errb

// 	err := cmd.Run()

// 	if err != nil {
// 		return errb.Bytes()
// 	}

// 	return outb.Bytes()
// }
