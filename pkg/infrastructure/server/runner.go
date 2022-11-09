package server

import "sync"

type runner interface {
	Dispatch()
}

type runnerServer struct {
	runners []runner
}

func NewRunner(runners ...runner) *runnerServer {
	return &runnerServer{runners}
}

func (s *runnerServer) Start() {
	var wg sync.WaitGroup
	for _, runner := range s.runners {
		wg.Add(1)
		go (func() {
			for {
				runner.Dispatch()
			}
		})()
	}
	wg.Wait()
}
