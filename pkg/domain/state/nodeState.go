package state

import (
	"math"

	"graft/pkg/domain"
	"graft/pkg/utils"
)

type nodeState struct {
	id       string
	leaderId string
	peers    domain.Peers
	role     domain.Role
	*fsmState
}

func NewNodeState(id string, peers domain.Peers, persistent *PersistentState) *nodeState {
	return &nodeState{
		id:       id,
		peers:    peers,
		role:     domain.Follower,
		fsmState: NewFsmState(persistent),
	}
}

func (n nodeState) Id() string {
	return n.id
}

func (n nodeState) LeaderId() string {
	return n.leaderId
}

func (n nodeState) HasLeader() bool {
	_, ok := n.peers[n.leaderId]
	return ok
}

func (n nodeState) Leader() domain.Peer {
	return n.peers[n.leaderId]
}

func (n nodeState) Role() domain.Role {
	return n.role
}

func (n nodeState) Peers() domain.Peers {
	return utils.CopyMap(n.peers)
}

func (n nodeState) WithInitializeLeader() nodeState {
	defaultNextIndex := n.LastLogIndex() // + 1
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

func (n nodeState) WithCurrentTerm(term uint32) nodeState {
	n.fsmState.PersistentState.currentTerm = term
	return n
}

func (n nodeState) WithVotedFor(vote string) nodeState {
	n.fsmState.PersistentState.votedFor = vote
	return n
}

func (n nodeState) WithRole(role domain.Role) nodeState {
	n.role = role
	return n
}

func (n nodeState) WithClusterLeader(leaderId string) nodeState {
	n.leaderId = leaderId
	return n
}

func (n nodeState) Quorum() int {
	totalPeers := len(n.peers) + 1
	return int(math.Ceil(float64(totalPeers) / 2.0))
}

func (n nodeState) IsLeader() bool {
	return n.id == n.leaderId
}

// Raft determines which of two logs is more up-to-date
// by comparing the index and term of the last entries in the
// logs. If the logs have last entries with different terms, then
// the log with the later term is more up-to-date. If the logs
// end with the same term, then whichever log is longer is
// more up-to-date.
// IsUpToDate from the caller point of view
func (n nodeState) IsUpToDate(lastLogIndex uint32, lastLogTerm uint32) bool {
	log := n.LastLog()

	if log.Term == lastLogTerm {
		return lastLogIndex >= n.LastLogIndex()
	}

	return lastLogTerm >= log.Term
}

func (n nodeState) CanGrantVote(peerId string, lastLogIndex uint32, lastLogTerm uint32) bool {
	// Unknown peer id
	if _, ok := n.peers[peerId]; !ok {
		return false
	}

	votedFor := n.VotedFor()
	voteAvailable := votedFor == "" || votedFor == peerId
	candidateUpToDate := n.IsUpToDate(lastLogIndex, lastLogTerm)

	return voteAvailable && candidateUpToDate
}

func (n *nodeState) AppendEntriesInput(peerId string) domain.AppendEntriesInput {
	matchIndex := n.MatchIndexForPeer(peerId)
	nextIndex := n.NextIndexForPeer(peerId)

	var prevLogTerm uint32
	var prevLogIndex uint32
	var entries []domain.LogEntry

	// No need to send new entries when peer
	// Already has matching log for lastLogIndex
	if matchIndex == n.LastLogIndex() {
		entries = []domain.LogEntry{}
		prevLogIndex = n.LastLogIndex()
		prevLog, _ := n.Log(prevLogIndex)
		prevLogTerm = prevLog.Term
	} else {
		entries = n.LogsFrom(nextIndex + 1)
		prevLogIndex = nextIndex
		prevLog, _ := n.Log(prevLogIndex)
		prevLogTerm = prevLog.Term
	}

	return domain.AppendEntriesInput{
		LeaderId:     n.Id(),
		Term:         n.CurrentTerm(),
		LeaderCommit: n.CommitIndex(),
		PrevLogIndex: prevLogIndex,
		PrevLogTerm:  prevLogTerm,
		Entries:      entries,
	}
}

func (n nodeState) RequestVoteInput() domain.RequestVoteInput {
	return domain.RequestVoteInput{
		CandidateId:  n.Id(),
		Term:         n.CurrentTerm(),
		LastLogIndex: n.LastLogIndex(),
		LastLogTerm:  n.LastLog().Term,
	}
}

func (n nodeState) ComputeNewCommitIndex() uint32 {
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
			if log, err := n.Log(N); err == nil {
				if log.Term == n.CurrentTerm() {
					return N
				}
			}
		}
	}

	return commitIndex
}
