package serversnew

import (
	"graft/app/entity"
	"graft/app/usecase/persister"
	"graft/app/usecase/runner"
	"log"
)

type Runner struct {
	server    *entity.Server
	persister *persister.Service
}

func NewRunner(srv *entity.Server, ps *persister.Service) *Runner {
	return &Runner{
		srv,
		ps,
	}
}

func (r *Runner) Start(service *runner.Service) {
	log.Println("START RUNNER NEW SERVER")
	srv := r.server

	for {
		select {

		case <-srv.ShiftRole:
			switch {
			case srv.IsFollower():
				//
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
			continue

		case <-srv.ResetLeaderTicker:
			continue
		}

	}
}
