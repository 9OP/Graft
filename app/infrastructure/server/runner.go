package server

type runner interface {
	Run()
}

type runnerServer struct {
	runners []runner
}

func NewRunner(runners ...runner) *runnerServer {
	return &runnerServer{runners}
}

func (s *runnerServer) Start() {
	for {
		for _, runner := range s.runners {
			go runner.Run()
		}
	}
}
