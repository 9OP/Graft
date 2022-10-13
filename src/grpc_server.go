package main

import (
	"context"
	rpc "graft/src/api/graft_rpc"
	"log"
	"net"

	"google.golang.org/grpc"
)

// Move to global const / or config file
// type graftServerStateCtxKeyType string

// const GRAFT_SERVER_STATE graftServerStateCtxKeyType = "graft_server_state"

func SendRequestVoteRpc(host string, input *rpc.RequestVoteInput) (*rpc.RequestVoteOutput, error) {
	var conn *grpc.ClientConn
	conn, _ = grpc.Dial("127.0.0.1:"+host, grpc.WithInsecure())
	defer conn.Close()

	c := rpc.NewRpcClient(conn)

	res, err := c.RequestVote(context.Background(), input)
	return res, err
}

func SendAppendEntriesRpc(host string, input *rpc.AppendEntriesInput) (*rpc.AppendEntriesOutput, error) {
	var conn *grpc.ClientConn
	conn, _ = grpc.Dial("127.0.0.1:"+host, grpc.WithInsecure())
	defer conn.Close()

	c := rpc.NewRpcClient(conn)

	res, err := c.AppendEntries(context.Background(), input)
	return res, err
}

func attachContextMiddleware(state *ServerState) grpc.ServerOption {
	middleware := func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		ctx = context.WithValue(ctx, "graft_server_state", state)
		h, err := handler(ctx, req)
		return h, err
	}

	return grpc.UnaryInterceptor(middleware)
}

type GrpcServer struct {
	serverState *ServerState
}

func NewGrpcServer(state *ServerState) *GrpcServer {
	return &GrpcServer{
		serverState: state,
	}
}

func (srv *GrpcServer) Start(port string) {
	log.Println("START GRPC SERVER")

	lis, err := net.Listen("tcp", "127.0.0.1:"+port)
	if err != nil {
		log.Fatalf("Failed to listen: \n\t%v\n", err)
	}

	server := grpc.NewServer(attachContextMiddleware(srv.serverState))
	rpc.RegisterRpcServer(server, &Service{})

	if err := server.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: \n\t%v\n", err)
	}
}
