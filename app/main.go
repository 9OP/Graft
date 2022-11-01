package main

import (
	"encoding/json"
	"fmt"
	"graft/app/domain/entity"
	"graft/app/infrastructure/utils"
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
	log.Info("INIT")
	utils.ConfigureLogger()
}

func main() {
	config, err := utils.GetConfig("config.yml")
	if err != nil {
		log.Fatalf("Cannot load config \n\t%s", err)
	}

	fmt.Println(config)
}

// func main() {
// 	fmt.Println(getConfig())
// 	args := parseArgs()

// 	// Config
// 	ELECTION_TIMEOUT := 300 // ms
// 	LEADER_TICKER := 25     // ms
// 	STATE_LOCATION := fmt.Sprintf("state_%s.json", args.id)

// 	// Driven port/adapter (domain -> infra)
// 	grpcClientAdapter := secondaryAdapter.NewGrpcClient()
// 	jsonAdapter := secondaryAdapter.NewJsonPersister()

// 	rpcClientPort := secondaryPort.NewRpcClientPort(grpcClientAdapter)
// 	persisterPort := secondaryPort.NewPersisterPort(STATE_LOCATION, jsonAdapter)

// 	// Domain
// 	persistent, _ := persisterPort.Load()
// 	timeout := entity.NewTimeout(ELECTION_TIMEOUT, LEADER_TICKER)
// 	srv := service.NewServer(args.id, args.peers, persistent)

// 	// Services
// 	runnerUsecase := runner.NewService(srv, timeout, rpcClientPort, persisterPort)
// 	receiverUsecase := receiver.NewService(srv)
// 	clusterUsecase := cluster.NewService(srv)

// 	// Driving port/adapter (infra -> domain)
// 	rpcServerPort := primaryPort.NewRpcServerPort(receiverUsecase)
// 	grpcServerAdapter := primaryAdapter.NewGrpcApi(rpcServerPort)

// 	// Infrastructure
// 	runnerServer := server.NewRunner(runnerUsecase)
// 	grpcServer := server.NewRpc(grpcServerAdapter)
// 	clusterServer := server.NewClusterServer(clusterUsecase)

// 	// Start servers
// 	go grpcServer.Start(args.port)
// 	go clusterServer.Start(args.apiPort)
// 	runnerServer.Start()
// }
