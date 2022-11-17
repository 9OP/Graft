package state

import (
	"bytes"
	"errors"
	"fmt"
	"os/exec"
	"sync"
	"sync/atomic"
	"unsafe"

	"graft/pkg/domain"

	log "github.com/sirupsen/logrus"
)

/*
Implementation note:
- nodeState is supposed to only expose pure function `Withers` that return a copy
  instead of mutating the instance.
- When a mutation needs to be saved in ClusterNode (for instance incrementing nodeState.CurrentTerm)
  You should call swapState, which atomically change the pointer of ClusterNode.nodeState to your
  new nodeState version.
- This means that every ClusterNode method that calls swapState, needs to accept a pointer receiver (c *ClusterNode)
  Because otherwise, swapState is called with a copy of the pointers instead of original pointers
*/

type signals struct {
	Commit             chan struct{}
	ShiftRole          chan struct{}
	SaveState          chan struct{}
	SynchronizeLogs    chan struct{}
	ResetLeaderTicker  chan struct{}
	ResetElectionTimer chan struct{}
}

func newSignals() signals {
	return signals{
		Commit:             make(chan struct{}, 1),
		ShiftRole:          make(chan struct{}, 1),
		SaveState:          make(chan struct{}, 1),
		SynchronizeLogs:    make(chan struct{}, 1),
		ResetLeaderTicker:  make(chan struct{}, 1),
		ResetElectionTimer: make(chan struct{}, 1),
	}
}

func (s *signals) commit() {
	s.Commit <- struct{}{}
}

func (s *signals) shiftRole() {
	s.ShiftRole <- struct{}{}
}

func (s *signals) saveState() {
	s.SaveState <- struct{}{}
}

func (s *signals) synchronizeLogs() {
	s.SynchronizeLogs <- struct{}{}
}

func (s *signals) resetLeaderTicker() {
	s.ResetLeaderTicker <- struct{}{}
}

func (s *signals) resetTimeout() {
	s.ResetElectionTimer <- struct{}{}
}

type ClusterNode struct {
	*nodeState
	signals
	fsmInit string
	fsmEval string
}

func NewClusterNode(
	id string,
	peers domain.Peers,
	fsmInit string,
	fsmEval string,
	persistent *PersistentState,
) *ClusterNode {
	return &ClusterNode{
		nodeState: NewNodeState(id, peers, persistent),
		signals:   newSignals(),
		fsmInit:   fsmInit,
		fsmEval:   fsmEval,
	}
}

func (c ClusterNode) GetState() nodeState {
	return *c.nodeState
}

func (c *ClusterNode) swapState(state interface{}) {
	var addr *unsafe.Pointer
	var new unsafe.Pointer

	switch newState := state.(type) {
	case *PersistentState:
		oldState := &c.nodeState.fsmState.PersistentState
		addr = (*unsafe.Pointer)(unsafe.Pointer(oldState))
		new = unsafe.Pointer(newState)
	case *fsmState:
		oldState := &c.nodeState.fsmState
		addr = (*unsafe.Pointer)(unsafe.Pointer(oldState))
		new = unsafe.Pointer(newState)
	case *nodeState:
		oldState := &c.nodeState
		addr = (*unsafe.Pointer)(unsafe.Pointer(oldState))
		new = (unsafe.Pointer(newState))
	default:
		return
	}

	atomic.SwapPointer(addr, new)
	c.saveState()
}

func (c ClusterNode) Heartbeat() {
	c.resetTimeout()
}

