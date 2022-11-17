package receiver

import (
	"graft/pkg/domain"
	"graft/pkg/domain/state"
	"graft/pkg/utils"
)

type service struct {
	clusterNode *state.ClusterNode
}

func NewService(clusterNode *state.ClusterNode) *service {
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

	node.SetClusterLeader(input.LeaderId)
	node.Heartbeat()

	localPrevLog, err := node.Log(input.PrevLogIndex)
	if localPrevLog.Term == input.PrevLogTerm && err == nil {
		node.AppendLogs(input.PrevLogIndex, input.Entries...)
		output.Success = true
	} else {
		node.DeleteLogsFrom(input.PrevLogIndex)
	}

	if input.LeaderCommit > node.CommitIndex() {
		newIndex := utils.Min(node.LastLogIndex(), input.LeaderCommit)
		node.SetCommitIndex(newIndex)
	}

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

	isUpToDate := node.IsUpToDate(input.LastLogIndex, input.LastLogTerm)
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
	isUpToDate := node.IsUpToDate(input.LastLogIndex, input.LastLogTerm)

	if !hasLeader && isUpToDate {
		output.VoteGranted = true
	}

	return output, nil
}
