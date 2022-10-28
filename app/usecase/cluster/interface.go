package cluster

import (
	"graft/app/domain/entity"
)

type repository interface {
	AppendEntries(peer entity.Peer, input *entity.AppendEntriesInput) (*entity.AppendEntriesOutput, error)
}

type UseCase interface {
	NewCmd(cmd string) (interface{}, error)
}
