package cluster

import (
	"errors"
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
		return nil, errors.New("NOT_LEADER")
	}

	timeout := time.NewTimer(2 * time.Second)
	applied := s.repository.Execute(command)

	select {
	case <-timeout.C:
		return nil, errors.New("TIMEOUT")
	case result := <-applied:
		return result, nil
	}
}
