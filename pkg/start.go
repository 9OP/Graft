package pkg

import (
	"graft/pkg/domain"

	primaryAdapter "graft/pkg/infrastructure/adapter/primary"
	secondaryAdapter "graft/pkg/infrastructure/adapter/secondary"
	primaryPort "graft/pkg/infrastructure/port/primary"
	secondaryPort "graft/pkg/infrastructure/port/secondary"
	"graft/pkg/infrastructure/server"

	"graft/pkg/services/api"
	"graft/pkg/services/core"
	"graft/pkg/services/rpc"
	"graft/pkg/utils"
)

func Start(
	id string,
	peers domain.Peers,
	persistentLocation string,
	electionTimeout int,
	heartbeatTimeout int,
	rpcPort string,
	apiPort string,
	fsmInit string,
	fsmEval string,
	logLevel string,
) {
	utils.ConfigureLogger(logLevel)

	// Driven port/adapter (domain -> infra)
	grpcClientAdapter := secondaryAdapter.NewGrpcClient()
	jsonPersisterAdapter := secondaryAdapter.NewJsonPersister()

	rpcClientPort := secondaryPort.NewRpcClientPort(grpcClientAdapter)
	persisterPort := secondaryPort.NewPersisterPort(persistentLocation, jsonPersisterAdapter)

	// Domain
	persistent, _ := persisterPort.Load()
	node := domain.NewNode(id, peers, persistent)

	// Services
	coreService := core.NewService(node, rpcClientPort, persisterPort, electionTimeout, heartbeatTimeout)
	rpcService := rpc.NewService(node)
	apiService := api.NewService(node)

	// Driving port/adapter (infra -> domain)
	rpcServerPort := primaryPort.NewRpcServerPort(rpcService)
	grpcServerAdapter := primaryAdapter.NewGrpcApi(rpcServerPort)

	// Infrastructure
	runnerServer := server.NewRunner(coreService)
	grpcServer := server.NewRpc(grpcServerAdapter)
	clusterServer := server.NewClusterServer(apiService)

	// Start servers: p2p rpc, API and runner
	go grpcServer.Start(rpcPort)
	go clusterServer.Start(apiPort)
	runnerServer.Start()
}
