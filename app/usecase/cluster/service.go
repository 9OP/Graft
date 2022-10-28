package cluster

import (
	"errors"
	srv "graft/app/domain/service"
)

type service struct {
	server     *srv.Server
	repository repository
}

func NewService(s *srv.Server, r repository) *service {
	return &service{
		server:     s,
		repository: r,
	}
}

func (s *service) NewCmd(cmd string) (interface{}, error) {
	if !s.server.IsLeader() {
		return nil, errors.New("NOT_LEADER")
	}

	/*
		- Append local logs with cmd
		- Start timer 3s
		- Select on signals(timer, applied):
			- If timer timesout => return nil, timeoutError
			- If applied => return applied result

	*/

	s.server.AppendLogs([]string{cmd})

	return nil, nil

	// // state := s.server.GetState()
	// // quorum := s.server.GetQuorum()
	// // accepted := 1

	// /*
	// 	Should retry until timeout 4s for instance
	// 	- when res is not success then add a query to the queue
	// 	- increment decrement logs in peers
	// 	- update leader nextIndex / matchIndex

	// */

	// // Broadcast new log to peers
	// var m sync.Mutex
	// appendLogs := func(p entity.Peer) {
	// 	input := s.server.GetAppendEntriesInput(entries)
	// 	if res, err := s.repository.AppendEntries(p, &input); err == nil {
	// 		if res.Term > state.CurrentTerm {
	// 			s.server.DowngradeFollower(res.Term, p.Id)
	// 			return
	// 		}
	// 		if res.Success {
	// 			m.Lock()
	// 			defer m.Unlock()
	// 			accepted += 1
	// 		}
	// 	}
	// }
	// s.server.Broadcast(appendLogs)

	// // Log is replicated to a majority of
	// // server, we can execute the command
	// if accepted >= quorum {
	// 	// Could be downgraded to Follower during broadcast
	// 	// We do not commit explicitly if server is not a leader anymore
	// 	// There will be an election and futur command will commit implicitly
	// 	// the current command.
	// 	if s.server.IsLeader() {
	// 		s.server.CommitCmd(cmd)
	// 	}
	// 	return s.server.ExecFsmCmd(cmd)
	// }

	// return nil, errors.New("QUORUM_NOT_REACHED")
}
