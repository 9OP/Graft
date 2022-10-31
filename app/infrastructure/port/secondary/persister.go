package secondaryPort

import (
	"graft/app/domain/entity"
	adapter "graft/app/infrastructure/adapter/secondary"
	"sync"

	log "github.com/sirupsen/logrus"
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
	pst, err := p.adapter.Load(p.location)
	if err != nil {
		log.Warn("CANNOT LOAD PERSISTENT, USING DEFAULT")
		return entity.NewPersistent(), err
	}
	return pst, err
}

func (p *persisterPort) Save(state *entity.Persistent) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.adapter.Save(state, p.location)
}
