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
		if s.clusterNode.HasLeader() {
			leader := s.clusterNode.Leader()
			return nil, entity.NewNotLeaderError(leader)
		}
		return nil, entity.NewUnknownLeaderError()
	}

	res := <-s.clusterNode.ExecuteCommand(command)
	return res.Out, res.Err
}

func (s *service) ExecuteQuery(query string, weakConsistency bool) ([]byte, error) {
	if !s.clusterNode.IsLeader() && !weakConsistency {
		if s.clusterNode.HasLeader() {
			leader := s.clusterNode.Leader()
			return nil, entity.NewNotLeaderError(leader)
		}
		return nil, entity.NewUnknownLeaderError()
	}

	// TODO: Should ensure that the node has the lastest state version
	// to avoid stale read...
	// If leader:
	//	Should get a quorum before returning value to client
	// If not leader:
	// 	Should get same commit index as knwown leader
	res := <-s.clusterNode.ExecuteQuery(query)
	return res.Out, res.Err
}
