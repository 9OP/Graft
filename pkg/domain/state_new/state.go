package statenew

import (
	"errors"
	"math"

	"graft/pkg/domain"
	"graft/pkg/utils"
)

type state struct {
	id       string
	leaderId string
	peers    domain.Peers
	role     domain.Role

	commitIndex uint32
	lastApplied uint32
	nextIndex   peerIndex
	matchIndex  peerIndex

	currentTerm uint32
	votedFor    string
	machineLogs []domain.LogEntry
}

type peerIndex map[string]uint32

type PersistentState struct {
	CurrentTerm uint32            `json:"current_term"`
	VotedFor    string            `json:"voted_for"`
	MachineLogs []domain.LogEntry `json:"machine_logs"`
}

func newState(id string, peers domain.Peers, persistent PersistentState) *state {
	return &state{
		id:    id,
		peers: peers,
		role:  domain.Follower,

		commitIndex: 0,
		lastApplied: 0,
		nextIndex:   map[string]uint32{},
		matchIndex:  map[string]uint32{},

		currentTerm: persistent.CurrentTerm,
		votedFor:    persistent.VotedFor,
		machineLogs: persistent.MachineLogs,
	}
}

func (s state) Id() string {
	return s.id
}

func (s state) Leader() domain.Peer {
	return s.peers[s.leaderId]
}

func (s state) HasLeader() bool {
	_, ok := s.peers[s.leaderId]
	return ok
}

func (s state) IsLeader() bool {
	return s.id == s.leaderId
}

func (s state) WithLeader(leaderId string) state {
	s.leaderId = leaderId
	return s
}

