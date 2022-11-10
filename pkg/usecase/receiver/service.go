package receiver

import (
	"fmt"

	utils "graft/pkg/domain"
	"graft/pkg/domain/entity"
	domain "graft/pkg/domain/service"
)

type service struct {
	clusterNode *domain.ClusterNode
}

func NewService(clusterNode *domain.ClusterNode) *service {
	return &service{clusterNode}
}

func (s *service) AppendEntries(input *entity.AppendEntriesInput) (*entity.AppendEntriesOutput, error) {
	node := s.clusterNode

	output := &entity.AppendEntriesOutput{
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

	localPrevLog, err := node.MachineLog(input.PrevLogIndex)
	if localPrevLog.Term == input.PrevLogTerm && err == nil {
		node.AppendLogs(input.PrevLogIndex, input.Entries...)
		output.Success = true
	} else {
		fmt.Println("delete logs")
		node.DeleteLogsFrom(input.PrevLogIndex)
	}

	if input.LeaderCommit >= node.CommitIndex() {
		newIndex := utils.Min(node.LastLogIndex(), input.LeaderCommit)
		node.SetCommitIndex(newIndex)
	}

	return output, nil
}

func (s *service) RequestVote(input *entity.RequestVoteInput) (*entity.RequestVoteOutput, error) {
	node := s.clusterNode

	output := &entity.RequestVoteOutput{
		Term:        node.CurrentTerm(),
		VoteGranted: false,
	}
	if input.Term < node.CurrentTerm() {
		return output, nil
	}
	if input.Term > node.CurrentTerm() {
		node.DowngradeFollower(input.Term)
	}

	if node.CanGrantVote(input.CandidateId, input.LastLogIndex, input.LastLogTerm) {
		node.Heartbeat()
		node.GrantVote(input.CandidateId)
		output.VoteGranted = true
	}

	return output, nil
}

func (s *service) PreVote(input *entity.RequestVoteInput) (*entity.RequestVoteOutput, error) {
	node := s.clusterNode

	output := &entity.RequestVoteOutput{
		Term:        node.CurrentTerm(),
		VoteGranted: false,
	}

	if node.IsLogUpToDate(input.LastLogIndex, input.LastLogTerm) {
		output.VoteGranted = true
	}

	return output, nil
}
