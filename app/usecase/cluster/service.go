package cluster

import (
	"graft/app/domain/entity"
	"time"
)

type service struct {
	repository repository
}

func NewService(repository repository) *service {
	return &service{repository}
}

func (s *service) ExecuteCommand(command string) (interface{}, error) {
	if !s.repository.IsLeader() {
		leader := s.repository.GetLeader()
		return nil, entity.NewNotLeaderError(leader)
	}

	applied := s.repository.Execute(command)

	select {
	case <-time.After(2 * time.Second):
		return nil, entity.NewTimeoutError()
	case result := <-applied:
		return result, nil
	}
}

func (s *service) ExecuteQuery(query string) (interface{}, error) {
	return nil, nil
}
