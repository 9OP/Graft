package persister

import (
	"graft/app/entity"
	"sync"
)

type Service struct {
	repository Repository
	location   string
	mu         sync.Mutex
}

func NewService(location string, repository Repository) *Service {
	return &Service{
		location:   location,
		repository: repository,
	}
}

func (s *Service) SaveState(state *entity.Persistent) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.repository.Save(s.location, state)
}

func (s *Service) LoadState() (*entity.Persistent, error) {
	state, err := s.repository.Load(s.location)

	var DEFAULT_STATE entity.Persistent = entity.Persistent{
		CurrentTerm: 0,
		VotedFor:    "",
		MachineLogs: []entity.MachineLog{{Term: 0, Value: ""}},
	}

	if err != nil {
		return &DEFAULT_STATE, err
	}

	return state, nil
}
