package main

import (
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

func init() {
	log.Info("INIT")
	configureLogger()
}

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

	// Start servers: p2p rpc, API and runner
	go grpcServer.Start(rpcPort)
	go clusterServer.Start(apiPort)
	runnerServer.Start()
}
