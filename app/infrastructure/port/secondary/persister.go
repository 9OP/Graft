package secondaryPort

import (
	"graft/app/domain/entity"
	adapter "graft/app/infrastructure/adapter/secondary"
	"sync"
)

type persisterPort struct {
	adapter  adapter.UseCaseJsonPersisterAdapter
	location string
	mu       sync.Mutex
}

func NewPersisterPort(location string, adapter adapter.UseCaseJsonPersisterAdapter) *persisterPort {
	return &persisterPort{
		adapter:  adapter,
		location: location,
	}
}

func (p *persisterPort) Load() (*entity.Persistent, error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.adapter.Load(p.location)
}

func (p *persisterPort) Save(state *entity.Persistent) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.adapter.Save(state, p.location)
}
