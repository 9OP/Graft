package server

import (
	"crypto/tls"
	"fmt"
	"net"

	"graft/pkg/infrastructure/adapter/p2pRpc"

	log "github.com/sirupsen/logrus"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
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

func loadTLSCredentials() (credentials.TransportCredentials, error) {
	// Load server's certificate and private key
	serverCert, err := tls.LoadX509KeyPair("cert/server.pem", "cert/server.key")
	if err != nil {
		return nil, err
	}

	// Create the credentials and return it
	config := &tls.Config{
		Certificates: []tls.Certificate{serverCert},
		ClientAuth:   tls.NoClientCert,
	}

	return credentials.NewTLS(config), nil
}

func getListennerOrFail(port string) net.Listener {
	addr := fmt.Sprintf("%s:%s", "127.0.0.1", port)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("FAILED TO LISTEN\n\t%v\n", err)
	}
	return lis
}

func serveOrFail(server *grpc.Server, lis net.Listener) {
	if err := server.Serve(lis); err != nil {
		log.Fatalf("FAILED TO SERVE\n\t%v\n", err)
	}
}

func createGrpcServer(apis ...p2pRpc.RpcServer) *grpc.Server {
	creds, err := loadTLSCredentials()
	if err != nil {
		log.Fatalf("Failed to setup TLS: %v", err)
	}
	server := grpc.NewServer(grpc.Creds(creds))
	for _, api := range apis {
		p2pRpc.RegisterRpcServer(server, api)
	}
	return server
}