func (c ClusterNode) Broadcast(fn func(p domain.Peer)) {
	var wg sync.WaitGroup
	for _, peer := range c.Peers() {
		wg.Add(1)
		go func(p domain.Peer, w *sync.WaitGroup) {
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

func (c *ClusterNode) AppendLogs(prevLogIndex uint32, entries ...domain.LogEntry) {
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
		log.Debug("SET COMMIT INDEX ")
		c.swapState(&newState)
		c.commit()
	}
}

func (c *ClusterNode) SetNextMatchIndex(peerId string, index uint32) {
	shouldUpdate := c.NextIndexForPeer(peerId) != index ||
		c.MatchIndexForPeer(peerId) != index
	if shouldUpdate {
		newState := c.
			WithNextIndex(peerId, index).
			WithMatchIndex(peerId, index)
		c.swapState(&newState)
	}
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
	newRole := domain.Follower
	newState := c.
		WithRole(newRole).
		WithCurrentTerm(term).
		WithVotedFor("")
	c.swapState(&newState)
	c.shiftRole()
	c.resetTimeout()
}

func (c *ClusterNode) IncrementCandidateTerm() {
	if c.Role() == domain.Candidate {
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
	if c.Role() == domain.Follower {
		log.Info("UPGRADE TO CANDIDATE")
		newRole := domain.Candidate
		newState := c.
			WithRole(newRole).
			WithClusterLeader("")
		c.swapState(&newState)
		c.shiftRole()
		c.resetTimeout()
		return
	}
	log.Warn("CANNOT UPGRADE CANDIDATE FOR ", c.Role())
}

func (c *ClusterNode) UpgradeLeader() {
	if c.Role() == domain.Candidate {
		log.Infof("UPGRADE TO LEADER TERM %d\n", c.CurrentTerm())
		newRole := domain.Leader
		newState := c.
			WithInitializeLeader().
			WithClusterLeader(c.Id()).
			WithRole(newRole)
		c.swapState(&newState)
		// noOp will force logs commit and trigger
		// ApplyLogs in the entire cluster
		noOp := domain.LogEntry{
			Term:  c.CurrentTerm(),
			Type:  "ADMIN",
			Value: "NO_OP",
		}
		c.AppendLogs(c.LastLogIndex(), noOp)
		c.shiftRole()
		c.resetLeaderTicker()
		return
	}
	log.Warn("CANNOT UPGRADE LEADER FOR ", c.Role())
}

func (c *ClusterNode) ApplyLogs() {
	commitIndex := c.CommitIndex()
	lastApplied := c.LastApplied()

	for lastApplied < commitIndex {
		if lastApplied == 0 {
			c.initFsm()
		}
		// Increment last applied first
		// because lastApplied = 0 is not a valid logEntry
		lastApplied += 1
		if log, err := c.Log(lastApplied); err == nil {
			if log.Type != "COMMAND" {
				continue
			}
			res := c.evalFsm(log.Value, "COMMAND")
			if log.C != nil {
				log.C <- res
			}
		}
	}
	// Swap only once
	newState := c.WithLastApplied(lastApplied)
	c.swapState(&newState)
}

func (c *ClusterNode) ExecuteCommand(command string) chan domain.EvalResult {
	result := make(chan domain.EvalResult, 1)
	newEntry := domain.LogEntry{
		Value: command,
		Term:  c.CurrentTerm(),
		Type:  "COMMAND",
		C:     result,
	}

	go (func() {
		c.AppendLogs(c.LastLogIndex(), newEntry)
		c.synchronizeLogs()
	})()

	return result
}

func (c ClusterNode) ExecuteQuery(query string) chan domain.EvalResult {
	result := make(chan domain.EvalResult, 1)
	go (func() { result <- c.evalFsm(query, "QUERY") })()
	return result
}

func (c ClusterNode) evalFsm(entry string, entryType string) domain.EvalResult {
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

	var outb, errb bytes.Buffer
	cmd.Stdout = &outb
	cmd.Stderr = &errb

	if err := cmd.Run(); err != nil {
		return domain.EvalResult{
			Err: errors.New(errb.String()),
		}
	}

	return domain.EvalResult{
		Out: outb.Bytes(),
	}
}

func (c ClusterNode) initFsm() {
	cmd := exec.Command(c.fsmInit, c.Id())
	cmd.Run()
}
