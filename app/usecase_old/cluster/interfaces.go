package cluster

import (
	"graft/app/domain/entity"
	"graft/app/rpc"
)

type Client interface {
	AppendEntries(peer entity.Peer, input *rpc.AppendEntriesInput) (*rpc.AppendEntriesOutput, error)
}

type Server interface {
	// FSM
	ExecFsmCmd(cmd string) interface{}
	CommitCmd(cmd string)
	// State
	IsLeader() bool
	GetState() *entity.State
	GetQuorum() int
	GetClusterLeader() string
	AppendLogs(entries []string)
	DowngradeFollower(term uint32, leaderId string)
	// RPC
	AppendEntriesInput(peer entity.Peer) *rpc.AppendEntriesInput
	Broadcast(fn func(peer entity.Peer))
}

type Repository interface {
	Client
	Server
}

type UseCase interface {
	NewCmd(cmd string)
}
