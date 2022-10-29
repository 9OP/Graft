package receiver

import (
	"fmt"
	"graft/app/domain/entity"
)

type service struct {
	repository
}

func NewService(repository repository) *service {
	return &service{repository}
}

func (s *service) AppendEntries(input *entity.AppendEntriesInput) (*entity.AppendEntriesOutput, error) {
	state := s.GetState()

	output := &entity.AppendEntriesOutput{
		Term:    state.CurrentTerm,
		Success: false,
	}
	if input.Term < state.CurrentTerm {
		return output, nil
	}
	if input.Term > state.CurrentTerm {

		s.DowngradeFollower(input.Term, input.LeaderId)
	}

	s.SetClusterLeader(input.LeaderId)
	s.Heartbeat()

	localPrevLog := state.GetLogByIndex(input.PrevLogIndex)
	if localPrevLog.Term == input.PrevLogTerm {
		if len(input.Entries) > 0 {
			fmt.Println("append logs", input.Entries, len(input.Entries))
			s.AppendLogs(input.Entries)
		}
		output.Success = true
	} else {
		fmt.Println("DLETE LOGS", input.PrevLogIndex)
		s.DeleteLogsFrom(input.PrevLogIndex)
	}

	if input.LeaderCommit > state.CommitIndex {
		s.SetCommitIndex(min(state.GetLastLogIndex(), input.LeaderCommit))
	}

	return output, nil
}

func (s *service) RequestVote(input *entity.RequestVoteInput) (*entity.RequestVoteOutput, error) {
	state := s.GetState()

	output := &entity.RequestVoteOutput{
		Term:        state.CurrentTerm,
		VoteGranted: false,
	}
	if input.Term < state.CurrentTerm {
		return output, nil
	}
	if input.Term > state.CurrentTerm {
		// Do not log follow cluster leader in req vote
		// FollowCluster leader only in appendEntries
		s.DowngradeFollower(input.Term, input.CandidateId)
	}

	if s.GrantVote(input.CandidateId, input.LastLogIndex, input.LastLogTerm) {
		output.VoteGranted = true
		s.Heartbeat()
	}

	return output, nil
}

// Move to utils ?
func min[K uint | uint8 | uint16 | uint32 | uint64 | int](value_0, value_1 K) K {
	if value_0 < value_1 {
		return value_0
	}
	return value_1
}
