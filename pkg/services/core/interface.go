package core

import (
	"graft/pkg/domain"
)

type client interface {
	AppendEntries(peer domain.Peer, input *domain.AppendEntriesInput) (*domain.AppendEntriesOutput, error)
	RequestVote(peer domain.Peer, input *domain.RequestVoteInput) (*domain.RequestVoteOutput, error)
	PreVote(peer domain.Peer, input *domain.RequestVoteInput) (*domain.RequestVoteOutput, error)
}

type persister interface {
	Load() (domain.PersistentState, error)
	Save(domain.PersistentState) error
}

type UseCase interface {
	Run()
}
