package domain

import (
	"errors"
	"math"

	"graft/pkg/utils"
)

type state struct {
	id       string
	leaderId string
	peers    Peers
	role     Role

	commitIndex uint32
	lastApplied uint32
	nextIndex   peerIndex
	matchIndex  peerIndex

	currentTerm uint32
	votedFor    string
	machineLogs []LogEntry
}

type peerIndex map[string]uint32

type PersistentState struct {
	CurrentTerm uint32     `json:"current_term"`
	VotedFor    string     `json:"voted_for"`
	MachineLogs []LogEntry `json:"machine_logs"`
}

func (n Node) ToPersistent() PersistentState {
	return PersistentState{
		CurrentTerm: n.currentTerm,
		VotedFor:    n.votedFor,
		MachineLogs: n.logs(),
	}
}

var DEFAULT_PERSISTENT_STATE = PersistentState{
	CurrentTerm: 0,
	VotedFor:    "",
	MachineLogs: []LogEntry{},
}

func newState(id string, peers Peers, persistent PersistentState) *state {
	return &state{
		id:    id,
		peers: peers,
		role:  Follower,

		commitIndex: 0,
		lastApplied: 0,
		nextIndex:   map[string]uint32{},
		matchIndex:  map[string]uint32{},

		currentTerm: persistent.CurrentTerm,
		votedFor:    persistent.VotedFor,
		machineLogs: persistent.MachineLogs,
	}
}

func (s state) Leader() Peer {
	return s.peers[s.leaderId]
}

func (s state) HasLeader() bool {
	_, ok := s.peers[s.leaderId]
	return ok
}

func (s state) IsLeader() bool {
	return s.id == s.leaderId
}

func (s state) withLeader(leaderId string) state {
	s.leaderId = leaderId
	return s
}

func (s state) withLeaderInitialization() state {
	defaultNextIndex := s.lastLogIndex()
	nextIndex := make(map[string]uint32, len(s.peers))
	matchIndex := make(map[string]uint32, len(s.peers))
	for _, peer := range s.peers {
		nextIndex[peer.Id] = defaultNextIndex
		matchIndex[peer.Id] = 0
	}
	s.nextIndex = nextIndex
	s.matchIndex = matchIndex
	return s
}

// Quorum is defined as the absolute majority of nodes:
// (N + 1) / 2, where N is the number of nodes in the cluster
// which can recover from up to (N - 1) / 2 failures
func (s state) Quorum() int {
	numberNodes := float64(len(s.peers) + 1) // add self
	return int(math.Ceil((numberNodes + 1) / 2.0))
}

func (s state) Role() Role {
	return s.role
}

func (s state) withRole(role Role) state {
	s.role = role
	return s
}

func (s state) withCommitIndex(index uint32) state {
	s.commitIndex = index
	return s
}

// Compute new commitIndex N such that:
// - N > commitIndex,
// - a majority of matchIndex[i] ≥ N
// - log[N].term == currentTerm
func (s state) computeNewCommitIndex() uint32 {
	quorum := s.Quorum()
	lastLogIndex := s.lastLogIndex() // Upper value of N
	commitIndex := s.commitIndex     // Lower value of N

	for N := lastLogIndex; N > commitIndex; N-- {
		count := 1 // count self
		for _, mIndex := range s.matchIndex {
			if mIndex >= N {
				count += 1
			}
		}
		log, _ := s.Log(N)
		if count >= quorum && log.Term == s.currentTerm {
			return N
		}
	}

	return commitIndex
}

func (s state) withLastApplied(lastApplied uint32) state {
	s.lastApplied = lastApplied
	return s
}

func (s state) nextIndexForPeer(peerId string) uint32 {
	if idx, ok := s.nextIndex[peerId]; ok {
		return idx
	}
	return 0
}

func (s state) withNextIndex(peerId string, index uint32) state {
	if idx, ok := s.nextIndex[peerId]; (idx == index && ok) || !ok {
		return s
	}
	nextIndex := utils.CopyMap(s.nextIndex)
	nextIndex[peerId] = index
	s.nextIndex = nextIndex
	return s
}

func (s state) withDecrementNextIndex(peerId string) state {
	if idx, ok := s.nextIndex[peerId]; ok && idx > 0 {
		return s.withNextIndex(peerId, idx-1)
	}
	return s
}

func (s state) matchIndexForPeer(peerId string) uint32 {
	if idx, ok := s.matchIndex[peerId]; ok {
		return idx
	}
	return 0
}

func (s state) withMatchIndex(peerId string, index uint32) state {
	if idx, ok := s.matchIndex[peerId]; (idx == index && ok) || !ok {
		return s
	}
	matchIndex := utils.CopyMap(s.matchIndex)
	matchIndex[peerId] = index
	s.matchIndex = matchIndex
	return s
}

func (s state) CurrentTerm() uint32 {
	return s.currentTerm
}

func (s state) withCurrentTerm(term uint32) state {
	s.currentTerm = term
	return s
}

func (s state) withVotedFor(vote string) state {
	s.votedFor = vote
	return s
}

func (s state) CanGrantVote(peerId string) bool {
	// Unknown peer id
	if _, ok := s.peers[peerId]; !ok {
		return false
	}
	voteAvailable := s.votedFor == "" || s.votedFor == peerId
	return voteAvailable
}

