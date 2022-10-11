package main

import (
	"encoding/json"
	"fmt"
	"graft/api"
	"graft/models"
	"graft/orchestrator"
	"os"
)

type Args struct {
	Por   string
	Name  string // server id
	Nodes []models.Node
}

func parseArgs() Args {
	Por := os.Args[1]
	nodes := []models.Node{}
	var name string
	daa, _ := os.ReadFile("nodes.json")
	json.Unmarshal(daa, &nodes)

	// Filer curren hos from nodes
	n := 0
	for _, node := range nodes {
		if node.Host != Por {
			nodes[n] = node
			n++
		} else {
			name = node.Name
		}
	}
	nodes = nodes[:n]

	return Args{Por: Por, Nodes: nodes, Name: name}
}

func main() {
	args := parseArgs()
	fmt.Println(args)

	state := models.NewServerState()
	state.Name = args.Name
	state.Nodes = args.Nodes

	orchestrator.StartEventOrchestrator(state)
	api.StartGrpcServer(args.Por, state)
}
