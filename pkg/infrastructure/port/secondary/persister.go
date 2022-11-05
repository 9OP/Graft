package secondaryPort

import (
	"graft/pkg/domain/entity"
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

func (p persisterPort) Load() (*entity.PersistentState, error) {
	state, err := p.adapter.Load(p.location)
	if err != nil {
		log.Warn("CANNOT LOAD PERSISTENT STATE, USING DEFAULT")
		res := entity.NewDefaultPersistentState()
		return &res, err
	}
	res := entity.NewPersistentState(
		state.CurrentTerm,
		state.VotedFor,
		state.MachineLogs,
	)
	return &res, nil
}

func (p persisterPort) Save(currentTerm uint32, votedFor string, machineLogs []entity.LogEntry) error {
	return p.adapter.Save(p.location, currentTerm, votedFor, machineLogs)
}
