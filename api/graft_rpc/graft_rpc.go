package graft_rpc

import (
	"context"
	"graft/models"
)

type Service struct {
	UnimplementedRpcServer
}

func (s *Service) AppendEntries(ctx context.Context, entries *AppendEntriesInput) (*AppendEntriesOutput, error) {
	state := ctx.Value("graft_server_state").(*models.ServerState)
	state.Ping()

	output := &AppendEntriesOutput{}

	if entries.Term > int32(state.CurrentTerm) {
		state.DowngradeToFollower(uint16(entries.Term))
	}

	return output, nil
}

func (s *Service) RequestVote(ctx context.Context, vote *RequestVoteInput) (*RequestVoteOutput, error) {
	state := ctx.Value("graft_server_state").(*models.ServerState)
	state.Ping()

	output := &RequestVoteOutput{
		Term:        int32(state.CurrentTerm),
		VoteGranted: false,
	}

	if vote.Term > int32(state.CurrentTerm) {
		state.DowngradeToFollower(uint16(vote.Term))
	}

	if vote.Term < int32(state.CurrentTerm) {
		return output, nil
	}

	if state.CanVote(vote.CandidateId, int(vote.LastLogIndex), int(vote.LastLogTerm)) {
		state.Vote(vote.CandidateId)
		output.VoteGranted = true
	}

	return output, nil
}
