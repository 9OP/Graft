package cluster

import (
	"graft/pkg/domain/entity"
)

type service struct {
	repository repository
}

func NewService(repository repository) *service {
	return &service{repository}
}

func (s *service) ExecuteCommand(command string) ([]byte, error) {
	if !s.repository.IsLeader() {
		leader := s.repository.Leader()
		return nil, entity.NewNotLeaderError(leader)
	}

	res := <-s.repository.ExecuteCommand(command)
	return res.Out, res.Err
}

func (s *service) ExecuteQuery(query string, weakConsistency bool) ([]byte, error) {
	if !s.repository.IsLeader() && !weakConsistency {
		leader := s.repository.Leader()
		return nil, entity.NewNotLeaderError(leader)
	}

	/* Consistency:
	- default
	- strong
	- weak
	*/

	res := <-s.repository.ExecuteQuery(query)
	return res.Out, res.Err
}
