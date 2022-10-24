package serversnew

import (
	"graft/app/entity"
	entitynew "graft/app/entity_new"
	"graft/app/usecase/persister"
	"graft/app/usecase/runner"
	"log"
)

type Runner struct {
	server    *entity.Server
	timeout   *entitynew.Timeout
	persister *persister.Service
}

func NewRunner(srv *entity.Server, tm *entitynew.Timeout, ps *persister.Service) *Runner {
	return &Runner{srv, tm, ps}
}

func (r *Runner) Start(service *runner.Service) {
	log.Println("START RUNNER NEW SERVER")
	srv := r.server
	log.Println("service")

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
			r.persister.SaveState(&state.Persistent)

		case <-srv.ResetElectionTimer:
			r.timeout.ResetElectionTimer()

		case <-srv.ResetLeaderTicker:
			r.timeout.ResetLeaderTicker()
		}

	}
}
