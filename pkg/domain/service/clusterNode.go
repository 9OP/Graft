package service

import (
	"graft/pkg/domain/entity"
	"sync"
	"sync/atomic"
	"unsafe"

	log "github.com/sirupsen/logrus"
)

/*
Implementation note:
- NodeState is supposed to only expose pure function `Withers` that return a copy
  instead of mutating the instance.
- When a mutation needs to be saved in ClusterNode (for instance incrementing NodeState.CurrentTerm)
  You should call SwapState, which atomically change the pointer of ClusterNode.NodeState to your
  new NodeState version.
- This means that every ClusterNode method that calls SwapState, needs to accept a pointer receiver (c *ClusterNode)
  Because otherwise, SwapState is called with a copy of the pointers instead of original pointers
*/

type ClusterNode struct {
	signals
	*entity.NodeState
}

func NewClusterNode(id string, peers entity.Peers, persistent *entity.PersistentState) *ClusterNode {
	nodeState := entity.NewNodeState(id, peers, persistent)
	signals := newSignals()
	return &ClusterNode{
		signals:   signals,
		NodeState: &nodeState,
	}
}

func (c *ClusterNode) GetState() entity.NodeState {
	return *c.NodeState
}

func (c *ClusterNode) SwapState(state interface{}) {
	var addr *unsafe.Pointer
	var new unsafe.Pointer

	switch newState := state.(type) {
	case *entity.PersistentState:
		oldState := &c.NodeState.FsmState.PersistentState
		addr = (*unsafe.Pointer)(unsafe.Pointer(oldState))
		new = unsafe.Pointer(newState)
	case *entity.FsmState:
		oldState := &c.NodeState.FsmState
		addr = (*unsafe.Pointer)(unsafe.Pointer(oldState))
		new = unsafe.Pointer(newState)
	case *entity.NodeState:
		oldState := &c.NodeState
		addr = (*unsafe.Pointer)(unsafe.Pointer(oldState))
		new = (unsafe.Pointer(newState))
	default:
		return
	}

	atomic.SwapPointer(addr, new)
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

func (c *ClusterNode) DeleteLogsFrom(index uint32) {
	newState := c.WithDeleteLogsFrom(index)
	c.SwapState(&newState)
	c.saveState()
}

func (c *ClusterNode) AppendLogs(prevLogIndex uint32, entries ...entity.LogEntry) {
	newState := c.WithAppendLogs(prevLogIndex, entries...)
	c.SwapState(&newState)
	c.saveState()
}

func (c *ClusterNode) SetClusterLeader(leaderId string) {
	if c.Leader().Id == leaderId {
		return
	}
	log.Debug("SET CLUSTER LEADER ", leaderId)
	newState := c.WithClusterLeader(leaderId)
	c.SwapState(&newState)
}

func (c *ClusterNode) SetCommitIndex(index uint32) {
	if c.CommitIndex() == index {
		return
	}
	log.Debug("SET COMMIT INDEX ", index)
	newState := c.WithCommitIndex(index)
	c.SwapState(&newState)
	c.saveState()
	c.commit()
}

func (c *ClusterNode) SetNextMatchIndex(peerId string, index uint32) {
	newState := c.
		WithNextIndex(peerId, index).
		WithMatchIndex(peerId, index)
	c.SwapState(&newState)
}

func (c *ClusterNode) DecrementNextIndex(peerId string) {
	newState := c.WithDecrementNextIndex(peerId)
	c.SwapState(&newState)
}

func (c *ClusterNode) GrantVote(peerId string) {
	newState := c.WithVotedFor(peerId)
	c.SwapState(&newState)
	c.saveState()
}

func (c *ClusterNode) DowngradeFollower(term uint32) {
	log.Infof("DOWNGRADE TO FOLLOWER TERM: %d\n", term)
	newState := c.
		WithRole(entity.Follower).
		WithCurrentTerm(term).
		WithVotedFor("")
	c.SwapState(&newState)
	c.shiftRole(entity.Follower)
	c.saveState()
	c.resetTimeout()
}

func (c *ClusterNode) IncrementCandidateTerm() {
	if c.Role() == entity.Candidate {
		term := c.CurrentTerm() + 1
		log.Debugf("INCREMENT CANDIDATE TERM %d\n", term)
		newState := c.
			WithCurrentTerm(term).
			WithVotedFor(c.Id())
		c.SwapState(&newState)
		c.saveState()
		c.resetTimeout()
		return
	}
	log.Warn("CANNOT INCREMENT TERM FOR ", c.Role())
}

func (c *ClusterNode) UpgradeCandidate() {
	if c.Role() == entity.Follower {
		log.Info("UPGRADE TO CANDIDATE")
		newState := c.WithRole(entity.Candidate)
		c.SwapState(&newState)
		c.shiftRole(entity.Candidate)
		return
	}
	log.Warn("CANNOT UPGRADE CANDIDATE FOR ", c.Role())
}

func (c *ClusterNode) UpgradeLeader() {
	if c.Role() == entity.Candidate {
		log.Infof("UPGRADE TO LEADER TERM %d\n", c.CurrentTerm())
		newState := c.
			WithInitializeLeader().
			WithClusterLeader(c.Id()).
			WithRole(entity.Leader)
		c.SwapState(&newState)
		c.shiftRole(entity.Leader)
		c.resetLeaderTicker()
		return
	}
	log.Warn("CANNOT UPGRADE LEADER FOR ", c.Role())
}

func (c *ClusterNode) ApplyLogs() {
	commitIndex := c.CommitIndex()
	state := *c.FsmState

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
