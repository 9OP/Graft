package lib

import (
	"sync/atomic"

	"graft/pkg/domain"
)

type Client interface {
	AppendEntries(peer domain.Peer, input *domain.AppendEntriesInput) (*domain.AppendEntriesOutput, error)
	RequestVote(peer domain.Peer, input *domain.RequestVoteInput) (*domain.RequestVoteOutput, error)
	PreVote(peer domain.Peer, input *domain.RequestVoteInput) (*domain.RequestVoteOutput, error)
}

func PreVote(node *domain.Node, client Client) bool {
	input := node.RequestVoteInput()
	quorum := node.Quorum()
	var prevotesGranted uint32 = 1 // vote for self

	gatherPreVote := func(p domain.Peer) {
		if res, err := client.PreVote(p, &input); err == nil {
			if res.VoteGranted {
				atomic.AddUint32(&prevotesGranted, 1)
			}
		}
	}
	node.Broadcast(gatherPreVote, domain.BroadcastActive)

	quorumReached := int(prevotesGranted) >= quorum
	return quorumReached
}

func RequestVote(node *domain.Node, client Client) bool {
	input := node.RequestVoteInput()
	quorum := node.Quorum()
	var votesGranted uint32 = 1 // vote for self

	gatherVote := func(p domain.Peer) {
		if res, err := client.RequestVote(p, &input); err == nil {
			if res.Term > node.CurrentTerm() {
				node.DowngradeFollower(res.Term)
				return
			}
			if res.VoteGranted {
				atomic.AddUint32(&votesGranted, 1)
			}
		}
	}
	node.Broadcast(gatherVote, domain.BroadcastActive)

	isCandidate := node.Role() == domain.Candidate
	quorumReached := int(votesGranted) >= quorum
	return isCandidate && quorumReached
}

func AppendEntries(node *domain.Node, client Client) bool {
	quorum := node.Quorum()
	var peersAlive uint32 = 1 // self

	synchroniseLogsRoutine := func(p domain.Peer) {
		input := node.AppendEntriesInput(p.Id)
		if res, err := client.AppendEntries(p, &input); err == nil {
			atomic.AddUint32(&peersAlive, 1)
			if res.Term > node.CurrentTerm() {
				node.DowngradeFollower(res.Term)
				return
			}
			if res.Success {
				// newPeerLastLogIndex is always the sum of len(entries) and prevLogIndex
				// Even if peer was send logs it already has, because prevLogIndex is the number
				// of logs already contained at least in peer, and len(entries) is the additional
				// entries accepted.
				newPeerLastLogIndex := input.PrevLogIndex + uint32(len(input.Entries))
				node.SetNextMatchIndex(p.Id, newPeerLastLogIndex)
			} else {
				node.DecrementNextIndex(p.Id)
			}
		}
	}
	node.Broadcast(synchroniseLogsRoutine, domain.BroadcastAll)
	quorumReached := int(peersAlive) >= quorum
	return quorumReached
}
