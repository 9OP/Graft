package receiver

import (
	"graft/pkg/domain"
)

type UseCase interface {
	AppendEntries(input *domain.AppendEntriesInput) (*domain.AppendEntriesOutput, error)
	RequestVote(input *domain.RequestVoteInput) (*domain.RequestVoteOutput, error)
	PreVote(input *domain.RequestVoteInput) (*domain.RequestVoteOutput, error)
}
