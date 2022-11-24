package api

import (
	"graft/pkg/domain"
)

type service struct {
	clusterNode *domain.Node
}

func NewService(clusterNode *domain.Node) *service {
	return &service{clusterNode}
}

func (s *service) ExecuteCommand(command domain.ExecuteInput) ([]byte, error) {
	if !s.clusterNode.IsLeader() {
		if s.clusterNode.HasLeader() {
			leader := s.clusterNode.Leader()
			return nil, domain.NewNotLeaderError(leader)
		}
		return nil, domain.NewUnknownLeaderError()
	}

	res := <-s.clusterNode.ExecuteCommand(command)
	return res.Out, res.Err
}

func (s *service) ExecuteQuery(query string, weakConsistency bool) ([]byte, error) {
	if !s.clusterNode.IsLeader() && !weakConsistency {
		if s.clusterNode.HasLeader() {
			leader := s.clusterNode.Leader()
			return nil, domain.NewNotLeaderError(leader)
		}
		return nil, domain.NewUnknownLeaderError()
	}

	// TODO: Should ensure that the node has the lastest state version
	// to avoid stale read...
	// If leader:
	//	Should get a quorum before returning value to client
	// If not leader:
	// 	Should get same commit index as knwown leader
	// res := <-s.clusterNode.ExecuteQuery(query)
	// return res.Out, res.Err
	return nil, nil
}
