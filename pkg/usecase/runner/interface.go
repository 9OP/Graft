package runner

import (
	"graft/pkg/domain/entity"
)

type repository interface {
	AppendEntries(peer entity.Peer, input *entity.AppendEntriesInput) (*entity.AppendEntriesOutput, error)
	RequestVote(peer entity.Peer, input *entity.RequestVoteInput) (*entity.RequestVoteOutput, error)
	PreVote(peer entity.Peer, input *entity.RequestVoteInput) (*entity.RequestVoteOutput, error)
}

type persister interface {
	Load() (*entity.PersistentState, error)
	Save(currentTerm uint32, votedFor string, machineLogs []entity.LogEntry) error
}

type UseCase interface {
	Run()
}
