package main

import (
	"encoding/json"
	"os"
)

type Args struct {
	port  string
	name  string // server id
	nodes []Node
}

func parseArgs() Args {
	port := os.Args[1]
	host := os.Args[2]
	nodes := []Node{}
	var name string
	daa, _ := os.ReadFile("nodes.json")
	json.Unmarshal(daa, &nodes)

	// Filter current host from nodes
	n := 0
	for _, node := range nodes {
		if node.host != host || node.port != port {
			nodes[n] = node
			n++
		} else {
			name = node.id
		}
	}
	nodes = nodes[:n]

	return Args{port: port, nodes: nodes, name: name}
}

func main() {
	args := parseArgs()
	state := NewServerState(args.name, &args.nodes)

	och := NewEventOrchestrator(state)
	go och.Start()

	srv := NewRpcServer(state)
	srv.Start(args.port)
}
