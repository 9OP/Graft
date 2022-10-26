package receiver

import (
	"graft/app/domain/entity"
)

type service struct {
	repo Repository
}

func NewService(repo Repository) *service {
	return &service{repo}
}

func (s *service) AppendEntries(input *entity.AppendEntriesInput) (*entity.AppendEntriesOutput, error) {
	srv := s.repo
	state := srv.GetState()
	log := state.GetLogByIndex(input.PrevLogIndex)

	output := &entity.AppendEntriesOutput{
		Term:    state.CurrentTerm,
		Success: false,
	}

	if input.Term < state.CurrentTerm {
		return output, nil
	}

	if input.Term > state.CurrentTerm {
		srv.DowngradeFollower(input.Term, input.LeaderId)
	}

	srv.SetClusterLeader(input.LeaderId)
	srv.Heartbeat()

	if log.Term == input.PrevLogTerm {
		srv.AppendLogs(input.Entries)
		output.Success = true
	} else {
		srv.DeleteLogsFrom(input.PrevLogIndex)
	}

	if input.LeaderCommit > state.CommitIndex {
		lastLogIndex := state.GetLastLogIndex()
		if input.LeaderCommit > lastLogIndex {
			srv.SetCommitIndex(lastLogIndex)
		} else {
			srv.SetCommitIndex(input.LeaderCommit)
		}
	}

	return output, nil
}

func (s *service) RequestVote(input *entity.RequestVoteInput) (*entity.RequestVoteOutput, error) {
	srv := s.repo
	state := srv.GetState()

	output := &entity.RequestVoteOutput{
		Term:        state.CurrentTerm,
		VoteGranted: false,
	}

	if input.Term < state.CurrentTerm {
		return output, nil
	}

	srv.Heartbeat()

	if input.Term > state.CurrentTerm {
		srv.DowngradeFollower(input.Term, input.CandidateId)
	}

	if srv.GrantVote(input.CandidateId, input.LastLogIndex, input.LastLogTerm) {
		output.VoteGranted = true
	}

	return output, nil
}
