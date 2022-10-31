package cluster

import (
	"errors"
	"graft/app/domain/entity"
	srv "graft/app/domain/service"
	"time"
)

type service struct {
	applied chan interface{}
	server  *srv.Server
}

func NewService(s *srv.Server, a chan interface{}) *service {
	return &service{
		applied: a,
		server:  s,
	}
}

func (s *service) ExecuteCommand(command string) (interface{}, error) {
	if !s.server.IsLeader() {
		return nil, errors.New("NOT_LEADER")
	}

	state := s.server.GetState()

	timeout := time.NewTimer(2 * time.Second)

	newLog := entity.LogEntry{Value: command, Term: state.CurrentTerm}
	// "github.com/google/uuid"
	//  id := uuid.New().String() // use id to track applied command
	//  newLog := entity.LogEntry{Value: command, Term: state.CurrentTerm, Uuid: id}
	s.server.AppendLogs([]entity.LogEntry{newLog}, state.GetLastLogIndex())

	select {
	case <-timeout.C:
		return nil, errors.New("TIMEOUT")
	case result := <-s.applied:
		return result, nil
		// case resUuid, result := <-s.applied:
		// if resUuid == id { return result, nil}
	}
}
