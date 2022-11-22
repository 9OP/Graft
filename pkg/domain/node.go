package domain

import (
	"encoding/json"
	"errors"
	"sync"
	"sync/atomic"
	"unsafe"

	"graft/pkg/utils"

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
		withRole(Follower).
		withCurrentTerm(term).
		withVotedFor("")
	n.swapState(&newState)
	n.shiftRole()
	n.resetTimeout()
}

func (n *Node) IncrementCandidateTerm() {
	if n.Role() == Candidate {
		term := n.currentTerm + 1
		log.Debugf("INCREMENT CANDIDATE TERM %d\n", term)
		newState := n.
			withCurrentTerm(term).
			withVotedFor(n.id)
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
			withRole(Candidate).
			withLeader("")
		n.swapState(&newState)
		n.shiftRole()
		n.resetTimeout()
		return
	}
	log.Warn("CANNOT UPGRADE CANDIDATE FOR ", n.Role())
}

func (n *Node) UpgradeLeader() {
	if n.Role() == Candidate {
		log.Infof("UPGRADE TO LEADER TERM %d\n", n.currentTerm)
		newState := n.
			withLeaderInitialization().
			withLeader(n.id).
			withRole(Leader)
		n.swapState(&newState)
		n.AppendLogs(
			n.lastLogIndex(),
			LogEntry{
				Term: n.currentTerm,
				Type: LogNoop,
			},
		)
		n.shiftRole()
		n.resetLeaderTicker()
		return
	}
	log.Warn("CANNOT UPGRADE LEADER FOR ", n.Role())
}

func (n Node) Heartbeat() {
	n.resetTimeout()
}

type BroadcastType uint8

const (
	BroadcastAll BroadcastType = iota
	BroadcastActive
	BroadcastInactive
)

func (n Node) Broadcast(fn func(p Peer), broadcastType BroadcastType) {
	var wg sync.WaitGroup
	var peers Peers

	switch broadcastType {
	case BroadcastAll:
		peers = n.peers
	case BroadcastActive:
		peers = n.activePeers()
	default:
		peers = Peers{}
	}

	for _, peer := range peers {
		wg.Add(1)
		go func(p Peer, w *sync.WaitGroup) {
			defer w.Done()
			fn(p)
		}(peer, &wg)
	}
	wg.Wait()
}

func (n *Node) DeleteLogsFrom(index uint32) {
	if newState, changed := n.withDeleteLogsFrom(index); changed {
		log.Debug("DELETE LOGS", index)
		n.swapState(&newState)
	}
}

func (n *Node) AppendLogs(prevLogIndex uint32, entries ...LogEntry) {
	if newState, changed := n.withAppendLogs(prevLogIndex, entries...); changed {
		log.Debug("APPEND LOGS", prevLogIndex)
		n.swapState(&newState)
	}
}

func (n *Node) SetLeader(leaderId string) {
	if n.Leader().Id != leaderId {
		log.Debug("SET LEADER", leaderId)
		newState := n.withLeader(leaderId)
		n.swapState(&newState)
	}
}

func (n *Node) setCommitIndex(index uint32) {
	if n.commitIndex != index {
		log.Debug("SET COMMIT INDEX", index)
		newState := n.withCommitIndex(index)
		n.swapState(&newState)
		n.commit()
	}
}

func (n *Node) UpdateLeaderCommitIndex(leaderCommitIndex uint32) {
	if leaderCommitIndex > n.commitIndex {
		newIndex := utils.Min(n.lastLogIndex(), leaderCommitIndex)
		n.setCommitIndex(newIndex)
	}
}

func (n *Node) UpdateNewCommitIndex() {
	newCommitIndex := n.computeNewCommitIndex()
	n.setCommitIndex(newCommitIndex)
}

