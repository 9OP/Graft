package rpc

import (
	"fmt"
	"time"

	"graft/pkg/domain"
)

type service struct {
	node *domain.Node
	quit chan struct{}
}

func NewService(node *domain.Node, quit chan struct{}) *service {
	return &service{node, quit}
}

func (s *service) AppendEntries(input *domain.AppendEntriesInput) (*domain.AppendEntriesOutput, error) {
	if !s.node.IsActivePeer(input.LeaderId) {
		return nil, domain.ErrNotActive
	}

	output := &domain.AppendEntriesOutput{
		Term:    s.node.CurrentTerm(),
		Success: false,
	}
	if input.Term < s.node.CurrentTerm() {
		return output, nil
	}
	if input.Term > s.node.CurrentTerm() {
		s.node.DowngradeFollower(input.Term)
	}

	s.node.SetLeader(input.LeaderId)
	s.node.Heartbeat()

	localPrevLog, err := s.node.Log(input.PrevLogIndex)
	if localPrevLog.Term == input.PrevLogTerm && err == nil {
		s.node.AppendLogs(input.PrevLogIndex, input.Entries...)
		output.Success = true
	} else {
		s.node.DeleteLogsFrom(input.PrevLogIndex)
	}

	if !s.node.IsShuttingDown() {
		// Prevent applying new logs when shutting down
		s.node.UpdateLeaderCommitIndex(input.LeaderCommit)
	}

	return output, nil
}

func (s *service) RequestVote(input *domain.RequestVoteInput) (*domain.RequestVoteOutput, error) {
	if s.node.IsShuttingDown() {
		return nil, domain.ErrShuttingDown
	}
	if !s.node.IsActivePeer(input.CandidateId) {
		return nil, domain.ErrNotActive
	}

	output := &domain.RequestVoteOutput{
		Term:        s.node.CurrentTerm(),
		VoteGranted: false,
	}
	if input.Term < s.node.CurrentTerm() {
		return output, nil
	}
	if input.Term > s.node.CurrentTerm() {
		s.node.DowngradeFollower(input.Term)
	}

	isUpToDate := s.node.IsLogUpToDate(input.LastLogIndex, input.LastLogTerm)
	canGrantVote := s.node.CanGrantVote(input.CandidateId)

	if canGrantVote && isUpToDate {
		s.node.Heartbeat()
		s.node.GrantVote(input.CandidateId)
		output.VoteGranted = true
	}

	return output, nil
}

func (s *service) PreVote(input *domain.RequestVoteInput) (*domain.RequestVoteOutput, error) {
	if s.node.IsShuttingDown() {
		return nil, domain.ErrShuttingDown
	}

	output := &domain.RequestVoteOutput{
		Term:        s.node.CurrentTerm(),
		VoteGranted: false,
	}

	if !s.node.IsActivePeer(input.CandidateId) {
		return output, nil
	}

	hasLeader := s.node.HasLeader()
	isUpToDate := s.node.IsLogUpToDate(input.LastLogIndex, input.LastLogTerm)

	if !hasLeader && isUpToDate {
		output.VoteGranted = true
	}

	return output, nil
}

func (s *service) Execute(input *domain.ExecuteInput) (*domain.ExecuteOutput, error) {
	if s.node.IsShuttingDown() {
		return nil, domain.ErrShuttingDown
	}

	// TODO: implements consistency
	if !s.node.IsLeader() {
		return nil, domain.ErrNotLeader
	}

	res := <-s.node.ExecuteCommand(*input)

	// Should we separate res.Err and error ? arent they the same ?

	return &res, nil
}

func (s *service) ClusterConfiguration() (*domain.ClusterConfiguration, error) {
	configuration := s.node.GetClusterConfiguration()
	return &configuration, nil
}

func (s *service) Shutdown() {
	fmt.Println("shutting down in 3s")

	s.node.Shutdown()

	time.AfterFunc(3*time.Second, func() {
		fmt.Println("shutdown")
		close(s.quit)
	})
}
