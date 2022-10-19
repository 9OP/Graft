package receiver

import (
	"context"
	"graft/app/rpc"
)

type Service struct {
	rpc.UnimplementedRpcServer
	repository Repository
}

func NewService(repo Repository) *Service {
	return &Service{repository: repo}
}

func (service *Service) AppendEntries(ctx context.Context, input *rpc.AppendEntriesInput) (*rpc.AppendEntriesOutput, error) {
	srv := service.repository
	state := srv.GetState()

	output := &rpc.AppendEntriesOutput{
		Term:    state.CurrentTerm,
		Success: false,
	}

	if input.Term < state.CurrentTerm {
		return output, nil
	}

	if input.Term > state.CurrentTerm {
		srv.DowngradeFollower(input.Term, input.LeaderId)
	}

	// Set cluster leader

	srv.Heartbeat()

	log := state.GetLogIndex(int(input.PrevLogIndex))

	// Refactor write to server directly

	if log.Term != input.PrevLogTerm {
		// srv.DeleteLogsFrom()
		state.DeleteLogFrom(int(input.PrevLogIndex))
	} else {
		output.Success = true
		// srv.AppendLogs()
		state.AppendLogs(input.Entries)
	}

	if input.LeaderCommit > state.CommitIndex {
		lastLogIndex := state.LastLogIndex()
		if input.LeaderCommit > lastLogIndex {
			state.CommitIndex = lastLogIndex
		} else {
			state.CommitIndex = input.LeaderCommit
		}
	}

	srv.SaveState()

	return output, nil
}

func (service *Service) RequestVote(ctx context.Context, input *rpc.RequestVoteInput) (*rpc.RequestVoteOutput, error) {
	srv := service.repository
	state := srv.GetState()

	output := &rpc.RequestVoteOutput{
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
