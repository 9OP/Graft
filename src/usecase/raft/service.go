package raft

import "graft/src/rpc"

type Service struct {
	repository Repository
}

func NewService(repo Repository) *Service {
	return &Service{
		repository: repo,
	}
}

func (service *Service) AppendEntries(input *rpc.AppendEntriesInput) (*rpc.AppendEntriesOutput, error) {
	state := service.repository.getState()
	state.Ping()

	output := &rpc.AppendEntriesOutput{}

	if input.Term > int32(state.CurrentTerm) {
		state.DowngradeToFollower(uint16(input.Term))
	}

	return output, nil
}

func (service *Service) RequestVote(input *rpc.RequestVoteInput) (*rpc.RequestVoteOutput, error) {
	//
}
