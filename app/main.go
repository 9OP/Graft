package main

import (
	"encoding/json"
	"fmt"
	"graft/app/domain/entity"
	"graft/app/domain/service"
	"graft/app/infrastructure/adapter"
	"graft/app/infrastructure/port"
	"graft/app/infrastructure/server"
	rpcReceiver "graft/app/usecase/rpcReceiver"
	rpcSender "graft/app/usecase/rpcSender"
	"os"
)

type Args struct {
	port  string
	id    string
	peers entity.Peers
}

func parseArgs() Args {
	var id string
	port := os.Args[1]
	peersList := []entity.Peer{}
	peers := entity.Peers{}

	data, _ := os.ReadFile("peers.json")
	json.Unmarshal(data, &peersList)

	for _, peer := range peersList {
		if peer.Port != port {
			peers[peer.Id] = peer
		} else {
			id = peer.Id
		}
	}

	return Args{port: port, peers: peers, id: id}
}

func main() {
	args := parseArgs()

	var ELECTION_TIMEOUT int = 350 // ms
	var LEADER_TICKER int = 35     // ms
	timeout := entity.NewTimeout(ELECTION_TIMEOUT, LEADER_TICKER)
	srv := service.NewServer(args.id, args.peers)

	persister := adapter.NewJsonPersister(fmt.Sprintf("state_%s.json", args.id))
	runner := server.NewRunnerServer(srv, timeout, persister)

	// RPC Client port/adapters
	rpcClient := adapter.NewRpcClient()
	rpcClientPort := port.NewRpcClientPort(rpcClient)

	// Use cases
	rpcReceiver := rpcReceiver.NewService(srv)
	rpcSender := rpcSender.NewService(rpcClientPort, timeout)

	// RPC Server port/adapters
	rpcServerPort := port.NewRpcServerPort(rpcReceiver)
	rpcServer := adapter.NewRpcApi(rpcServerPort)

	server := server.NewRpcServer(rpcServer)

	go server.Start(args.port)
	runner.Start(rpcSender)

}
