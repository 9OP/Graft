package main

import (
	"fmt"
	"graft/app/domain/entity"
	"graft/app/domain/service"
	primaryAdapter "graft/app/infrastructure/adapter/primary"
	secondaryAdapter "graft/app/infrastructure/adapter/secondary"
	primaryPort "graft/app/infrastructure/port/primary"
	secondaryPort "graft/app/infrastructure/port/secondary"
	"graft/app/infrastructure/server"
	"graft/app/infrastructure/utils"
	"graft/app/usecase/cluster"
	"graft/app/usecase/receiver"
	"graft/app/usecase/runner"

	log "github.com/sirupsen/logrus"
)

func init() {
	log.Info("INIT")
	utils.ConfigureLogger()
}

func boot(config utils.Config) {
	// Driven port/adapter (domain -> infra)
	grpcClientAdapter := secondaryAdapter.NewGrpcClient()
	jsonAdapter := secondaryAdapter.NewJsonPersister()

	rpcClientPort := secondaryPort.NewRpcClientPort(grpcClientAdapter)
	persisterPort := secondaryPort.NewPersisterPort(STATE_LOCATION, jsonAdapter)

	// Domain
	persistent, _ := persisterPort.Load()
	timeout := entity.NewTimeout(config.Timeouts.Election, config.Timeouts.Heartbeat)
	srv := service.NewServer(args.id, config.Peers, persistent)

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

func main() {
	config, err := utils.GetConfig("config.yml")
	if err != nil {
		log.Fatalf("Cannot load config \n\t%s", err)
	}

	fmt.Println(config)
}
