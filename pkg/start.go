package pkg

import (
	"fmt"
	"net/netip"

	"graft/pkg/domain"
	"graft/pkg/utils"
	"graft/pkg/utils/log"

	primaryAdapter "graft/pkg/infrastructure/adapter/primary"
	secondaryAdapter "graft/pkg/infrastructure/adapter/secondary"
	primaryPort "graft/pkg/infrastructure/port/primary"
	secondaryPort "graft/pkg/infrastructure/port/secondary"
	"graft/pkg/infrastructure/server"

	"graft/pkg/services/core"
	"graft/pkg/services/rpc"
)

func Start(
	id string,
	host netip.AddrPort,
	peers domain.Peers,
	fsm string,
	electionTimeout int,
	leaderHeartbeat int,
) chan struct{} {
	quit := make(chan struct{}, 1)

	// Driven port/adapter (domain -> infra)
	grpcClientAdapter := secondaryAdapter.NewClusterClient()
	jsonPersisterAdapter := secondaryAdapter.NewJsonPersister()

	rpcClientPort := secondaryPort.NewRpcClientPort(grpcClientAdapter)
	persisterPort := secondaryPort.NewPersisterPort(fmt.Sprintf("%s/.%s.json", utils.GraftPath(), id), jsonPersisterAdapter)

	// Domain
	persistent, _ := persisterPort.Load()
	config := domain.NodeConfig{
		Id:              id,
		Host:            host,
		ElectionTimeout: electionTimeout,
		LeaderHeartbeat: leaderHeartbeat,
	}
	node := domain.NewNode(config, fsm, peers, persistent)

	// Services
	coreService := core.NewService(node, rpcClientPort, persisterPort)
	rpcService := rpc.NewService(node, rpcClientPort, quit)

	// Driving port/adapter (infra -> domain)
	rpcServerPort := primaryPort.NewRpcServerPort(rpcService)
	grpcServerAdapter := primaryAdapter.NewGrpcApi(rpcServerPort)

	// Infrastructure
	core := server.NewRunner(coreService)
	rpc := server.NewRpc(grpcServerAdapter, host.Port())

	// Start servers
	go (func() {
		err := rpc.Start()
		if err != nil {
			log.Fatalf("cannot start rpc server: %v", err)
		}
	})()
	go core.Start()

	log.Infof("start")

	return quit
}
