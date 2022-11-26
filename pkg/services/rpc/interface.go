package rpc

import (
	"graft/pkg/domain"
)

type UseCase interface {
	AppendEntries(input *domain.AppendEntriesInput) (*domain.AppendEntriesOutput, error)
	RequestVote(input *domain.RequestVoteInput) (*domain.RequestVoteOutput, error)
	PreVote(input *domain.RequestVoteInput) (*domain.RequestVoteOutput, error)
	//
	Execute(input *domain.ExecuteInput) (*domain.ExecuteOutput, error)
	Configuration() (*domain.ClusterConfiguration, error)
	LeadershipTransfer() error
	Shutdown()
	Ping() error
}

type client interface {
	PreVote(peer domain.Peer, input *domain.RequestVoteInput) (*domain.RequestVoteOutput, error)
	Ping(peer domain.Peer) error
}
