package secondaryPort

import (
	"graft/pkg/domain"

	log "github.com/sirupsen/logrus"
)

type AdapterPersister interface {
	Load(location string) (*domain.PersistentState, error)
	Save(location string, currentTerm uint32, votedFor string, machineLogs []domain.LogEntry) error
}

type persisterPort struct {
	adapter  AdapterPersister
	location string
}

func NewPersisterPort(location string, adapter AdapterPersister) *persisterPort {
	return &persisterPort{
		adapter:  adapter,
		location: location,
	}
}

func (p persisterPort) Load() (domain.PersistentState, error) {
	state, err := p.adapter.Load(p.location)
	if err != nil {
		log.Warn("CANNOT LOAD PERSISTENT STATE, USING DEFAULT")
		return domain.DEFAULT_PERSISTENT_STATE, err
	}
	return *state, nil
}

func (p persisterPort) Save(state domain.PersistentState) error {
	return nil // for development
	// return p.adapter.Save(p.location, state.CurrentTerm, state.VotedFor, state.MachineLogs)
}
