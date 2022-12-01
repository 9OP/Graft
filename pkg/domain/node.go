package domain

import (
	"encoding/json"
	"errors"
	"net/netip"
	"sync"
	"sync/atomic"
	"unsafe"

	"graft/pkg/utils"
	"graft/pkg/utils/log"
)

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
		Commit:             make(chan struct{}, 10),
		ShiftRole:          make(chan struct{}, 10),
		SaveState:          make(chan struct{}, 10),
		SynchronizeLogs:    make(chan struct{}, 10),
		ResetLeaderTicker:  make(chan struct{}, 10),
		ResetElectionTimer: make(chan struct{}, 10),
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

type NodeConfig struct {
	Id              string
	Host            netip.AddrPort
	ElectionTimeout int
	LeaderHeartbeat int
}

type Node struct {
	*state
	signals
	config NodeConfig
	fsm    fsm
	exit   bool
	apply  sync.Mutex
}

func NewNode(
	config NodeConfig,
	fsm string,
	peers Peers,
	persistent PersistentState,
) *Node {
	return &Node{
		state:   newState(config.Id, peers, persistent),
		signals: newSignals(),
		config:  config,
		exit:    false,
		fsm:     NewFsm(fsm),
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

func (n *Node) GetState() state {
	return *n.state
}

func (n *Node) GetClusterConfiguration() ClusterConfiguration {
	return ClusterConfiguration{
		Peers:           utils.CopyMap(n.peers),
		LeaderId:        n.leaderId,
		Fsm:             n.fsm.path,
		ElectionTimeout: n.config.ElectionTimeout,
		LeaderHeartbeat: n.config.LeaderHeartbeat,
	}
}

func (n *Node) Shutdown() {
	n.exit = true
	n.DowngradeFollower(n.currentTerm)
}

func (n *Node) IsShuttingDown() bool {
	return n.exit
}

func (n *Node) DowngradeFollower(term uint32) {
	log.Infof("downgrade Follower %v", term)
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
		log.Debugf("increment term %v", term)
		newState := n.
			withCurrentTerm(term).
			withVotedFor(n.id)
		n.swapState(&newState)
		n.resetTimeout()
		return
	}
	log.Warnf("cannot increment: %v", n.Role())
}

func (n *Node) UpgradeCandidate() {
	if n.exit {
		return
	}

	if n.Role() == Follower {
		log.Infof("upgrade Candidate")
		newState := n.
			withRole(Candidate).
			withLeader("")
		n.swapState(&newState)
		n.shiftRole()
		n.resetTimeout()
		return
	}
	log.Warnf("cannot upgrade candidate: %v", n.Role())
}

func (n *Node) UpgradeLeader() {
	if n.Role() == Candidate {
		log.Infof("upgrade Leader %v", n.currentTerm)
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
	log.Warnf("cannot upgrade leader: %v", n.Role())
}

func (n *Node) Heartbeat() {
	n.resetTimeout()
}

type BroadcastType uint8

const (
	BroadcastAll BroadcastType = iota
	BroadcastActive
	BroadcastInactive
)

func (n *Node) Broadcast(fn func(p Peer), broadcastType BroadcastType) {
	var wg sync.WaitGroup
	var peers Peers

	switch broadcastType {
	case BroadcastAll:
		peers = utils.CopyMap(n.peers)
	case BroadcastActive:
		peers = n.activePeers()
	default:
		peers = Peers{}
	}

	// Prevent sending requests to itself
	delete(peers, n.id)

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
		log.Warnf("delete logs from index %v", index)
		n.swapState(&newState)
	}
}

func (n *Node) AppendLogs(prevLogIndex uint32, entries ...LogEntry) {
	if newState, changed := n.withAppendLogs(prevLogIndex, entries...); changed {
		log.Debugf("append logs from index %v", prevLogIndex)
		n.swapState(&newState)
	}
}

func (n *Node) SetLeader(leaderId string) {
	if n.Leader().Id != leaderId {
		log.Infof("follow leader: %v", leaderId)
		newState := n.withLeader(leaderId)
		n.swapState(&newState)
	}
}

func (n *Node) setCommitIndex(index uint32) {
	if index > n.lastLogIndex() {
		log.Errorf("cannot set commit index (%v) > last log index (%v)", index, n.lastLogIndex())
		return
	}

	if n.commitIndex != index {
		log.Debugf("commit index: %v", index)
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

func (n *Node) ComputeNewCommitIndex() {
	if n.Role() == Leader {
		newCommitIndex := n.computeNewCommitIndex()
		n.setCommitIndex(newCommitIndex)
		return
	}
	log.Warnf("cannot compute new commit index for: %v", n.Role())
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
	log.Debugf("apply logs to commitIndex %v", n.commitIndex)

	if n.lastApplied >= n.commitIndex || n.exit {
		log.Warnf("not applying logs")
		return
	}

	/*
		Critical section:
		Must not run ApplyLogs concurrently to prevent applying
		the same log multiple time to the FSM.
		ApplyLogs should be run by only one goroutine at any given time.
	*/
	n.apply.Lock()
	defer n.apply.Unlock()

	// Must set lastApplied and commitIndex
	// After calling Lock()
	lastApplied := n.lastApplied
	commitIndex := n.commitIndex

	if lastApplied == 0 {
		input := EvalInput{
			Id:       n.id,
			EvalType: INIT,
		}
		n.fsm.eval(input)
	}

	for lastApplied < commitIndex {
		lastApplied += 1
		n.applyLog(lastApplied)
	}

	// Swap only once
	newState := n.withLastApplied(lastApplied)
	n.swapState(&newState)
}

func (n *Node) applyLog(index uint32) {
	if lg, err := n.Log(index); err == nil {
		switch lg.Type {
		case LogCommand:
			input := EvalInput{
				Id:       n.id,
				Data:     string(lg.Data),
				EvalType: COMMAND,
			}
			res := n.fsm.eval(input)
			log.Debugf("apply\n%s", string(lg.Data))
			if lg.C != nil {
				lg.C <- res
			}

		case LogConfiguration:
			var config ConfigurationUpdate
			if err := json.Unmarshal(lg.Data, &config); err == nil {
				if res := n.configurationUpdate(config); lg.C != nil {
					lg.C <- res
				}
			}

		case LogNoop:
			return

		default:
			log.Errorf("log type unknown: %v", lg.Type)
		}
	}
}

func (n *Node) Execute(cmd ExecuteInput) chan ExecuteOutput {
	result := make(chan ExecuteOutput, 1)

	switch cmd.Type {
	case LogCommand, LogConfiguration:
		newEntry := LogEntry{
			Index: uint64(n.lastLogIndex()),
			Term:  n.currentTerm,
			Type:  cmd.Type,
			Data:  cmd.Data,
			C:     result,
		}
		dispatch := func() {
			n.AppendLogs(n.lastLogIndex(), newEntry)
			n.synchronizeLogs()
		}

		go dispatch()

	case LogQuery:
		dispatch := func() {
			input := EvalInput{
				Id:       n.id,
				Data:     string(cmd.Data),
				EvalType: QUERY,
			}
			result <- n.fsm.eval(input)
		}
		go dispatch()

	default:
		log.Errorf("log type unknown: %v", cmd.Type)
		result <- ExecuteOutput{}
	}

	return result
}

var (
	errPeerAlreadyExists = errors.New("peer already exists")
	errPeerDoesNotExist  = errors.New("peer does not exist")
)

func (n *Node) configurationUpdate(config ConfigurationUpdate) (res ExecuteOutput) {
	res = ExecuteOutput{
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
		log.Infof("configuration: add %v", peer.Id)
		peer.Active = false
		peers = n.peers.addPeer(peer)

	case ConfActivatePeer:
		if _, ok := n.peers[peer.Id]; !ok {
			res.Err = errPeerDoesNotExist
			return
		}
		log.Infof("configuration: activate %v", peer.Id)
		peers = n.peers.activatePeer(peer.Id)

	case ConfDeactivatePeer:
		if _, ok := n.peers[peer.Id]; !ok {
			res.Err = errPeerDoesNotExist
			return
		}
		log.Infof("configuration: deactivate %v", peer.Id)
		peers = n.peers.deactivatePeer(peer.Id)

	case ConfRemovePeer:
		if _, ok := n.peers[peer.Id]; !ok {
			res.Err = errPeerDoesNotExist
			return
		}
		log.Infof("configuration: remove %v", peer.Id)
		peers = n.peers.removePeer(peer.Id)

	default:
		return
	}

	newState := n.withPeers(peers)
	n.swapState(&newState)
	return
}
