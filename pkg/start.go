package pkg

import (
	"graft/pkg/domain"

	primaryAdapter "graft/pkg/infrastructure/adapter/primary"
	secondaryAdapter "graft/pkg/infrastructure/adapter/secondary"
	primaryPort "graft/pkg/infrastructure/port/primary"
	secondaryPort "graft/pkg/infrastructure/port/secondary"
	"graft/pkg/infrastructure/server"

	"graft/pkg/usecase/cluster"
	"graft/pkg/usecase/receiver"
	"graft/pkg/usecase/runner"
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
	runnerUsecase := runner.NewService(node, rpcClientPort, persisterPort, electionTimeout, heartbeatTimeout)
	receiverUsecase := receiver.NewService(node)
	clusterUsecase := cluster.NewService(node)

	// Driving port/adapter (infra -> domain)
	rpcServerPort := primaryPort.NewRpcServerPort(receiverUsecase)
	grpcServerAdapter := primaryAdapter.NewGrpcApi(rpcServerPort)

	// Infrastructure
	runnerServer := server.NewRunner(runnerUsecase)
	grpcServer := server.NewRpc(grpcServerAdapter)
	clusterServer := server.NewClusterServer(clusterUsecase)

	// Start servers: p2p rpc, API and runner
	go grpcServer.Start(rpcPort)
	go clusterServer.Start(apiPort)
	runnerServer.Start()
}
