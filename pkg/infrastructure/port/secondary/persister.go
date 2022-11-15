package secondaryPort

import (
	"graft/pkg/domain"
	"graft/pkg/domain/state"
	adapter "graft/pkg/infrastructure/adapter/secondary"

	log "github.com/sirupsen/logrus"
)

type persisterPort struct {
	adapter  adapter.UseCaseJsonPersisterAdapter
	location string
}

func NewPersisterPort(location string, adapter adapter.UseCaseJsonPersisterAdapter) *persisterPort {
	return &persisterPort{
		adapter:  adapter,
		location: location,
	}
}

func (p persisterPort) Load() (*state.PersistentState, error) {
	ps, err := p.adapter.Load(p.location)
	if err != nil {
		log.Warn("CANNOT LOAD PERSISTENT STATE, USING DEFAULT")
		res := state.NewDefaultPersistentState()
		return &res, err
	}
	res := state.NewPersistentState(
		ps.CurrentTerm,
		ps.VotedFor,
		ps.Logs,
	)
	return &res, nil
}

func (p persisterPort) Save(currentTerm uint32, votedFor string, machineLogs []domain.LogEntry) error {
	return p.adapter.Save(p.location, currentTerm, votedFor, machineLogs)
}
