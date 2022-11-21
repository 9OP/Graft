package rpc

import (
	"graft/pkg/domain"
)

type service struct {
	clusterNode *domain.Node
}

func NewService(clusterNode *domain.Node) *service {
	return &service{clusterNode}
}

func (s *service) AppendEntries(input *domain.AppendEntriesInput) (*domain.AppendEntriesOutput, error) {
	node := s.clusterNode

	output := &domain.AppendEntriesOutput{
		Term:    node.CurrentTerm(),
		Success: false,
	}
	if input.Term < node.CurrentTerm() {
		return output, nil
	}
	if input.Term > node.CurrentTerm() {
		node.DowngradeFollower(input.Term)
	}

	node.SetLeader(input.LeaderId)
	node.Heartbeat()

	localPrevLog, err := node.Log(input.PrevLogIndex)
	if localPrevLog.Term == input.PrevLogTerm && err == nil {
		node.AppendLogs(input.PrevLogIndex, input.Entries...)
		output.Success = true
	} else {
		node.DeleteLogsFrom(input.PrevLogIndex)
	}

	node.UpdateLeaderCommitIndex(input.LeaderCommit)

	return output, nil
}

func (s *service) RequestVote(input *domain.RequestVoteInput) (*domain.RequestVoteOutput, error) {
	node := s.clusterNode

	output := &domain.RequestVoteOutput{
		Term:        node.CurrentTerm(),
		VoteGranted: false,
	}
	if input.Term < node.CurrentTerm() {
		return output, nil
	}
	if input.Term > node.CurrentTerm() {
		node.DowngradeFollower(input.Term)
	}

	isUpToDate := node.IsLogUpToDate(input.LastLogIndex, input.LastLogTerm)
	canGrantVote := node.CanGrantVote(input.CandidateId)

	if canGrantVote && isUpToDate {
		node.Heartbeat()
		node.GrantVote(input.CandidateId)
		output.VoteGranted = true
	}

	return output, nil
}

func (s *service) PreVote(input *domain.RequestVoteInput) (*domain.RequestVoteOutput, error) {
	node := s.clusterNode

	output := &domain.RequestVoteOutput{
		Term:        node.CurrentTerm(),
		VoteGranted: false,
	}

	hasLeader := node.HasLeader()
	isUpToDate := node.IsLogUpToDate(input.LastLogIndex, input.LastLogTerm)

	if !hasLeader && isUpToDate {
		output.VoteGranted = true
	}

	return output, nil
}