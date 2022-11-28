package rpc

import (
	"graft/pkg/domain"
	"graft/pkg/services/lib"
)

type UseCase interface {
	AppendEntries(input *domain.AppendEntriesInput) (*domain.AppendEntriesOutput, error)
	RequestVote(input *domain.RequestVoteInput) (*domain.RequestVoteOutput, error)
	PreVote(input *domain.RequestVoteInput) (*domain.RequestVoteOutput, error)
	Execute(input *domain.ExecuteInput) (*domain.ExecuteOutput, error)
	Configuration() (*domain.ClusterConfiguration, error)
	LeadershipTransfer() error
	Shutdown()
	Ping() error
}

type client interface {
	lib.Client
	Ping(target string) error
}
