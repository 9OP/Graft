package server

type runner interface {
	Run()
}

type runnerServer struct {
	runner
}

func NewRunner(runner runner) *runnerServer {
	return &runnerServer{runner}
}

func (s *runnerServer) Start() {
	s.runner.Run()
}
