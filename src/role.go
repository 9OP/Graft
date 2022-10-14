package main

type ROLE struct{}

type Voter interface {
	Vote(candidateId string)
	CanVote(candidateId string, candidateLastLogIndex int, candidateLastLogTerm int) bool
}

type FOLLOWER interface {
	Voter
	StartElection()
	Timeout()
}

/*
DDD
roles should not be attribute, but should be either struct or interface
role as interface: defined behaviour
role as struct: define state

behaviours:

Follower:
- vote
- start election
- raise to candidate

Candidate:
- vote
- restart election
- raise to leader
- downgrade to follower

Leader:
- broadcast
- downgrade to follower

*/
