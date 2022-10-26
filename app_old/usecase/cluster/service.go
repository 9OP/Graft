package cluster

import (
	"errors"
	"graft/app/domain/entity"
	"sync"
)

type Service struct {
	repository Repository
}

func NewService(repo Repository) *Service {
	return &Service{repository: repo}
}

type FsmResponse = interface{}

func (s *Service) NewCmd(cmd string) (FsmResponse, error) {
	if !s.repository.IsLeader() {
		return nil, errors.New("NOT_LEADER")
	}

	// Append cmd to leader logs
	state := s.repository.GetState()
	quorum := s.repository.GetQuorum()
	entries := []string{cmd}
	accepted := 1
	s.repository.AppendLogs(entries)

	// Broadcast new log to peers
	var m sync.Mutex
	appendLogs := func(p entity.Peer) {
		input := s.repository.AppendEntriesInput(p)
		if res, err := s.repository.AppendEntries(p, input); err == nil {
			if res.Term > state.CurrentTerm {
				s.repository.DowngradeFollower(res.Term, p.Id)
				return
			}
			if res.Success {
				m.Lock()
				defer m.Unlock()
				accepted += 1
			}
		}
	}
	s.repository.Broadcast(appendLogs)

	if accepted >= quorum {
		// Could be downgraded to Follower during broadcast
		if s.repository.IsLeader() {
			s.repository.CommitCmd(cmd)
			return s.repository.ExecFsmCmd(cmd), nil
		}
		return nil, errors.New("NOT_LEADER")
	}

	return nil, errors.New("QUORUM_NOT_REACHED")
}
