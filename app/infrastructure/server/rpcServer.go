package server

import (
	"fmt"
	"graft/app/infrastructure/adapter"
	"graft/app/infrastructure/adapter/rpc"
	"log"
	"net"

	"google.golang.org/grpc"
)

type rpcServer struct {
	rpcApi adapter.RpcApi
}

func NewRpcServer(rpcApi adapter.RpcApi) *rpcServer {
	return &rpcServer{rpcApi}
}

func (r *rpcServer) Start(port string) {
	log.Println("START RECEIVER SERVER")

	addr := fmt.Sprintf("%s:%s", "127.0.0.1", port)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("Failed to listen: \n\t%v\n", err)
	}

	server := grpc.NewServer()
	rpc.RegisterRpcServer(server, &r.rpcApi)

	if err := server.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: \n\t%v\n", err)
	}
}
