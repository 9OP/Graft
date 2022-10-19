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
		srv.DowngradeFollower(input.Term)
	}

	srv.Heartbeat()
	srv.SetClusterLeader(input.LeaderId)
	state.LastLogIndex()
	state.LastLogTerm()

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
		srv.DowngradeFollower(input.Term)
	}

	if srv.GrantVote(input.CandidateId, input.LastLogIndex, input.LastLogTerm) {
		output.VoteGranted = true
	}

	return output, nil
}
