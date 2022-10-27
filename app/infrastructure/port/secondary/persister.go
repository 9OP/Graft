package secondaryPort

import (
	"graft/app/domain/entity"
	adapter "graft/app/infrastructure/adapter/secondary"
)

type persisterPort struct {
	adapter  adapter.UseCaseJsonPersisterAdapter
	location string
}

func NewPersisterPort(location string, adapter adapter.UseCaseJsonPersisterAdapter) *persisterPort {
	return &persisterPort{adapter, location}
}

func (p *persisterPort) Load() (*entity.Persistent, error) {
	return p.adapter.Load(p.location)
}

func (p *persisterPort) Save(state *entity.Persistent) error {
	return p.adapter.Save(state, p.location)
}
