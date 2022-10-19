package main

import (
	"encoding/json"
	"fmt"
	"graft/src2/entity"
	"graft/src2/repository"
	"graft/src2/repository/clients"
	"graft/src2/repository/servers"
	"graft/src2/usecase/persister"
	"graft/src2/usecase/receiver"
	"graft/src2/usecase/runner"
	"os"
)

type Args struct {
	port  string
	id    string
	peers []entity.Peer
}

func parseArgs() Args {
	var id string
	port := os.Args[1]
	peers := []entity.Peer{}

	daa, _ := os.ReadFile("peers.json")
	json.Unmarshal(daa, &peers)

	// Filter current host from nodes
	n := 0
	for _, peer := range peers {
		if peer.Port != port {
			peers[n] = peer
			n++
		} else {
			id = peer.Id
		}
	}
	peers = peers[:n]

	return Args{port: port, peers: peers, id: id}
}

func main() {
	args := parseArgs()

	var ELECTION_TIMEOUT int = 350 // ms
	var LEADER_TICKER int = 35     // ms
	timeout := entity.NewTimeout(ELECTION_TIMEOUT)
	ticker := entity.NewTicker(LEADER_TICKER)

	var stateService *persister.Service
	var runnerService *runner.Service
	var receiverService *receiver.Service

	var clientRunner = &clients.Runner{}
	var runnerServer *servers.Runner = &servers.Runner{}
	var receiverServer *servers.Receiver

	stateService = persister.NewService(fmt.Sprintf("state_%s.json", args.id), &repository.JsonPersister{})
	runnerService = runner.NewService(clientRunner, timeout, ticker)
	receiverService = receiver.NewService(runnerServer)

	*runnerServer = *servers.NewRunner(args.id, args.peers, stateService, timeout, ticker)
	receiverServer = servers.NewReceiver(args.port)

	go receiverServer.Start(receiverService)
	runnerServer.Start(runnerService)
}
