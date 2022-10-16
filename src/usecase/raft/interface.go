package raft

import (
	"graft/src/entity"
	"graft/src/rpc"
)

type Server interface {
	getState() entity.State
	heartbeat()
	grantVote()
	downgradeFollower()
}

type Repository interface {
	Server
}

type UseCase interface {
	AppendEntries(input *rpc.AppendEntriesInput) (*rpc.AppendEntriesOutput, error)
	RequestVote(input *rpc.RequestVoteInput) (*rpc.RequestVoteOutput, error)
}
