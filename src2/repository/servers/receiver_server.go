package servers

import (
	"fmt"
	"graft/src2/rpc"
	"graft/src2/usecase/receiver"
	"log"
	"net"

	"google.golang.org/grpc"
)

type Receiver struct {
	service *receiver.Service
	port    string
}

func NewReceiver(service *receiver.Service, port string) *Receiver {
	return &Receiver{service: service, port: port}
}

func (r *Receiver) Start() {
	log.Println("START RECEIVER SERVER")

	addr := fmt.Sprintf("%s:%s", "127.0.0.1", r.port)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("Failed to listen: \n\t%v\n", err)
	}

	server := grpc.NewServer()
	rpc.RegisterRpcServer(server, r.service)

	if err := server.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: \n\t%v\n", err)
	}
}
