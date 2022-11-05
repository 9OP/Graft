package cluster

import (
	"graft/pkg/domain/entity"
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
		leader := s.repository.Leader()
		return nil, entity.NewNotLeaderError(leader)
	}

	applied := s.repository.ExecuteCommand(command)

	select {
	case <-time.After(2 * time.Second):
		return nil, entity.NewTimeoutError()
	case result := <-applied:
		return result, nil
	}
}

func (s *service) ExecuteQuery(query string, weakConsistency bool) (interface{}, error) {
	if !s.repository.IsLeader() && !weakConsistency {
		leader := s.repository.Leader()
		return nil, entity.NewNotLeaderError(leader)
	}

	applied := s.repository.ExecuteQuery(query)

	select {
	case <-time.After(2 * time.Second):
		return nil, entity.NewTimeoutError()
	case result := <-applied:
		return result, nil
	}
}