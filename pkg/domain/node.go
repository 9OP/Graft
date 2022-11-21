package domain

import (
	"sync"
	"sync/atomic"
	"unsafe"

	log "github.com/sirupsen/logrus"
)

/*
Implementation note:
- nodeState is supposed to only expose pure function `Withers` that return a copy
  instead of mutating the instance.
- When a mutation needs to be saved in ClusterNode (for instance incrementing nodeState.CurrentTerm)
  You should call swapState, which atomically change the pointer of ClusterNode.nodeState to your
  new nodeState version.
- This means that every ClusterNode method that calls swapState, needs to accept a pointer receiver (n *Node)
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

func (s signals) commit() {
	s.Commit <- struct{}{}
}

func (s signals) shiftRole() {
	s.ShiftRole <- struct{}{}
}

func (s signals) saveState() {
	s.SaveState <- struct{}{}
}

func (s signals) synchronizeLogs() {
	s.SynchronizeLogs <- struct{}{}
}

func (s signals) resetLeaderTicker() {
	s.ResetLeaderTicker <- struct{}{}
}

func (s signals) resetTimeout() {
	s.ResetElectionTimer <- struct{}{}
}

type Node struct {
	*state
	signals
}

func NewNode(
	id string,
	peers Peers,
	persistent PersistentState,
) *Node {
	return &Node{
		state:   newState(id, peers, persistent),
		signals: newSignals(),
	}
}

func (n *Node) swapState(s interface{}) {
	var addr *unsafe.Pointer
	var new unsafe.Pointer

	switch newState := s.(type) {
	case *state:
		addr = (*unsafe.Pointer)(unsafe.Pointer(&n.state))
		new = (unsafe.Pointer(newState))
	default:
		return
	}

	atomic.SwapPointer(addr, new)
	n.saveState()
}

func (n Node) GetState() state {
	return *n.state
}

func (n *Node) DowngradeFollower(term uint32) {
	log.Infof("DOWNGRADE TO FOLLOWER TERM: %d\n", term)
	newState := n.
		WithRole(Follower).
		WithCurrentTerm(term).
		WithVotedFor("")
	n.swapState(&newState)
	n.shiftRole()
	n.resetTimeout()
}

func (n *Node) IncrementCandidateTerm() {
	if n.Role() == Candidate {
		term := n.CurrentTerm() + 1
		log.Debugf("INCREMENT CANDIDATE TERM %d\n", term)
		newState := n.
			WithCurrentTerm(term).
			WithVotedFor(n.Id())
		n.swapState(&newState)
		n.resetTimeout()
		return
	}
	log.Warn("CANNOT INCREMENT TERM FOR ", n.Role())
}

func (n *Node) UpgradeCandidate() {
	if n.Role() == Follower {
		log.Info("UPGRADE TO CANDIDATE")
		newState := n.
			WithRole(Candidate).
			WithLeader("")
		n.swapState(&newState)
		n.shiftRole()
		n.resetTimeout()
		return
	}
	log.Warn("CANNOT UPGRADE CANDIDATE FOR ", n.Role())
}

func (n *Node) UpgradeLeader() {
	if n.Role() == Candidate {
		log.Infof("UPGRADE TO LEADER TERM %d\n", n.CurrentTerm())
		newState := n.
			WithLeaderInitialization().
			WithLeader(n.id).
			WithRole(Leader)
		n.swapState(&newState)
		// noOp will force logs commit and trigger
		// ApplyLogs in the entire cluster
		noop := LogEntry{
			Term:  n.CurrentTerm(),
			Type:  "ADMIN",
			Value: "NO_OP",
		}
		n.AppendLogs(n.LastLogIndex(), noop)
		n.shiftRole()
		n.resetLeaderTicker()
		return
	}
	log.Warn("CANNOT UPGRADE LEADER FOR ", n.Role())
}

func (n Node) Heartbeat() {
	n.resetTimeout()
}

func (n Node) Broadcast(fn func(p Peer)) {
	var wg sync.WaitGroup
	for _, peer := range n.Peers() {
		wg.Add(1)
		go func(p Peer, w *sync.WaitGroup) {
			defer w.Done()
			fn(p)
		}(peer, &wg)
	}
	wg.Wait()
}

func (n *Node) DeleteLogsFrom(index uint32) {
	if newState, changed := n.WithDeleteLogsFrom(index); changed {
		log.Debug("DELETE LOGS", index)
		n.swapState(&newState)
	}
}

func (n *Node) AppendLogs(prevLogIndex uint32, entries ...LogEntry) {
	if newState, changed := n.WithAppendLogs(prevLogIndex, entries...); changed {
		log.Debug("APPEND LOGS", prevLogIndex)
		n.swapState(&newState)
	}
}

func (n *Node) SetLeader(leaderId string) {
	if n.Leader().Id != leaderId {
		log.Debug("SET LEADER", leaderId)
		newState := n.WithLeader(leaderId)
		n.swapState(&newState)
	}
}

func (n *Node) SetCommitIndex(index uint32) {
	if n.CommitIndex() != index {
		log.Debug("SET COMMIT INDEX", index)
		newState := n.WithCommitIndex(index)
		n.swapState(&newState)
		n.commit()
	}
}

func (n *Node) SetNextMatchIndex(peerId string, index uint32) {
	if n.NextIndexForPeer(peerId) != index || n.MatchIndexForPeer(peerId) != index {
		newState := n.
			WithNextIndex(peerId, index).
			WithMatchIndex(peerId, index)
		n.swapState(&newState)
	}
}

func (n *Node) DecrementNextIndex(peerId string) {
	if n.NextIndexForPeer(peerId) > 0 {
		newState := n.WithDecrementNextIndex(peerId)
		n.swapState(&newState)
	}
}

func (n *Node) GrantVote(peerId string) {
	if n.VotedFor() != peerId {
		newState := n.WithVotedFor(peerId)
		n.swapState(&newState)
	}
}

func (n *Node) ApplyLogs() {
	commitIndex := n.CommitIndex()
	lastApplied := n.LastApplied()

	for lastApplied < commitIndex {
		// if lastApplied == 0 {
		// 	c.initFsm()
		// }
		// Increment last applied first
		// because lastApplied = 0 is not a valid logEntry
		lastApplied += 1
		if log, err := n.Log(lastApplied); err == nil {
			if log.Type != "COMMAND" {
				continue
			}
			// res := c.evalFsm(log.Value, "COMMAND")
			// if log.C != nil {
			// 	log.C <- res
			// }
		}
	}

	// Swap only once
	if lastApplied > n.LastApplied() {
		newState := n.WithLastApplied(lastApplied)
		n.swapState(&newState)
	}
}
