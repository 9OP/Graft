package secondaryPort

import (
	"graft/pkg/domain"
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

func (p persisterPort) Load() (domain.PersistentState, error) {
	ps, err := p.adapter.Load(p.location)
	if err != nil {
		log.Warn("CANNOT LOAD PERSISTENT STATE, USING DEFAULT")
		return domain.DEFAULT_PERSISTENT_STATE, err
	}
	res := domain.PersistentState{
		CurrentTerm: ps.CurrentTerm,
		VotedFor:    ps.VotedFor,
		MachineLogs: ps.Logs,
	}
	return res, nil
}

func (p persisterPort) Save(currentTerm uint32, votedFor string, machineLogs []domain.LogEntry) error {
	return p.adapter.Save(p.location, currentTerm, votedFor, machineLogs)
}
