package main

import (
	"context"
	"graft/graft_rpc"

	"log"

	"google.golang.org/grpc"
)

func main() {
	var conn *grpc.ClientConn

	conn, err := grpc.Dial(":7679", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect %v", err)
	}
	defer conn.Close()

	c := graft_rpc.NewRpcClient(conn)

	_, err = c.AppendEntries(context.Background(), &graft_rpc.AppendEntriesInput{})
	if err != nil {
		log.Fatalf("Error when calling AppendEntries: %s", err)
	}
	log.Println("received response from server")
}
