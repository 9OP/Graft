package server

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/netip"

	"graft/pkg/infrastructure/adapter/p2pRpc"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type rpcServer struct {
	api  p2pRpc.RpcServer
	port uint16
}

func NewRpc(api p2pRpc.RpcServer, port uint16) *rpcServer {
	return &rpcServer{api, port}
}

func (r *rpcServer) Start() error {
	addr, err := netip.ParseAddrPort(fmt.Sprintf("127.0.0.1:%v", r.port))
	if err != nil {
		return err
	}

	creds, err := loadTLSCredentials()
	if err != nil {
		return err
	}

	server := grpc.NewServer(grpc.Creds(creds))
	p2pRpc.RegisterRpcServer(server, r.api)

	lis, err := net.Listen("tcp", addr.String())
	if err != nil {
		return err
	}

	if err := server.Serve(lis); err != nil {
		return err
	}

	return nil
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
