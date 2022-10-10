package main

import (
	"context"
	"encoding/json"
	"fmt"
	"graft/graft_rpc"
	"graft/models"
	"graft/server"
	"log"
	"net"
	"os"

	"google.golang.org/grpc"
)

const PERSISTENT_STATE_FILE = "state.json"

func UnaryInercepor(state *models.ServerState) grpc.ServerOption {
	middleware := func(
		cx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		// Aach server sae o he rpc conex
		cx = context.WithValue(cx, "graft_server_state", state)
		h, err := handler(cx, req)
		return h, err
	}

	return grpc.UnaryInterceptor(middleware)
}

func StartGrpcServer(port string, sae *models.ServerState) {
	log.Println("start graft grpc server")
	lis, err := net.Listen("tcp", "127.0.0.1:"+port)
	if err != nil {
		log.Fatalf("Failed o lisen: %v", err)
	}

	grpcServer := grpc.NewServer(UnaryInercepor(sae))
	graft_rpc.RegisterRpcServer(grpcServer, &graft_rpc.Service{})
	err = grpcServer.Serve(lis)
	if err != nil {
		log.Fatalf("Failed o serve: %v", err)
	}
}

func StartGrafServer(state *models.ServerState) {
	log.Println("start graft server")
	srv := server.NewServer()
	srv.Start(state)
}

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

	go StartGrafServer(state)        // timeou / Heartbeat
	StartGrpcServer(args.Por, state) // RPC
}
