package infrastructure

import (
	"context"
	"fmt"
	"graft/app/infrastructure/adapter/rpc"
	"log"
	"net"

	"google.golang.org/grpc"
)

// Port ?
type serverApi struct {
	rpc.UnimplementedRpcServer
}

func (s *serverApi) AppendEntries(ctx context.Context, input *rpc.AppendEntriesInput) (*rpc.AppendEntriesOutput, error) {
	return nil, nil
}

func (s *serverApi) RequesVote(ctx context.Context, input *rpc.RequestVoteInput) (*rpc.RequestVoteOutput, error) {
	return nil, nil
}

// Adapter
func Start(port string, api *serverApi) {
	log.Println("START RECEIVER SERVER")

	addr := fmt.Sprintf("%s:%s", "127.0.0.1", port)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("Failed to listen: \n\t%v\n", err)
	}

	server := grpc.NewServer()
	rpc.RegisterRpcServer(server, api)

	if err := server.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: \n\t%v\n", err)
	}
}
