package service

type signals struct {
	Commit             chan struct{}
	ShiftRole          chan struct{}
	SaveState          chan struct{}
	SynchronizeLogs    chan struct{}
	ResetLeaderTicker  chan struct{}
	ResetElectionTimer chan struct{}
}

func newSignals() signals {
	sng := signals{
		Commit:             make(chan struct{}, 1),
		ShiftRole:          make(chan struct{}, 1),
		SaveState:          make(chan struct{}, 1),
		SynchronizeLogs:    make(chan struct{}, 1),
		ResetLeaderTicker:  make(chan struct{}, 1),
		ResetElectionTimer: make(chan struct{}, 1),
	}
	sng.resetTimeout()
	return sng
}

func (s *signals) commit() {
	s.Commit <- struct{}{}
}

func (s *signals) shiftRole() {
	s.ShiftRole <- struct{}{}
}

func (s *signals) saveState() {
	s.SaveState <- struct{}{}
}

func (s *signals) synchronizeLogs() {
	s.SynchronizeLogs <- struct{}{}
}

func (s *signals) resetLeaderTicker() {
	s.ResetLeaderTicker <- struct{}{}
}

func (s *signals) resetTimeout() {
	s.ResetElectionTimer <- struct{}{}
}
