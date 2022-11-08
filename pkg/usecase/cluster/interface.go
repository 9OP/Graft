package cluster

import "graft/pkg/domain/entity"

type repository interface {
	ExecuteCommand(command string) chan entity.EvalResult
	ExecuteQuery(query string) chan entity.EvalResult
	IsLeader() bool
	Leader() entity.Peer
}

type UseCase interface {
	// Writes with ExecuteCommand (leader only)
	ExecuteCommand(command string) (interface{}, error)
	// Reads with ExecuteQuery (leader + followers)
	ExecuteQuery(query string, weakConsistency bool) (interface{}, error)
}
