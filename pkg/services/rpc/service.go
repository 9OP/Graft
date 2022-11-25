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

func (s service) AppendEntries(input *domain.AppendEntriesInput) (*domain.AppendEntriesOutput, error) {
	node := s.node

	output := &domain.AppendEntriesOutput{
		Term:    node.CurrentTerm(),
		Success: false,
	}
	if input.Term < node.CurrentTerm() || !s.node.IsActivePeer(input.LeaderId) {
		return output, nil
	}
	if input.Term > node.CurrentTerm() {
		node.DowngradeFollower(input.Term)
	}

	node.SetLeader(input.LeaderId)
	node.Heartbeat()

	localPrevLog, err := node.Log(input.PrevLogIndex)
	if localPrevLog.Term == input.PrevLogTerm && err == nil {
		node.AppendLogs(input.PrevLogIndex, input.Entries...)
		output.Success = true
	} else {
		node.DeleteLogsFrom(input.PrevLogIndex)
	}

	node.UpdateLeaderCommitIndex(input.LeaderCommit)

	return output, nil
}

func (s service) RequestVote(input *domain.RequestVoteInput) (*domain.RequestVoteOutput, error) {
	node := s.node
	if s.node.IsShuttingDown() {
		return nil, domain.ErrShuttingDown
	}

	output := &domain.RequestVoteOutput{
		Term:        node.CurrentTerm(),
		VoteGranted: false,
	}
	if input.Term < node.CurrentTerm() || !s.node.IsActivePeer(input.CandidateId) {
		return output, nil
	}
	if input.Term > node.CurrentTerm() {
		node.DowngradeFollower(input.Term)
	}

	isUpToDate := node.IsLogUpToDate(input.LastLogIndex, input.LastLogTerm)
	canGrantVote := node.CanGrantVote(input.CandidateId)

	if canGrantVote && isUpToDate {
		node.Heartbeat()
		node.GrantVote(input.CandidateId)
		output.VoteGranted = true
	}

	return output, nil
}

func (s service) PreVote(input *domain.RequestVoteInput) (*domain.RequestVoteOutput, error) {
	node := s.node
	if s.node.IsShuttingDown() {
		return nil, domain.ErrShuttingDown
	}

	output := &domain.RequestVoteOutput{
		Term:        node.CurrentTerm(),
		VoteGranted: false,
	}

	if !s.node.IsActivePeer(input.CandidateId) {
		return output, nil
	}

	hasLeader := node.HasLeader()
	isUpToDate := node.IsLogUpToDate(input.LastLogIndex, input.LastLogTerm)

	if !hasLeader && isUpToDate {
		output.VoteGranted = true
	}

	return output, nil
}

func (s service) Execute(input *domain.ExecuteInput) (*domain.ExecuteOutput, error) {
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

func (s service) ClusterConfiguration() (*domain.ClusterConfiguration, error) {
	configuration := s.node.GetClusterConfiguration()
	return &configuration, nil
}

func (s service) Shutdown() {
	fmt.Println("shutting down in 3s")

	s.node.Shutdown()

	time.AfterFunc(3*time.Second, func() {
		fmt.Println("shutdown")
		close(s.quit)
	})
}
