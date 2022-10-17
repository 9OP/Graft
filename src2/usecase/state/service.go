package state

import "graft/src2/entity"

type Service struct {
	repository Repository
	location   string
}

func NewService(location string, repository Repository) *Service {
	return &Service{
		location:   location,
		repository: repository,
	}
}

func (s *Service) SaveState(state *entity.PersistentState) error {
	return s.repository.Save(s.location, state)
}

func (s *Service) LoadState() (*entity.PersistentState, error) {
	state, err := s.repository.Load(s.location)

	var DEFAULT_STATE entity.PersistentState = entity.PersistentState{
		CurrentTerm: 0,
		VotedFor:    "",
		MachineLogs: []entity.MachineLog{{Term: 0, Value: ""}},
	}

	if err != nil {
		return &DEFAULT_STATE, err
	}

	return state, nil
}
