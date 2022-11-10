package cluster

import (
	"graft/pkg/domain/entity"
	domain "graft/pkg/domain/service"
)

type service struct {
	clusterNode *domain.ClusterNode
}

func NewService(clusterNode *domain.ClusterNode) *service {
	return &service{clusterNode}
}

func (s *service) ExecuteCommand(command string) ([]byte, error) {
	if !s.clusterNode.IsLeader() {
		leader := s.clusterNode.Leader()
		return nil, entity.NewNotLeaderError(leader)
	}

	res := <-s.clusterNode.ExecuteCommand(command)
	return res.Out, res.Err
}

func (s *service) ExecuteQuery(query string, weakConsistency bool) ([]byte, error) {
	if !s.clusterNode.IsLeader() && !weakConsistency {
		leader := s.clusterNode.Leader()
		return nil, entity.NewNotLeaderError(leader)
	}

	res := <-s.clusterNode.ExecuteQuery(query)
	return res.Out, res.Err
}
