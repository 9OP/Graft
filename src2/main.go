package main

import (
	"graft/src2/servers"
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

	stateService = persister.NewService("state.json", &JsonPersister{})
	runnerService = runner.NewService()
	receiverService = receiver.NewService(runnerServer)

	runnerServer = servers.NewRunner(
		"id",
		nil,
		stateService,
		runnerService,
	)
	receiverServer = servers.NewReceiver(
		receiverService,
		"10000",
	)

	go receiverServer.Start()
	runnerServer.Start()
}
