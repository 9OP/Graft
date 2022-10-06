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
	state := ctx.Value("graft_server_state").(*models.ServerState)
	log.Printf("context value state: %v\n", state)
	state.Heartbeat <- true

	// state.CommitIndex += 1

	// log.Println("sleep 8s")
	// time.Sleep(8 * time.Second)
	// log.Printf("RPC -> AppendEntries %v\n", entries)
	return &AppendEntriesOutput{}, nil
}
