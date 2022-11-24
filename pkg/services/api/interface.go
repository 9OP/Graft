package api

import "graft/pkg/domain"

type UseCase interface {
	ExecuteCommand(command domain.ExecuteInput) ([]byte, error)
	ExecuteQuery(query string, weakConsistency bool) ([]byte, error)
}
