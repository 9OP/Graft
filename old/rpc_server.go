package main

import (
	"context"
	"graft/src/rpc"
	"log"
	"net"

	"google.golang.org/grpc"
)

type graftServerStateCtxKeyType string

const GRAFT_SERVER_STATE graftServerStateCtxKeyType = "graft_server_state"

type Service struct {
	rpc.UnimplementedRpcServer
}

func (s *Service) AppendEntries(ctx context.Context, entries *rpc.AppendEntriesInput) (*rpc.AppendEntriesOutput, error) {
	state := ctx.Value(GRAFT_SERVER_STATE).(*ServerState)
	state.Ping()

	output := &rpc.AppendEntriesOutput{}

	if entries.Term > int32(state.CurrentTerm) {
		state.DowngradeToFollower(uint16(entries.Term))
	}

	return output, nil
}

func (s *Service) RequestVote(ctx context.Context, vote *rpc.RequestVoteInput) (*rpc.RequestVoteOutput, error) {
	state := ctx.Value(GRAFT_SERVER_STATE).(*ServerState)
	state.Ping() // retard receiver to raise to candidate

	output := &rpc.RequestVoteOutput{
		Term:        int32(state.CurrentTerm),
		VoteGranted: false,
	}

	if vote.Term > int32(state.CurrentTerm) {
		state.DowngradeToFollower(uint16(vote.Term))
	}

	if vote.Term < int32(state.CurrentTerm) {
		return output, nil
	}

	if state.CanVote(vote.CandidateId, int(vote.LastLogIndex), int(vote.LastLogTerm)) {
		state.Vote(vote.CandidateId)
		output.VoteGranted = true
	}

	return output, nil
}

func attachContextMiddleware(state *ServerState) grpc.ServerOption {
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

type RpcServer struct {
	serverState *ServerState
}

func NewRpcServer(state *ServerState) *RpcServer {
	return &RpcServer{
		serverState: state,
	}
}

func (srv *RpcServer) Start(port string) {
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