func (n *Node) SetNextMatchIndex(peerId string, index uint32) {
	if n.nextIndexForPeer(peerId) != index || n.matchIndexForPeer(peerId) != index {
		newState := n.
			withNextIndex(peerId, index).
			withMatchIndex(peerId, index)
		n.swapState(&newState)
	}
}

func (n *Node) DecrementNextIndex(peerId string) {
	if n.nextIndexForPeer(peerId) > 0 {
		newState := n.withDecrementNextIndex(peerId)
		n.swapState(&newState)
	}
}

func (n *Node) GrantVote(peerId string) {
	if n.votedFor != peerId {
		newState := n.withVotedFor(peerId)
		n.swapState(&newState)
	}
}

func (n *Node) ApplyLogs() {
	commitIndex := n.commitIndex
	lastApplied := n.lastApplied

	for lastApplied < commitIndex {
		// if lastApplied == 0 {
		// 	c.initFsm()
		// }

		// Increment last applied first
		// because lastApplied = 0 is not a valid logEntry
		lastApplied += 1
		if log, err := n.Log(lastApplied); err == nil {
			switch log.Type {
			case LogCommand:
				// do command
				// res := c.evalFsm(log.Value, "COMMAND")
				// if log.C != nil {
				// 	log.C <- res
				// }
			case LogConfiguration:
				var config ConfigurationUpdate
				if err := json.Unmarshal(log.Data, &config); err == nil {
					if res := n.configurationUpdate(config); log.C != nil {
						log.C <- res
					}
				}
			case LogNoop:
				continue
			default:
				continue
			}
		}
	}

	// Swap only once
	if lastApplied > n.lastApplied {
		newState := n.withLastApplied(lastApplied)
		n.swapState(&newState)
	}
}

func (n Node) ExecuteCommand(cmd ApiCommand) chan EvalResult {
	result := make(chan EvalResult, 1)
	newEntry := LogEntry{
		Index: uint64(n.lastLogIndex()),
		Term:  n.currentTerm,
		Type:  cmd.Type,
		Data:  cmd.Data,
		C:     result,
	}

	go (func() {
		n.AppendLogs(n.lastLogIndex(), newEntry)
		n.synchronizeLogs()
	})()

	return result
}

// func (c ClusterNode) ExecuteQuery(query string) chan domain.EvalResult {
// 	result := make(chan domain.EvalResult, 1)
// 	go (func() { result <- c.evalFsm(query, "QUERY") })()
// 	return result
// }

var (
	errPeerAlreadyExists = errors.New("peer already exists")
	errPeerDoesNotExist  = errors.New("peer does not exist")
)

func (n *Node) configurationUpdate(config ConfigurationUpdate) (res EvalResult) {
	res = EvalResult{
		Out: nil,
		Err: nil,
	}

	var peers Peers
	peer := config.Peer

	switch config.Type {
	case ConfAddPeer:
		if p, ok := n.peers[peer.Id]; ok && p.Active {
			res.Err = errPeerAlreadyExists
			return
		}
		log.Info("configuration update: add peer", peer.Id)
		peer.Active = false
		peers = n.peers.addPeer(peer)

	case ConfActivatePeer:
		if _, ok := n.peers[peer.Id]; !ok {
			res.Err = errPeerDoesNotExist
			return
		}
		log.Info("configuration update: activate peer", peer.Id)
		peers = n.peers.activatePeer(peer.Id)

	case ConfDeactivatePeer:
		if _, ok := n.peers[peer.Id]; !ok {
			res.Err = errPeerDoesNotExist
			return
		}
		log.Info("configuration update: deactivate peer", peer.Id)
		peers = n.peers.deactivatePeer(peer.Id)

	case ConfRemovePeer:
		if _, ok := n.peers[peer.Id]; !ok {
			res.Err = errPeerDoesNotExist
			return
		}
		log.Info("configuration update: remove peer", peer.Id)
		peers = n.peers.removePeer(peer.Id)

	default:
		return
	}

	newState := n.withPeers(peers)
	n.swapState(&newState)
	return
}
