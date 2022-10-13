package main

import (
	"context"
	rpc "graft/src/api/graft_rpc"
)

type Service struct {
	rpc.UnimplementedRpcServer
}

func (s *Service) AppendEntries(ctx context.Context, entries *rpc.AppendEntriesInput) (*rpc.AppendEntriesOutput, error) {
	state := ctx.Value("graft_server_state").(*ServerState)
	state.Ping()

	output := &rpc.AppendEntriesOutput{}

	if entries.Term > int32(state.CurrentTerm) {
		state.DowngradeToFollower(uint16(entries.Term))
	}

	return output, nil
}

func (s *Service) RequestVote(ctx context.Context, vote *rpc.RequestVoteInput) (*rpc.RequestVoteOutput, error) {
	state := ctx.Value("graft_server_state").(*ServerState)
	state.Ping() // retard receiver to raise to candidate

	output := &rpc.RequestVoteOutput{
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
