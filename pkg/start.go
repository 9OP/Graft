package pkg

import (
	"graft/pkg/domain/entity"
	"graft/pkg/domain/service"
	primaryAdapter "graft/pkg/infrastructure/adapter/primary"
	secondaryAdapter "graft/pkg/infrastructure/adapter/secondary"
	primaryPort "graft/pkg/infrastructure/port/primary"
	secondaryPort "graft/pkg/infrastructure/port/secondary"
	"graft/pkg/infrastructure/server"

	"graft/pkg/usecase/cluster"
	"graft/pkg/usecase/receiver"
	"graft/pkg/usecase/runner"
	"os"

	log "github.com/sirupsen/logrus"
)

func configureLogger() {
	log.SetLevel(log.DebugLevel)
	log.SetOutput(os.Stdout)
	log.SetFormatter(&log.TextFormatter{
		ForceColors:     true,
		DisableColors:   false,
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02T15:04:05.000-07:00",
	})
}

func Start(
	id string,
	peers entity.Peers,
	persistentLocation string,
	electionTimeout int,
	heartbeatTimeout int,
	rpcPort string,
	apiPort string,
) {
	// Driven port/adapter (domain -> infra)
	grpcClientAdapter := secondaryAdapter.NewGrpcClient()
	jsonPersisterAdapter := secondaryAdapter.NewJsonPersister()

	rpcClientPort := secondaryPort.NewRpcClientPort(grpcClientAdapter)
	persisterPort := secondaryPort.NewPersisterPort(persistentLocation, jsonPersisterAdapter)

	// Domain
	persistent, _ := persisterPort.Load()
	timeout := entity.NewTimeout(electionTimeout, heartbeatTimeout)
	srv := service.NewServer(id, peers, persistent)

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

	// Start logger
	configureLogger()
	log.Info("START")

	// Start servers: p2p rpc, API and runner
	go grpcServer.Start(rpcPort)
	go clusterServer.Start(apiPort)
	runnerServer.Start()
}
