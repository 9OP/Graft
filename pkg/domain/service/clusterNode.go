package service

import (
	"graft/pkg/domain/entity"
	"sync"
	"sync/atomic"
	"unsafe"
)

type ClusterNode struct {
	signals
	entity.NodeState
}

func NewClusterNode(id string, peers entity.Peers, persistent entity.PersistentState) ClusterNode {
	nodeState := entity.NewNodeState(id, peers, persistent)
	signals := newSignals()
	return ClusterNode{
		signals:   signals,
		NodeState: nodeState,
	}
}

func (c ClusterNode) GetState() entity.NodeState {
	return c.NodeState
}

func (c ClusterNode) SwapState(state interface{}) {
	var addr, new unsafe.Pointer

	// Using reflection to get the addr pointer for
	// atomic swap
	switch newState := state.(type) {
	case *entity.PersistentState:
		addr = unsafe.Pointer(&c.NodeState.FsmState.PersistentState)
		new = unsafe.Pointer(newState)
	case *entity.FsmState:
		addr = unsafe.Pointer(&c.NodeState.FsmState)
		new = unsafe.Pointer(newState)
	case *entity.NodeState:
		addr = unsafe.Pointer(&c.NodeState)
		new = unsafe.Pointer(newState)
	default:
		return
	}

	atomic.SwapPointer(&addr, new)
}

func (c ClusterNode) Heartbeat() {
	c.resetTimeout()
}

func (c ClusterNode) Broadcast(fn func(p entity.Peer)) {
	var wg sync.WaitGroup
	for _, peer := range c.Peers() {
		wg.Add(1)
		go func(p entity.Peer, w *sync.WaitGroup) {
			defer w.Done()
			fn(p)
		}(peer, &wg)
	}
	wg.Wait()
}

func (c ClusterNode) DeleteLogsFrom(index uint32) {
	newState := c.WithDeleteLogsFrom(index)
	c.SwapState(&newState)
	c.saveState()
}

func (c ClusterNode) AppendLogs(prevLogIndex uint32, entries ...entity.LogEntry) {
	newState := c.WithAppendLogs(prevLogIndex, entries...)
	c.SwapState(&newState)
	c.saveState()
}

func (c ClusterNode) SetClusterLeader(leaderId string) {
	newState := c.WithClusterLeader(leaderId)
	c.SwapState(&newState)
}

func (c ClusterNode) SetCommitIndex(index uint32) {
	newState := c.WithCommitIndex(index)
	c.SwapState(&newState)
	c.saveState()
	c.commit()
}

func (c ClusterNode) SetNextMatchIndex(peerId string, index uint32) {
	newState := c.
		WithNextIndex(peerId, index).
		WithMatchIndex(peerId, index)
	c.SwapState(&newState)
}

func (c ClusterNode) DecrementNextIndex(peerId string) {
	newState := c.WithDecrementNextIndex(peerId)
	c.SwapState(&newState)
}

func (c ClusterNode) GrantVote(peerId string) {
	newState := c.WithVotedFor(peerId)
	c.SwapState(&newState)
	c.saveState()
}

func (c ClusterNode) DowngradeFollower(term uint32) {
	// log.Infof("DOWNGRADE TO FOLLOWER TERM: %d\n", term)
	newState := c.
		WithRole(entity.Follower).
		WithCurrentTerm(term).
		WithVotedFor("")
	c.SwapState(&newState)
	c.saveState()
	c.resetTimeout()
}

func (c ClusterNode) IncrementCandidateTerm() {
	if c.Role() == entity.Candidate {
		term := c.CurrentTerm() + 1
		// log.Debugf("INCREMENT CANDIDATE TERM %d\n", term)
		newState := c.
			WithCurrentTerm(term).
			WithVotedFor(c.Id())
		c.SwapState(&newState)
		c.saveState()
		c.resetTimeout()
		return
	}
	// log.Warn("CANNOT INCREMENT TERM FOR ", s.Role)
}

func (c ClusterNode) UpgradeCandidate() {
	if c.Role() == entity.Follower {
		// 		log.Info("UPGRADE TO CANDIDATE")
		newState := c.WithRole(entity.Candidate)
		c.SwapState(&newState)
		return
	}
	// 	log.Warn("CANNOT UPGRADE CANDIDATE FOR ", s.Role)

}

func (c ClusterNode) UpgradeLeader() {
	if c.Role() == entity.Leader {
		// 		log.Infof("UPGRADE TO LEADER TERM %d\n", s.CurrentTerm)
		newState := c.
			WithInitializeLeader().
			WithClusterLeader(c.Id()).
			WithRole(entity.Leader)
		c.SwapState(&newState)
		c.resetLeaderTicker()
		return
	}
	// 	log.Warn("CANNOT UPGRADE LEADER FOR ", s.Role)
}

func (c ClusterNode) ApplyLogs() {
	commitIndex := c.CommitIndex()
	state := c.FsmState

	for lastApplied := state.LastApplied(); commitIndex > lastApplied; {
		state = state.WithIncrementLastApplied()
		if log, err := state.MachineLog(lastApplied); err == nil {
			res := eval(log.Value)
			if log.C != nil {
				log.C <- res
			}
		}
	}
	// Swap only once
	c.SwapState(&state)
}

func (c ClusterNode) ExecuteCommand(command string) chan interface{} {
	result := make(chan interface{}, 1)
	newEntry := entity.LogEntry{
		Value: command,
		Term:  c.CurrentTerm(),
		C:     result,
	}
	c.AppendLogs(c.LastLogIndex(), newEntry)
	return result
}

func (c ClusterNode) ExecuteQuery(query string) chan interface{} {
	result := make(chan interface{}, 1)
	go (func() { result <- eval(query) })()
	return result
}

func eval(entry string) interface{} {
	// log.Debug("EXECUTE: ", entry)
	// cmd := exec.Command(entry)
	// var outb, errb bytes.Buffer
	// cmd.Stdout = &outb
	// cmd.Stderr = &errb
	// err := cmd.Run()

	// if err != nil {
	// 	return errb.Bytes()
	// }
	// return outb.Bytes()
	return entry
}
