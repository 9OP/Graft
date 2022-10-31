package main

import (
	"encoding/json"
	"fmt"
	"graft/app/domain/entity"
	"graft/app/domain/service"
	primaryAdapter "graft/app/infrastructure/adapter/primary"
	secondaryAdapter "graft/app/infrastructure/adapter/secondary"
	primaryPort "graft/app/infrastructure/port/primary"
	secondaryPort "graft/app/infrastructure/port/secondary"
	"graft/app/infrastructure/server"
	"graft/app/usecase/cluster"
	"graft/app/usecase/receiver"
	"graft/app/usecase/runner"
	"os"

	log "github.com/sirupsen/logrus"
)

type Args struct {
	port    string
	apiPort string
	id      string
	peers   entity.Peers
}

func parseArgs() Args {
	var id string
	port := os.Args[1]
	apiPort := os.Args[2]
	peersList := []entity.Peer{}
	peers := entity.Peers{}

	data, _ := os.ReadFile("peers.json")
	json.Unmarshal(data, &peersList)

	for _, peer := range peersList {
		if peer.Port != port {
			peers[peer.Id] = peer
		} else {
			id = peer.Id
		}
	}

	return Args{port: port, apiPort: apiPort, peers: peers, id: id}
}

func init() {
	// Configure logrus
	log.SetLevel(log.DebugLevel)
	log.SetOutput(os.Stdout)
	log.SetFormatter(&log.TextFormatter{
		ForceColors:     true,
		DisableColors:   false,
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02T15:04:05.000-07:00",
	})

	log.Info("INIT")
}

func main() {
	args := parseArgs()

	// Config
	ELECTION_TIMEOUT := 300 // ms
	LEADER_TICKER := 25     // ms
	STATE_LOCATION := fmt.Sprintf("state_%s.json", args.id)

	// Driven port/adapter (domain -> infra)
	grpcClientAdapter := secondaryAdapter.NewGrpcClient()
	jsonAdapter := secondaryAdapter.NewJsonPersister()

	rpcClientPort := secondaryPort.NewRpcClientPort(grpcClientAdapter)
	persisterPort := secondaryPort.NewPersisterPort(STATE_LOCATION, jsonAdapter)

	// Domain
	persistent, _ := persisterPort.Load()
	timeout := entity.NewTimeout(ELECTION_TIMEOUT, LEADER_TICKER)
	srv := service.NewServer(args.id, args.peers, persistent)

	// Services
	runnerUsecase := runner.NewService(srv, timeout, rpcClientPort, persisterPort)
	receiverUsecase := receiver.NewService(srv)
	clusterUsecase := cluster.NewService(srv)

	// Driving port/adapter (infra -> domain)
	rpcServerPort := primaryPort.NewRpcServerPort(receiverUsecase)
	grpcServerAdapter := primaryAdapter.NewGrpcApi(rpcServerPort)

	// Infrastructure
	runnerServer := server.NewRunner(runnerUsecase)
	grpcServer := server.NewRpc(grpcServerAdapter)
	clusterServer := server.NewClusterServer(clusterUsecase)

	// Start servers
	go grpcServer.Start(args.port)
	go clusterServer.Start(args.apiPort)
	runnerServer.Start()
}
