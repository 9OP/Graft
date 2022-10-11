package api

import (
	"context"
	"graft/api/graft_rpc"
	"graft/models"
	"log"
	"net"

	"google.golang.org/grpc"
)

// Move to global const / or config file
type graftServerStateCtxKeyType string

const GRAFT_SERVER_STATE graftServerStateCtxKeyType = "graft_server_state"

func attachContextMiddleware(state *models.ServerState) grpc.ServerOption {
	middleware := func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		ctx = context.WithValue(ctx, GRAFT_SERVER_STATE, state)
		h, err := handler(ctx, req)
		return h, err
	}

	return grpc.UnaryInterceptor(middleware)
}

func StartGrpcServer(port string, state *models.ServerState) *grpc.Server {
	log.Println("START GRPC SERVER")
	lis, err := net.Listen("tcp", "127.0.0.1:"+port)
	if err != nil {
		log.Fatalf("Failed to listen: \n\t%v\n", err)
	}

	server := grpc.NewServer(attachContextMiddleware(state))
	graft_rpc.RegisterRpcServer(server, &graft_rpc.Service{})

	if err := server.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: \n\t%v\n", err)
	}

	return server
}
