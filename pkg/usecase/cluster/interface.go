package cluster

import "graft/pkg/domain/entity"

type repository interface {
	ExecuteCommand(command string) chan interface{}
	ExecuteQuery(query string) chan interface{}
	IsLeader() bool
	GetLeader() entity.Peer
}

type UseCase interface {
	// Writes with ExecuteCommand (leader only)
	ExecuteCommand(command string) (interface{}, error)
	// Reads with ExecuteQuery (leader + followers)
	ExecuteQuery(query string, weakConsistency bool) (interface{}, error)
}