func (s state) WithLeaderInitialization() state {
	defaultNextIndex := s.LastLogIndex()
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

func (s state) Peers() domain.Peers {
	return utils.CopyMap(s.peers)
}

// Quorum is defined as the absolute majority of nodes:
// (N + 1) / 2, where N is the number of nodes in the cluster
// which can recover from up to (N - 1) / 2 failures
func (s state) Quorum() int {
	numberNodes := float64(len(s.peers) + 1) // add self
	return int(math.Ceil((numberNodes + 1) / 2.0))
}

func (s state) Role() domain.Role {
	return s.role
}

func (s state) WithRole(role domain.Role) state {
	s.role = role
	return s
}

func (s state) CommitIndex() uint32 {
	return s.commitIndex
}

func (s state) WithCommitIndex(index uint32) state {
	s.commitIndex = index
	return s
}

// Compute new commitIndex N such that:
// - N > commitIndex,
// - a majority of matchIndex[i] ≥ N
// - log[N].term == currentTerm
func (s state) ComputeNewCommitIndex() uint32 {
	quorum := s.Quorum()
	lastLogIndex := s.LastLogIndex() // Upper value of N
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

func (s state) LastApplied() uint32 {
	return s.lastApplied
}

func (s state) WithLastApplied(lastApplied uint32) state {
	s.lastApplied = lastApplied
	return s
}

func (s state) NextIndexForPeer(peerId string) uint32 {
	if idx, ok := s.nextIndex[peerId]; ok {
		return idx
	}
	return 0
}

func (s state) WithNextIndex(peerId string, index uint32) state {
	if idx, ok := s.nextIndex[peerId]; (idx == index && ok) || !ok {
		return s
	}
	nextIndex := utils.CopyMap(s.nextIndex)
	nextIndex[peerId] = index
	s.nextIndex = nextIndex
	return s
}

func (s state) WithDecrementNextIndex(peerId string) state {
	if idx, ok := s.nextIndex[peerId]; ok && idx > 0 {
		return s.WithNextIndex(peerId, idx-1)
	}
	return s
}

func (s state) MatchIndexForPeer(peerId string) uint32 {
	if idx, ok := s.matchIndex[peerId]; ok {
		return idx
	}
	return 0
}

func (s state) WithMatchIndex(peerId string, index uint32) state {
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

func (s state) WithCurrentTerm(term uint32) state {
	s.currentTerm = term
	return s
}

func (s state) VotedFor() string {
	return s.votedFor
}

func (s state) WithVotedFor(vote string) state {
	s.votedFor = vote
	return s
}

func (s state) CanGrantVote(peerId string) bool {
	// Unknown peer id
	if _, ok := s.peers[peerId]; !ok {
		return false
	}
	votedFor := s.VotedFor()
	voteAvailable := votedFor == "" || votedFor == peerId
	return voteAvailable
}

var errIndexOutOfRange = errors.New("index out of range")

func (s state) Log(index uint32) (domain.LogEntry, error) {
	if index == 0 {
		return domain.LogEntry{}, nil
	}
	if index <= s.LastLogIndex() {
		return s.machineLogs[index-1], nil
	}
	return domain.LogEntry{}, errIndexOutOfRange
}

func (s state) Logs() []domain.LogEntry {
	machineLogs := make([]domain.LogEntry, len(s.machineLogs))
	copy(machineLogs, s.machineLogs)
	return machineLogs
}

func (s state) LastLogIndex() uint32 {
	return uint32(len(s.machineLogs))
}

func (s state) LastLog() domain.LogEntry {
	lastLog, _ := s.Log(s.LastLogIndex())
	return lastLog
}

func (s state) LogsFrom(index uint32) []domain.LogEntry {
	logs := s.Logs()
	if index == 0 {
		return logs
	}
	if index <= s.LastLogIndex() {
		return logs[index-1:]
	}
	return []domain.LogEntry{}
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
	log := s.LastLog()

	if log.Term == lastLogTerm {
		return lastLogIndex >= s.LastLogIndex()
	}

	return lastLogTerm >= log.Term
}

func (s state) WithDeleteLogsFrom(index uint32) (state, bool) {
	if index == 0 {
		s.machineLogs = []domain.LogEntry{}
		return s, true
	}
	logs := s.Logs()
	if index <= s.LastLogIndex() && index >= 1 {
		s.machineLogs = logs[:index-1]
		return s, true
	}
	s.machineLogs = logs
	return s, false
}

func (s state) WithAppendLogs(prevLogIndex uint32, entries ...domain.LogEntry) (state, bool) {
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

	lastLogIndex := s.LastLogIndex()
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
	logs := make([]domain.LogEntry, lastLogIndex, lenEntries+lastLogIndex)

	// Append new logs
	logs = append(logs, entries[newLogsFromIndex:]...)
	copy(logs, s.Logs())

	s.machineLogs = logs
	return s, changed
}

func (s state) AppendEntriesInput(peerId string) domain.AppendEntriesInput {
	matchIndex := s.MatchIndexForPeer(peerId)
	nextIndex := s.NextIndexForPeer(peerId)

	var prevLogTerm uint32
	var prevLogIndex uint32
	var entries []domain.LogEntry

	// When peer has commitIndex matching
	// leader lastLogIndex, there is no need to send
	// new entries.
	if matchIndex == s.LastLogIndex() {
		entries = []domain.LogEntry{}
		prevLogIndex = s.LastLogIndex()
		prevLogTerm = s.LastLog().Term
	} else {
		entries = s.LogsFrom(nextIndex + 1)
		prevLogIndex = nextIndex

		// When prevLogIndex >= lastLogIndex, prevLogTerm
		// is at least >= CurrentTerm logically
		// It also prevent to match the receiver condition:
		// prevLogTerm == log(prevLogIndex).Term
		prevLog, err := s.Log(prevLogIndex)
		if err == errIndexOutOfRange {
			prevLogTerm = s.CurrentTerm()
		} else {
			prevLogTerm = prevLog.Term
		}
	}

	return domain.AppendEntriesInput{
		LeaderId:     s.Id(),
		Term:         s.CurrentTerm(),
		LeaderCommit: s.CommitIndex(),
		PrevLogIndex: prevLogIndex,
		PrevLogTerm:  prevLogTerm,
		Entries:      entries,
	}
}

func (s state) RequestVoteInput() domain.RequestVoteInput {
	return domain.RequestVoteInput{
		CandidateId:  s.Id(),
		Term:         s.CurrentTerm(),
		LastLogIndex: s.LastLogIndex(),
		LastLogTerm:  s.LastLog().Term,
	}
}
