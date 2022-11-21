package runner

import (
	"graft/pkg/domain"
)

type repository interface {
	AppendEntries(peer domain.Peer, input *domain.AppendEntriesInput) (*domain.AppendEntriesOutput, error)
	RequestVote(peer domain.Peer, input *domain.RequestVoteInput) (*domain.RequestVoteOutput, error)
	PreVote(peer domain.Peer, input *domain.RequestVoteInput) (*domain.RequestVoteOutput, error)
}

type persister interface {
	Load() (domain.PersistentState, error)
	Save(currentTerm uint32, votedFor string, machineLogs []domain.LogEntry) error
}

type UseCase interface {
	Run()
}
