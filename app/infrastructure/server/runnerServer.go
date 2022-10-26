package server

import (
	"graft/app/domain/entity"
	domainService "graft/app/domain/service"
	"graft/app/infrastructure/port"
	"graft/app_old/usecase/runner"
	"log"
)

type runnerServer struct {
	server    *domainService.Server
	timeout   *entity.Timeout
	persister port.Persister
}

func NewRunnerServer(server *domainService.Server, timeout *entity.Timeout, persister port.Persister) *runnerServer {
	return &runnerServer{server, timeout, persister}
}

func (s *runnerServer) Start(service *runner.Service) {
	log.Println("START RUNNER NEW SERVER")
	srv := s.server

	for {
		select {

		case <-srv.ShiftRole:
			switch {
			case srv.IsFollower():
				go service.RunFollower(srv)
			case srv.IsCandidate():
				go service.RunCandidate(srv)
			case srv.IsLeader():
				go service.RunLeader(srv)
			}

		case <-srv.SaveState:
			state := srv.GetState()
			s.persister.Save(&state.Persistent)

		case <-srv.ResetElectionTimer:
			s.timeout.ResetElectionTimer()

		case <-srv.ResetLeaderTicker:
			s.timeout.ResetLeaderTicker()
		}

	}
}
