package graft_rpc

import (
	"context"
	"graft/models"
	"log"
)

type Service struct {
	UnimplementedRpcServer
}

func (s *Service) AppendEntries(ctx context.Context, entries *AppendEntriesInput) (*AppendEntriesOutput, error) {
	log.Println("->AppendEntries")
	state := ctx.Value("graft_server_state").(*models.ServerState)
	state.Heartbeat <- true

	// Convert receiver back to follower
	if entries.Term > int32(state.CurrentTerm) {
		log.Println("->AppendEntries, convert to follower")
		state.SwitchRole(models.Follower)
		// TODO: should set in lock
		state.VotedFor = ""
		state.CurrentTerm = uint16(entries.Term)
	}

	return &AppendEntriesOutput{}, nil
}

func (s *Service) RequestVote(ctx context.Context, vote *RequestVoteInput) (*RequestVoteOutput, error) {
	log.Println("->RequestVote")
	state := ctx.Value("graft_server_state").(*models.ServerState)
	// state.Heartbeat <- true // Prevent receiver to turn into a candidate

	// Leader does not grant vote
	if state.IsRole(models.Leader) {
		return &RequestVoteOutput{Term: int32(state.CurrentTerm), VoteGranted: false}, nil
	}

	// Convert receiver back to follower
	if vote.Term > int32(state.CurrentTerm) {
		log.Println("->RequestVote, convert to follower")
		state.SwitchRole(models.Follower)
		state.VotedFor = ""
		state.CurrentTerm = uint16(vote.Term)
	}

	// Candidate's term is outdated
	if vote.Term < int32(state.CurrentTerm) {
		log.Println("->RequestVote: Candidate term is outdated")
		return &RequestVoteOutput{Term: int32(state.CurrentTerm), VoteGranted: false}, nil
	}

	// Receiver has not voted yet
	if state.VotedFor == "" || state.VotedFor == vote.CandidateId {
		currentLogIndex := len(state.PersistentState.Log)
		currentLogTerm := state.PersistentState.Log[currentLogIndex-1].Term

		if currentLogTerm <= uint16(vote.LastLogTerm) && currentLogIndex <= int(vote.LastLogIndex) {
			// Candidate's log is at least up-to-date as receiver
			state.VotedFor = vote.CandidateId
			return &RequestVoteOutput{Term: int32(state.CurrentTerm), VoteGranted: true}, nil
		}

		// Candidate's log is not up-to-date
		return &RequestVoteOutput{Term: int32(state.CurrentTerm), VoteGranted: false}, nil
	}

	// Receiver already voted
	log.Println("->RequestVote: recevier already voted")
	return &RequestVoteOutput{Term: int32(state.CurrentTerm), VoteGranted: false}, nil
}
