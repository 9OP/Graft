package receiver

import (
	"graft/pkg/domain/entity"
)

type UseCase interface {
	AppendEntries(input *entity.AppendEntriesInput) (*entity.AppendEntriesOutput, error)
	RequestVote(input *entity.RequestVoteInput) (*entity.RequestVoteOutput, error)
	PreVote(input *entity.RequestVoteInput) (*entity.RequestVoteOutput, error)
}
