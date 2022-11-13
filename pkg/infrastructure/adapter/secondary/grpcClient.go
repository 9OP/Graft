package secondaryAdapter

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"
	"os"
	"time"

	"graft/pkg/infrastructure/adapter/p2pRpc"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type UseCaseGrpcClient interface {
	AppendEntries(target string, input *p2pRpc.AppendEntriesInput) (*p2pRpc.AppendEntriesOutput, error)
	RequestVote(target string, input *p2pRpc.RequestVoteInput) (*p2pRpc.RequestVoteOutput, error)
	PreVote(target string, input *p2pRpc.RequestVoteInput) (*p2pRpc.RequestVoteOutput, error)
}

type grpcClient struct{}

func NewGrpcClient() *grpcClient {
	return &grpcClient{}
}

func loadTLSCredentials() (credentials.TransportCredentials, error) {
	// Load certificate of the CA who signed server's certificate
	pemServerCA, err := os.ReadFile("cert/server.pem")
	if err != nil {
		return nil, err
	}

	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(pemServerCA) {
		return nil, fmt.Errorf("failed to add server CA's certificate")
	}

	// Create the credentials and return it
	config := &tls.Config{
		RootCAs:            certPool,
		InsecureSkipVerify: true,
	}

	return credentials.NewTLS(config), nil
}

func withClient[K *p2pRpc.AppendEntriesOutput | *p2pRpc.RequestVoteOutput](target string, fn func(c p2pRpc.RpcClient) (K, error)) (K, error) {
	// Dial options
	creds, err := loadTLSCredentials()
	if err != nil {
		log.Fatalf("could not process the credentials: %v", err)
	}

	conn, err := grpc.Dial(target, grpc.WithTransportCredentials(creds))
	if err != nil {
		return nil, err
	}

	defer conn.Close()
	c := p2pRpc.NewRpcClient(conn)
	return fn(c)
}

func (r *grpcClient) AppendEntries(target string, input *p2pRpc.AppendEntriesInput) (*p2pRpc.AppendEntriesOutput, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 350*time.Millisecond)
	defer cancel()

	return withClient(
		target,
		func(c p2pRpc.RpcClient) (*p2pRpc.AppendEntriesOutput, error) {
			return c.AppendEntries(ctx, input)
		})
}

func (r *grpcClient) RequestVote(target string, input *p2pRpc.RequestVoteInput) (*p2pRpc.RequestVoteOutput, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 350*time.Millisecond)
	defer cancel()

	return withClient(
		target,
		func(c p2pRpc.RpcClient) (*p2pRpc.RequestVoteOutput, error) {
			return c.RequestVote(ctx, input)
		})
}

func (r *grpcClient) PreVote(target string, input *p2pRpc.RequestVoteInput) (*p2pRpc.RequestVoteOutput, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 350*time.Millisecond)
	defer cancel()

	return withClient(
		target,
		func(c p2pRpc.RpcClient) (*p2pRpc.RequestVoteOutput, error) {
			return c.PreVote(ctx, input)
		})
}
