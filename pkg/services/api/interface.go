package api

import "graft/pkg/domain"

type UseCase interface {
	ExecuteCommand(command domain.ApiCommand) ([]byte, error)
	ExecuteQuery(query string, weakConsistency bool) ([]byte, error)
}
