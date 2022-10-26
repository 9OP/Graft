package entity

import (
	"log"
	"math"
)

// Write here every function that does not need to send external signal

// DDD aggregate over role, peers, fsmState(persistent)
type Node struct {
	id        string
	leaderId  string
	Peers     Peers
	role      Role
	*fsmState // ensure that we can call public method of private fsmState
}

func NewNode(id string, peers Peers) *Node {
	return &Node{
		id:       id,
		Peers:    peers,
		fsmState: NewFsmState(),
	}
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
	totalClusterNodes := len(n.Peers) + 1 // add self
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

func (n Node) GrantVote(id string, lastLogIndex uint32, lastLogTerm uint32) bool {
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

// func (s *Server) RequestVoteInput() *rpc.RequestVoteInput {
// 	state := s.GetState()
// 	return &rpc.RequestVoteInput{
// 		CandidateId:  s.Id,
// 		Term:         state.CurrentTerm,
// 		LastLogIndex: state.LastLogIndex(),
// 		LastLogTerm:  state.LastLogTerm(),
// 	}
// }

// func (s *Server) AppendEntriesInput() *rpc.AppendEntriesInput {
// 	state := s.GetState()
// 	return &rpc.AppendEntriesInput{
// 		LeaderId: s.Id,
// 		Term:     state.CurrentTerm,
// 	}
// }
