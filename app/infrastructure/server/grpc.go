package server

import (
	"fmt"
	"graft/app/infrastructure/adapter/p2pRpc"
	"log"
	"net"

	"google.golang.org/grpc"
)

type rpcServer struct {
	apis []p2pRpc.RpcServer
}

func NewRpc(apis ...p2pRpc.RpcServer) *rpcServer {
	return &rpcServer{apis}
}

func (r *rpcServer) Start(port string) {
	grpcServer := createGrpcServer(r.apis...)
	lis := getListennerOrFail(port)
	serveOrFail(grpcServer, lis)
}

func getListennerOrFail(port string) net.Listener {
	addr := fmt.Sprintf("%s:%s", "127.0.0.1", port)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("Failed to listen: \n\t%v\n", err)
	}
	return lis
}

func serveOrFail(server *grpc.Server, lis net.Listener) {
	if err := server.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: \n\t%v\n", err)
	}
}

func createGrpcServer(apis ...p2pRpc.RpcServer) *grpc.Server {
	server := grpc.NewServer()
	for _, api := range apis {
		p2pRpc.RegisterRpcServer(server, api)
	}
	return server
}
