package service

import (
	"fmt"
	"os/exec"
	"sync"
	"sync/atomic"
	"unsafe"

	"graft/pkg/domain/entity"

	log "github.com/sirupsen/logrus"
)

/*
Implementation note:
- NodeState is supposed to only expose pure function `Withers` that return a copy
  instead of mutating the instance.
- When a mutation needs to be saved in ClusterNode (for instance incrementing NodeState.CurrentTerm)
  You should call swapState, which atomically change the pointer of ClusterNode.NodeState to your
  new NodeState version.
- This means that every ClusterNode method that calls swapState, needs to accept a pointer receiver (c *ClusterNode)
  Because otherwise, swapState is called with a copy of the pointers instead of original pointers
*/

type ClusterNode struct {
	*entity.NodeState
	signals
	fsmInit string
	fsmEval string
}

func NewClusterNode(
	id string,
	peers entity.Peers,
	fsmInit string,
	fsmEval string,
	persistent *entity.PersistentState,
) *ClusterNode {
	nodeState := entity.NewNodeState(id, peers, persistent)
	signals := newSignals()
	return &ClusterNode{
		NodeState: &nodeState,
		signals:   signals,
		fsmInit:   fsmInit,
		fsmEval:   fsmEval,
	}
}

func (c ClusterNode) GetState() entity.NodeState {
	return *c.NodeState
}

func (c *ClusterNode) swapState(state interface{}) {
	var addr *unsafe.Pointer
	var new unsafe.Pointer

	switch newState := state.(type) {
	case *entity.PersistentState:
		oldState := &c.NodeState.FsmState.PersistentState
		addr = (*unsafe.Pointer)(unsafe.Pointer(oldState))
		new = unsafe.Pointer(newState)
		defer c.saveState()
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
	if newState, changed := c.WithDeleteLogsFrom(index); changed {
		c.swapState(&newState)
	}
}

func (c *ClusterNode) AppendLogs(prevLogIndex uint32, entries ...entity.LogEntry) {
	if newState, changed := c.WithAppendLogs(prevLogIndex, entries...); changed {
		log.Debug("APPEND LOGS")
		c.swapState(&newState)
	}
}

func (c *ClusterNode) SetClusterLeader(leaderId string) {
	if c.Leader().Id != leaderId {
		log.Debug("SET CLUSTER LEADER ", leaderId)
		newState := c.WithClusterLeader(leaderId)
		c.swapState(&newState)
	}
}

func (c *ClusterNode) SetCommitIndex(index uint32) {
	if newState, changed := c.WithCommitIndex(index); changed {
		log.Debug("SET COMMIT INDEX ", index)
		c.swapState(&newState)
		c.commit()
	}
}

func (c *ClusterNode) SetNextMatchIndex(peerId string, index uint32) {
	newState := c.
		WithNextIndex(peerId, index).
		WithMatchIndex(peerId, index)
	c.swapState(&newState)
}

func (c *ClusterNode) DecrementNextIndex(peerId string) {
	if newState, changed := c.WithDecrementNextIndex(peerId); changed {
		c.swapState(&newState)
	}
}

func (c *ClusterNode) GrantVote(peerId string) {
	if c.VotedFor() != peerId {
		newState := c.WithVotedFor(peerId)
		c.swapState(&newState)
	}
}

func (c *ClusterNode) DowngradeFollower(term uint32) {
	log.Infof("DOWNGRADE TO FOLLOWER TERM: %d\n", term)
	newRole := entity.Follower
	newState := c.
		WithRole(newRole).
		WithCurrentTerm(term).
		WithVotedFor("")
	c.swapState(&newState)
	c.shiftRole(newRole)
	c.resetTimeout()
}

func (c *ClusterNode) IncrementCandidateTerm() {
	if c.Role() == entity.Candidate {
		term := c.CurrentTerm() + 1
		log.Debugf("INCREMENT CANDIDATE TERM %d\n", term)
		newState := c.
			WithCurrentTerm(term).
			WithVotedFor(c.Id())
		c.swapState(&newState)
		c.resetTimeout()
		return
	}
	log.Warn("CANNOT INCREMENT TERM FOR ", c.Role())
}

func (c *ClusterNode) UpgradeCandidate() {
	if c.Role() == entity.Follower {
		log.Info("UPGRADE TO CANDIDATE")
		newRole := entity.Candidate
		newState := c.WithRole(newRole)
		c.swapState(&newState)
		c.shiftRole(newRole)
		c.resetTimeout()
		return
	}
	log.Warn("CANNOT UPGRADE CANDIDATE FOR ", c.Role())
}

func (c *ClusterNode) UpgradeLeader() {
	if c.Role() == entity.Candidate {
		log.Infof("UPGRADE TO LEADER TERM %d\n", c.CurrentTerm())
		newRole := entity.Leader
		newState := c.
			WithInitializeLeader().
			WithClusterLeader(c.Id()).
			WithRole(newRole)
		c.swapState(&newState)
		c.shiftRole(newRole)
		c.resetLeaderTicker()
		return
	}
	log.Warn("CANNOT UPGRADE LEADER FOR ", c.Role())
}

func (c *ClusterNode) ApplyLogs() {
	state := *c.FsmState
	commitIndex := c.CommitIndex()

	for state.LastApplied() < commitIndex {
		if state.LastApplied() == 0 {
			c.initFsm()
		}
		// Increment last applied first
		// because lastApplied = 0 is not a valid logEntry
		state = state.WithIncrementLastApplied()
		if log, err := c.MachineLog(state.LastApplied()); err == nil {
			res := c.evalFsm(log.Value, "COMMAND")
			if log.C != nil {
				log.C <- res
			}
		}
	}
	// Swap only once
	c.swapState(&state)
}

func (c *ClusterNode) ExecuteCommand(command string) chan entity.EvalResult {
	result := make(chan entity.EvalResult, 1)
	newEntry := entity.LogEntry{
		Value: command,
		Term:  c.CurrentTerm(),
		C:     result,
	}
	go c.AppendLogs(c.LastLogIndex(), newEntry)
	return result
}

func (c ClusterNode) ExecuteQuery(query string) chan entity.EvalResult {
	result := make(chan entity.EvalResult, 1)
	go (func() { result <- c.evalFsm(query, "QUERY") })()
	return result
}

func (c ClusterNode) evalFsm(entry string, entryType string) entity.EvalResult {
	cmd := exec.Command(
		c.fsmEval,
		entry,
		entryType,
		c.Id(),
		c.LeaderId(),
		c.Role().String(),
		fmt.Sprint(c.LastLogIndex()),
		fmt.Sprint(c.CurrentTerm()),
		fmt.Sprint(c.VotedFor()),
	)
	out, err := cmd.Output()
	log.Debugf("EVAL:\n\t%s\n\t%s\n\t%s", entry, string(out), err.Error())
	return entity.EvalResult{
		Out: out,
		Err: err,
	}
}

func (c ClusterNode) initFsm() {
	cmd := exec.Command(c.fsmInit)
	out, err := cmd.Output()
	log.Debugf("EVAL:\n\t%s\n\t%s", string(out), err.Error())
}
