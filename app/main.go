package main

import (
	"encoding/json"
	"graft/app/domain/entity"
	"graft/app/domain/service"
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
	_ = parseArgs()

	var ELECTION_TIMEOUT int = 350 // ms
	var LEADER_TICKER int = 35     // ms
	_ = entity.NewTimeout(ELECTION_TIMEOUT, LEADER_TICKER)

	_ = service.Server{}

}
