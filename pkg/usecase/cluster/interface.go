package cluster

import "graft/pkg/domain/entity"

type repository interface {
	ExecuteCommand(command string) chan entity.EvalResult
	ExecuteQuery(query string) chan entity.EvalResult
	IsLeader() bool
	Leader() entity.Peer
}

type UseCase interface {
	ExecuteCommand(command string) ([]byte, error)
	ExecuteQuery(query string, weakConsistency bool) ([]byte, error)
}
