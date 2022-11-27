package core

import (
	"graft/pkg/domain"
)

type persister interface {
	Load() (domain.PersistentState, error)
	Save(domain.PersistentState) error
}

type UseCase interface {
	Run()
}
