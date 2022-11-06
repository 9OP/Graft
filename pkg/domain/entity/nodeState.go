package entity

import (
	utils "graft/pkg/domain"
	"math"
)

type NodeState struct {
	id       string
	leaderId string
	peers    Peers
	role     Role
	*FsmState
}

func NewNodeState(id string, peers Peers, persistent *PersistentState) NodeState {
	fsmState := NewFsmState(persistent)
	return NodeState{
		id:       id,
		peers:    peers,
		role:     Follower,
		FsmState: &fsmState,
	}
}

func (n NodeState) Id() string {
	return n.id
}

func (n NodeState) LeaderId() string {
	return n.leaderId
}

func (n NodeState) Leader() Peer {
	leaderId := n.leaderId
	leader := n.peers[leaderId]
	return leader
}

func (n NodeState) Role() Role {
	return n.role
}

func (n NodeState) Peers() Peers {
	return utils.CopyMap(n.peers)
}

func (n NodeState) WithInitializeLeader() NodeState {
	defaultNextIndex := n.LastLogIndex() + 1
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

func (n NodeState) WithCurrentTerm(term uint32) NodeState {
	n.FsmState.PersistentState.currentTerm = term
	return n
}
func (n NodeState) WithVotedFor(vote string) NodeState {
	n.FsmState.PersistentState.votedFor = vote
	return n
}

func (n NodeState) WithRole(role Role) NodeState {
	n.role = role
	return n
}

func (n NodeState) WithClusterLeader(leaderId string) NodeState {
	n.leaderId = leaderId
	return n
}

func (n NodeState) Quorum() int {
	totalPeers := len(n.peers) + 1
	return int(math.Ceil(float64(totalPeers) / 2.0))
}

func (n NodeState) IsLeader() bool {
	return n.Id() == n.LeaderId()
}

func (n NodeState) CanGrantVote(peerId string, lastLogIndex uint32, lastLogTerm uint32) bool {
	// Unknown peer id
	if _, ok := n.peers[peerId]; !ok {
		return false
	}

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

func (n NodeState) AppendEntriesInput(peerId string) AppendEntriesInput {
	matchIndex := n.MatchIndexForPeer(peerId)
	nextIndex := n.NextIndexForPeer(peerId)

	var prevLogTerm uint32
	var prevLogIndex uint32
	var entries []LogEntry

	// No need to send new entries when peer
	// Already has matching log for lastLogIndex
	if matchIndex == n.LastLogIndex() {
		entries = []LogEntry{}
		prevLogIndex = n.LastLogIndex()
		prevLog, _ := n.MachineLog(prevLogIndex)
		prevLogTerm = prevLog.Term
	} else {
		entries = n.MachineLogsFrom(nextIndex)
		prevLogIndex = nextIndex - 1
		prevLog, _ := n.MachineLog(prevLogIndex)
		prevLogTerm = prevLog.Term
	}

	return AppendEntriesInput{
		LeaderId:     n.Id(),
		Term:         n.CurrentTerm(),
		LeaderCommit: n.CommitIndex(),
		PrevLogIndex: prevLogIndex,
		PrevLogTerm:  prevLogTerm,
		Entries:      entries,
	}
}

func (n NodeState) RequestVoteInput() RequestVoteInput {
	return RequestVoteInput{
		CandidateId:  n.Id(),
		Term:         n.CurrentTerm(),
		LastLogIndex: n.LastLogIndex(),
		LastLogTerm:  n.LastLog().Term,
	}
}

func (n NodeState) ComputeNewCommitIndex() uint32 {
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
