package service

import "graft/pkg/domain/entity"

type signals struct {
	ShiftRole          chan entity.Role
	SaveState          chan struct{}
	ResetElectionTimer chan struct{}
	ResetLeaderTicker  chan struct{}
	Commit             chan struct{}
}

func newSignals() signals {
	sng := signals{
		SaveState:          make(chan struct{}, 1),
		ShiftRole:          make(chan entity.Role, 1),
		ResetElectionTimer: make(chan struct{}, 1),
		ResetLeaderTicker:  make(chan struct{}, 1),
		Commit:             make(chan struct{}, 1),
	}
	sng.shiftRole(entity.Follower)
	sng.resetTimeout()
	return sng
}

func (s *signals) saveState() {
	s.SaveState <- struct{}{}
}

func (s *signals) shiftRole(role entity.Role) {
	s.ShiftRole <- role
}

func (s *signals) resetTimeout() {
	s.ResetElectionTimer <- struct{}{}
}

func (s *signals) resetLeaderTicker() {
	s.ResetLeaderTicker <- struct{}{}
}

func (s *signals) commit() {
	s.Commit <- struct{}{}
}
