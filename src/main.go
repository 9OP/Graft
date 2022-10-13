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
	nodes := []Node{}
	var name string
	daa, _ := os.ReadFile("nodes.json")
	json.Unmarshal(daa, &nodes)

	// Filtert current host from nodes
	n := 0
	for _, node := range nodes {
		if node.Host != port {
			nodes[n] = node
			n++
		} else {
			name = node.Name
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

	srv := NewGrpcServer(state)
	srv.Start(args.port)
}
