package persister

import "graft/app/entity"

type Repository interface {
	Load(location string) (*entity.Persistent, error)
	Save(location string, state *entity.Persistent) error
}

type UseCase interface {
	LoadState() (*entity.Persistent, error)
	SaveState(state *entity.Persistent) error
}
