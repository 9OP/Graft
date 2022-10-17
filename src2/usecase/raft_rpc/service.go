package raft_rpc

import "graft/src2/rpc"

type Service struct {
	repository Repository
}

func NewService(repo Repository) *Service {
	return &Service{
		repository: repo,
	}
}

func (service *Service) AppendEntries(input *rpc.AppendEntriesInput) (*rpc.AppendEntriesOutput, error) {
	// state := service.repository.getState()
	// state.Ping()

	// output := &rpc.AppendEntriesOutput{}

	// if input.Term > int32(state.CurrentTerm) {
	// 	state.DowngradeToFollower(uint16(input.Term))
	// }

	// return output, nil
	return nil, nil
}

func (service *Service) RequestVote(input *rpc.RequestVoteInput) (*rpc.RequestVoteOutput, error) {
	return nil, nil
}
