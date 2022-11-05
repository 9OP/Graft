package receiver

import (
	utils "graft/pkg/domain"
	"graft/pkg/domain/entity"
)

type service struct {
	repository
}

func NewService(repository repository) *service {
	return &service{repository}
}

func (s *service) AppendEntries(input *entity.AppendEntriesInput) (*entity.AppendEntriesOutput, error) {
	state := s.GetState()

	output := &entity.AppendEntriesOutput{
		Term:    state.CurrentTerm(),
		Success: false,
	}
	if input.Term < state.CurrentTerm() {
		return output, nil
	}
	if input.Term > state.CurrentTerm() {
		s.DowngradeFollower(input.Term)
	}

	s.SetClusterLeader(input.LeaderId)
	s.Heartbeat()

	localPrevLog, err := state.MachineLog(input.PrevLogIndex)
	if localPrevLog.Term == input.PrevLogTerm && err == nil {
		s.AppendLogs(input.PrevLogIndex, input.Entries...)
		output.Success = true
	} else {
		s.DeleteLogsFrom(input.PrevLogIndex)
	}

	if input.LeaderCommit > state.CommitIndex() {
		s.SetCommitIndex(utils.Min(state.LastLogIndex(), input.LeaderCommit))
	}

	return output, nil
}

func (s *service) RequestVote(input *entity.RequestVoteInput) (*entity.RequestVoteOutput, error) {
	state := s.GetState()

	output := &entity.RequestVoteOutput{
		Term:        state.CurrentTerm(),
		VoteGranted: false,
	}
	if input.Term < state.CurrentTerm() {
		return output, nil
	}
	if input.Term > state.CurrentTerm() {
		s.DowngradeFollower(input.Term)
	}

	if state.CanGrantVote(input.CandidateId, input.LastLogIndex, input.LastLogTerm) {
		s.Heartbeat()
		s.GrantVote(input.CandidateId)
		output.VoteGranted = true
	}

	return output, nil
}