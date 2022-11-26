package rpc

import (
	"encoding/json"
	"fmt"
	"sync/atomic"
	"time"

	"graft/pkg/domain"
)

type service struct {
	node   *domain.Node
	quit   chan struct{}
	client client
}

func NewService(node *domain.Node, client client, quit chan struct{}) *service {
	return &service{node, quit, client}
}

func (s service) AppendEntries(input *domain.AppendEntriesInput) (*domain.AppendEntriesOutput, error) {
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

func (s service) RequestVote(input *domain.RequestVoteInput) (*domain.RequestVoteOutput, error) {
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

func (s service) PreVote(input *domain.RequestVoteInput) (*domain.RequestVoteOutput, error) {
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

	isUpToDate := s.node.IsLogUpToDate(input.LastLogIndex, input.LastLogTerm)

	if isUpToDate {
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

	if err := s.validateExecuteInput(input); err != nil {
		return nil, err
	}

	res := <-s.node.ExecuteCommand(*input)

	// Should we separate res.Err and error ? arent they the same ?

	return &res, nil
}

func (s service) validateExecuteInput(input *domain.ExecuteInput) error {
	// Before activating a node
	// We need to confirm that the node is up
	// And that we can communicate with it.
	// Example: the new node is behing a NAT and we can connect to it.

	// This is necessary because adding multiple node behing NATs
	// Could break the cluster by updating the Quorum with unreachable
	// nodes.
	if input.Type == domain.LogConfiguration {
		var config domain.ConfigurationUpdate
		json.Unmarshal(input.Data, &config)

		if config.Type == domain.ConfActivatePeer {
			if err := s.client.Ping(config.Peer); err != nil {
				return domain.ErrUnreachable
			}
		}
	}

	return nil
}

func (s service) LeadershipTransfer() error {
	if s.node.IsShuttingDown() {
		return domain.ErrShuttingDown
	}

	if !s.preVote() {
		return domain.ErrPreVoteFailed
	}

	s.node.UpgradeCandidate()
	s.node.IncrementCandidateTerm()

	return nil
}

// Copied from core/service
func (s service) preVote() bool {
	input := s.node.RequestVoteInput()
	quorum := s.node.Quorum()
	var prevotesGranted uint32 = 1 // vote for self

	preVoteRoutine := func(p domain.Peer) {
		if res, err := s.client.PreVote(p, &input); err == nil {
			if res.VoteGranted {
				atomic.AddUint32(&prevotesGranted, 1)
			}
		}
	}
	s.node.Broadcast(preVoteRoutine, domain.BroadcastActive)

	quorumReached := int(prevotesGranted) >= quorum
	return quorumReached
}

func (s service) Configuration() (*domain.ClusterConfiguration, error) {
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

func (s service) Ping() error {
	if s.node.IsShuttingDown() {
		return domain.ErrShuttingDown
	}
	return nil
}
