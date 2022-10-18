package main

import (
	"graft/src2/repository"
	"graft/src2/repository/clients"
	"graft/src2/repository/servers"
	"graft/src2/usecase/persister"
	"graft/src2/usecase/receiver"
	"graft/src2/usecase/runner"
)

func main() {
	// Parse args

	var stateService *persister.Service
	var runnerService *runner.Service
	var receiverService *receiver.Service
	var runnerServer *servers.Runner
	var receiverServer *servers.Receiver

	stateService = persister.NewService("state.json", &repository.JsonPersister{})
	runnerService = runner.NewService(&clients.Runner{})
	receiverService = receiver.NewService(runnerServer)

	id := "id"

	runnerServer = servers.NewRunner(id, nil, stateService)
	receiverServer = servers.NewReceiver("10000")

	go receiverServer.Start(receiverService)
	runnerServer.Start(runnerService)
}
