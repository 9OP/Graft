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
	log.Println("->req append entries")
	state := ctx.Value("graft_server_state").(*models.ServerState)
	state.Heartbeat <- true

	// Convert receiver back to follower
	if entries.Term > int32(state.CurrentTerm) {
		log.Println("convert back to follower, heartbeat with higher term")
		state.Role = models.Follower
	}

	// state.CommitIndex += 1

	// log.Println("sleep 8s")
	// time.Sleep(8 * time.Second)
	// log.Printf("RPC -> AppendEntries %v\n", entries)
	return &AppendEntriesOutput{}, nil
}

func (s *Service) RequestVote(ctx context.Context, vote *RequestVoteInput) (*RequestVoteOutput, error) {
	log.Println("->req requestVote")
	state := ctx.Value("graft_server_state").(*models.ServerState)
	state.Heartbeat <- true // Prevent receiver to turn into a candidate

	// Convert receiver back to follower
	if vote.Term > int32(state.CurrentTerm) {
		log.Println("convert back to follower, vote with higher term")
		state.Role = models.Follower
	}

	// Candidate's term is outdated
	if vote.Term < int32(state.CurrentTerm) {
		return &RequestVoteOutput{Term: int32(state.CurrentTerm), VoteGranted: false}, nil
	}

	// Receiver has not voted yet
	if state.VotedFor == "" || state.VotedFor == vote.CandidateId {
		currentLogIndex := len(state.PersistentState.Log)
		currentLogTerm := state.PersistentState.Log[currentLogIndex-1].Term

		if currentLogTerm <= uint16(vote.LastLogTerm) && currentLogIndex <= int(vote.LastLogIndex) {
			// Candidate's log is at least up-to-date as receiver
			return &RequestVoteOutput{Term: int32(state.CurrentTerm), VoteGranted: true}, nil
		}

		// Candidate's log is not up-to-date
		return &RequestVoteOutput{Term: int32(state.CurrentTerm), VoteGranted: false}, nil
	}

	// Receiver already voted
	return &RequestVoteOutput{Term: int32(state.CurrentTerm), VoteGranted: false}, nil
}
