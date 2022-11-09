package service

type signals struct {
	ShiftRole          chan struct{}
	SaveState          chan struct{}
	ResetElectionTimer chan struct{}
	ResetLeaderTicker  chan struct{}
	Commit             chan struct{}
}

func newSignals() signals {
	sng := signals{
		SaveState:          make(chan struct{}, 1),
		ShiftRole:          make(chan struct{}, 1),
		ResetElectionTimer: make(chan struct{}, 1),
		ResetLeaderTicker:  make(chan struct{}, 1),
		Commit:             make(chan struct{}, 1),
	}
	sng.resetTimeout()
	return sng
}

func (s *signals) saveState() {
	s.SaveState <- struct{}{}
}

func (s *signals) shiftRole() {
	s.ShiftRole <- struct{}{}
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