var errIndexOutOfRange = errors.New("index out of range")

func (s state) Log(index uint32) (LogEntry, error) {
	if index == 0 {
		return LogEntry{}, nil
	}
	if index <= s.lastLogIndex() {
		return s.machineLogs[index-1], nil
	}
	return LogEntry{}, errIndexOutOfRange
}

func (s state) logs() []LogEntry {
	machineLogs := make([]LogEntry, len(s.machineLogs))
	copy(machineLogs, s.machineLogs)
	return machineLogs
}

func (s state) lastLogIndex() uint32 {
	return uint32(len(s.machineLogs))
}

func (s state) lastLog() LogEntry {
	lastLog, _ := s.Log(s.lastLogIndex())
	return lastLog
}

func (s state) logsFrom(index uint32) []LogEntry {
	logs := s.logs()
	if index == 0 {
		return logs
	}
	if index <= s.lastLogIndex() {
		return logs[index-1:]
	}
	return []LogEntry{}
}

// IsUpToDate from the caller point of view
//
// Note:
// Raft determines which of two logs is more up-to-date
// by comparing the index and term of the last entries in the
// logs. If the logs have last entries with different terms, then
// the log with the later term is more up-to-date. If the logs
// end with the same term, then whichever log is longer is
// more up-to-date.
func (s state) IsLogUpToDate(lastLogIndex uint32, lastLogTerm uint32) bool {
	log := s.lastLog()

	if log.Term == lastLogTerm {
		return lastLogIndex >= s.lastLogIndex()
	}

	return lastLogTerm >= log.Term
}

func (s state) withDeleteLogsFrom(index uint32) (state, bool) {
	if index == 0 {
		s.machineLogs = []LogEntry{}
		return s, true
	}
	logs := s.logs()
	if index <= s.lastLogIndex() && index >= 1 {
		s.machineLogs = logs[:index-1]
		return s, true
	}
	s.machineLogs = logs
	return s, false
}

func (s state) withAppendLogs(prevLogIndex uint32, entries ...LogEntry) (state, bool) {
	// Should append only new entries
	changed := false
	if len(entries) == 0 {
		return s, changed
	}

	/*
		Prevent appending logs twice
		Example:
		logs  -> [log1, log2, log3, log4, log5]
		term  ->   t1    t1    t2    t4    t4
		index ->   i1    i2    i3    i4    i5

		arguments -> [log3, log4, log5, log6], i2

		result -> [log1, log2, log3, log4, log5, log6]

		---
		We should increment prevLogIndex up to lastLogIndex.
		While incrementing a pointer in the entries slice.
		and then only copy the remaining "new" logs.

		prevLogIndex is necessay as it gives the offset of the
		entries argument, relative to s.Logs
	*/

	lastLogIndex := s.lastLogIndex()
	lenEntries := uint32(len(entries))

	// Find index of newLogs
	var newLogsFromIndex uint32
	if lastLogIndex >= prevLogIndex {
		newLogsFromIndex = utils.Min(lastLogIndex-prevLogIndex, lenEntries)
	} else {
		newLogsFromIndex = lenEntries
	}
	changed = len(entries[newLogsFromIndex:]) > 0

	if !changed {
		return s, changed
	}

	// Copy existings logs
	logs := make([]LogEntry, lastLogIndex, lenEntries+lastLogIndex)

	// Append new logs
	logs = append(logs, entries[newLogsFromIndex:]...)
	copy(logs, s.logs())

	s.machineLogs = logs
	return s, changed
}

func (s state) AppendEntriesInput(peerId string) AppendEntriesInput {
	matchIndex := s.matchIndexForPeer(peerId)
	nextIndex := s.nextIndexForPeer(peerId)

	var prevLogTerm uint32
	var prevLogIndex uint32
	var entries []LogEntry

	// When peer has commitIndex matching
	// leader lastLogIndex, there is no need to send
	// new entries.
	if matchIndex == s.lastLogIndex() {
		entries = []LogEntry{}
		prevLogIndex = s.lastLogIndex()
		prevLogTerm = s.lastLog().Term
	} else {
		entries = s.logsFrom(nextIndex + 1)
		prevLogIndex = nextIndex

		// When prevLogIndex >= lastLogIndex, prevLogTerm
		// is at least >= CurrentTerm logically
		// It also prevent to match the receiver condition:
		// prevLogTerm == log(prevLogIndex).Term
		prevLog, err := s.Log(prevLogIndex)
		if err == errIndexOutOfRange {
			prevLogTerm = s.currentTerm
		} else {
			prevLogTerm = prevLog.Term
		}
	}

	return AppendEntriesInput{
		LeaderId:     s.id,
		Term:         s.currentTerm,
		LeaderCommit: s.commitIndex,
		PrevLogIndex: prevLogIndex,
		PrevLogTerm:  prevLogTerm,
		Entries:      entries,
	}
}

func (s state) RequestVoteInput() RequestVoteInput {
	return RequestVoteInput{
		CandidateId:  s.id,
		Term:         s.currentTerm,
		LastLogIndex: s.lastLogIndex(),
		LastLogTerm:  s.lastLog().Term,
	}
}
