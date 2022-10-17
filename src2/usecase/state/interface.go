package state

import "graft/src2/entity"

type Repository interface {
	Load(location string) (*entity.PersistentState, error)
	Save(location string, state *entity.PersistentState) error
}

type UseCase interface {
	LoadState() (*entity.PersistentState, error)
	SaveState(state *entity.PersistentState) error
}
