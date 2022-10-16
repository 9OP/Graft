package raft_rpc

import (
	"graft/src/entity"
	"graft/src/rpc"
)

type Server interface {
	GetState() entity.State
	Heartbeat()
	GrantVote()
	DowngradeFollower()
}

type Repository interface {
	Server
}

type UseCase interface {
	AppendEntries(input *rpc.AppendEntriesInput) (*rpc.AppendEntriesOutput, error)
	RequestVote(input *rpc.RequestVoteInput) (*rpc.RequestVoteOutput, error)
}
