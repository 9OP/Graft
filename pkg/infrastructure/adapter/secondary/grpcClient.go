package secondaryAdapter

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"
	"os"
	"time"

	"graft/pkg/infrastructure/adapter/clusterRpc"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type clusterClient struct{}

func NewClusterClient() *clusterClient {
	return &clusterClient{}
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

func withClient[
	K *clusterRpc.AppendEntriesOutput |
		*clusterRpc.RequestVoteOutput |
		*clusterRpc.ConfigurationOutput |
		*clusterRpc.ExecuteOutput |
		*clusterRpc.Nil](target string, fn func(c clusterRpc.ClusterClient) (K, error),
) (K, error) {
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
	c := clusterRpc.NewClusterClient(conn)
	return fn(c)
}

func (r *clusterClient) AppendEntries(target string, input *clusterRpc.AppendEntriesInput) (*clusterRpc.AppendEntriesOutput, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 350*time.Millisecond)
	defer cancel()

	return withClient(
		target,
		func(c clusterRpc.ClusterClient) (*clusterRpc.AppendEntriesOutput, error) {
			return c.AppendEntries(ctx, input)
		})
}

func (r *clusterClient) RequestVote(target string, input *clusterRpc.RequestVoteInput) (*clusterRpc.RequestVoteOutput, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 350*time.Millisecond)
	defer cancel()

	return withClient(
		target,
		func(c clusterRpc.ClusterClient) (*clusterRpc.RequestVoteOutput, error) {
			return c.RequestVote(ctx, input)
		})
}

func (r *clusterClient) PreVote(target string, input *clusterRpc.RequestVoteInput) (*clusterRpc.RequestVoteOutput, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 350*time.Millisecond)
	defer cancel()

	return withClient(
		target,
		func(c clusterRpc.ClusterClient) (*clusterRpc.RequestVoteOutput, error) {
			return c.PreVote(ctx, input)
		})
}

func (r *clusterClient) Execute(target string, input *clusterRpc.ExecuteInput) (*clusterRpc.ExecuteOutput, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 350*time.Millisecond)
	defer cancel()

	return withClient(
		target,
		func(c clusterRpc.ClusterClient) (*clusterRpc.ExecuteOutput, error) {
			return c.Execute(ctx, input)
		})
}

func (r *clusterClient) LeadershipTransfer(target string, input *clusterRpc.Nil) (*clusterRpc.Nil, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 350*time.Millisecond)
	defer cancel()

	return withClient(
		target,
		func(c clusterRpc.ClusterClient) (*clusterRpc.Nil, error) {
			return c.LeadershipTransfer(ctx, input)
		})
}

func (r *clusterClient) ClusterConfiguration(target string, input *clusterRpc.Nil) (*clusterRpc.ConfigurationOutput, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 350*time.Millisecond)
	defer cancel()

	return withClient(
		target,
		func(c clusterRpc.ClusterClient) (*clusterRpc.ConfigurationOutput, error) {
			return c.Configuration(ctx, input)
		})
}

func (r *clusterClient) Shutdown(target string, input *clusterRpc.Nil) (*clusterRpc.Nil, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 350*time.Millisecond)
	defer cancel()

	return withClient(
		target,
		func(c clusterRpc.ClusterClient) (*clusterRpc.Nil, error) {
			return c.Shutdown(ctx, input)
		})
}

func (r *clusterClient) Ping(target string, input *clusterRpc.Nil) (*clusterRpc.Nil, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 350*time.Millisecond)
	defer cancel()

	return withClient(
		target,
		func(c clusterRpc.ClusterClient) (*clusterRpc.Nil, error) {
			return c.Ping(ctx, input)
		})
}
