package receiver

import (
	"graft/app/domain/entity"
)

type service struct {
	repository
}

func NewService(repository repository) *service {
	return &service{repository}
}

func (s *service) AppendEntries(input *entity.AppendEntriesInput) (*entity.AppendEntriesOutput, error) {
	state := s.GetState()
	log := state.GetLogByIndex(input.PrevLogIndex)

	output := &entity.AppendEntriesOutput{
		Term:    state.CurrentTerm,
		Success: false,
	}

	if input.Term < state.CurrentTerm {
		return output, nil
	}

	if input.Term > state.CurrentTerm {
		s.DowngradeFollower(input.Term, input.LeaderId)
	}

	s.SetClusterLeader(input.LeaderId)
	s.Heartbeat()

	if log.Term == input.PrevLogTerm {
		s.AppendLogs(input.Entries)
		output.Success = true
	} else {
		s.DeleteLogsFrom(input.PrevLogIndex)
	}

	if input.LeaderCommit > state.CommitIndex {
		lastLogIndex := state.GetLastLogIndex()
		if input.LeaderCommit > lastLogIndex {
			s.SetCommitIndex(lastLogIndex)
		} else {
			s.SetCommitIndex(input.LeaderCommit)
		}
	}

	return output, nil
}

func (s *service) RequestVote(input *entity.RequestVoteInput) (*entity.RequestVoteOutput, error) {
	state := s.GetState()

	output := &entity.RequestVoteOutput{
		Term:        state.CurrentTerm,
		VoteGranted: false,
	}

	if input.Term < state.CurrentTerm {
		return output, nil
	}

	s.Heartbeat()

	if input.Term > state.CurrentTerm {
		s.DowngradeFollower(input.Term, input.CandidateId)
	}

	if s.GrantVote(input.CandidateId, input.LastLogIndex, input.LastLogTerm) {
		output.VoteGranted = true
	}

	return output, nil
}
