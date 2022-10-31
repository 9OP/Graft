package cluster

import "graft/app/domain/entity"

type repository interface {
	Execute(entry string) chan interface{}
	IsLeader() bool
	GetLeader() entity.Peer
}

type UseCase interface {
	// Writes with ExecuteCommand (leader only)
	ExecuteCommand(command string) (interface{}, error)
	// Reads with ExecuteQuery (leader + followers)
	ExecuteQuery(query string) (interface{}, error)
}
