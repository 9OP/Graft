package main

import (
	"context"
	"fmt"
	"graft/graft_rpc"
	"graft/models"
	"log"
	"net"
	"time"

	"google.golang.org/grpc"
)

const PERSISTENT_STATE_FILE = "state.json"
const HOST = "127.0.0.1:7679"
const ELECTION_TIMEOUT = 300 // ms

func UnaryInterceptor(state *models.ServerState) grpc.ServerOption {
	middleware := func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		// Attach server state to the rpc context
		ctx = context.WithValue(ctx, "graft_server_state", state)

		// start := time.Now()

		// Calls the handler
		h, err := handler(ctx, req)

		// log.Printf("Request - Method:%s\tDuration:%s\tError:%v\n",
		// 	info.FullMethod,
		// 	time.Since(start),
		// 	err)

		return h, err
	}

	return grpc.UnaryInterceptor(middleware)
}

func StartGrpcServer(state *models.ServerState) {
	log.Println("start grpc server")
	lis, err := net.Listen("tcp", HOST)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer(UnaryInterceptor(state))
	graft_rpc.RegisterRpcServer(grpcServer, &graft_rpc.Service{})
	err = grpcServer.Serve(lis)
	if err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

func StartGraftServer(state *models.ServerState) {
	fmt.Println("start graft server")
	srv := models.Server{Timeout: time.NewTicker(3000 * time.Millisecond)}
	srv.Start(state)
}

func main() {
	state := models.NewServerState()

	go StartGraftServer(state) // Timeout / Heartbeat
	StartGrpcServer(state)     // RPC
}
