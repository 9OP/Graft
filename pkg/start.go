package pkg

import (
	"net/netip"

	"graft/pkg/domain"

	primaryAdapter "graft/pkg/infrastructure/adapter/primary"
	secondaryAdapter "graft/pkg/infrastructure/adapter/secondary"
	primaryPort "graft/pkg/infrastructure/port/primary"
	secondaryPort "graft/pkg/infrastructure/port/secondary"
	"graft/pkg/infrastructure/server"

	"graft/pkg/services/core"
	"graft/pkg/services/rpc"
	"graft/pkg/utils"
)

func Start(
	id string,
	host netip.AddrPort,
	peers domain.Peers,
	persistentLocation string,
	electionTimeout int,
	leaderHeartbeat int,
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
	config := domain.NodeConfig{
		Id:              id,
		Host:            host,
		ElectionTimeout: electionTimeout,
		LeaderHeartbeat: leaderHeartbeat,
	}
	node := domain.NewNode(config, peers, persistent)

	// Services
	coreService := core.NewService(node, rpcClientPort, persisterPort)
	rpcService := rpc.NewService(node)

	// Driving port/adapter (infra -> domain)
	rpcServerPort := primaryPort.NewRpcServerPort(rpcService)
	grpcServerAdapter := primaryAdapter.NewGrpcApi(rpcServerPort)

	// Infrastructure
	core := server.NewRunner(coreService)
	rpc := server.NewRpc(grpcServerAdapter)

	// Start servers: p2p rpc, API and runner
	go rpc.Start(host.Port())
	core.Start()
}
