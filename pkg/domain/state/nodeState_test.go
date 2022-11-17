package state

import (
	"fmt"
	"reflect"
	"testing"

	"graft/pkg/domain"
)

func TestHasLeader(t *testing.T) {
	peer := domain.Peer{}
	res := []struct {
		// state
		leaderId string
		peers    domain.Peers
		// output
		hasLeader bool
	}{
		{"leaderId", domain.Peers{}, false},
		{"leaderId", domain.Peers{"leaderId": peer}, true},
		{"leaderId", domain.Peers{"id": peer}, false},
		{"", domain.Peers{"leaderId": peer}, false},
	}
	for _, tt := range res {
		t.Run(tt.leaderId, func(t *testing.T) {
			state := nodeState{leaderId: tt.leaderId, peers: tt.peers}
			ok := state.HasLeader()
			if ok != tt.hasLeader {
				t.Errorf("HasLeader got %v, want %v", ok, tt.hasLeader)
			}
		})
	}
}

func TestWithInitializeLeader(t *testing.T) {
}

func TestQuorum(t *testing.T) {
	p := domain.Peer{}
	res := []struct {
		// state
		peers domain.Peers
		// output
		quorum int
	}{
		{domain.Peers{"peer1": p, "peer2": p}, 2},                         // tolerate 1 failure
		{domain.Peers{"peer1": p, "peer2": p, "peer3": p}, 3},             // tolerate 1 failure
		{domain.Peers{"peer1": p, "peer2": p, "peer3": p, "peer4": p}, 3}, // tolerate 2 failures
		{domain.Peers{"peer1": p}, 2},                                     // tolerate 0 failure
		{domain.Peers{}, 1},                                               // tolerate 0 failure
	}
	for _, tt := range res {
		t.Run(fmt.Sprintf("%v", tt.peers), func(t *testing.T) {
			state := nodeState{peers: tt.peers}
			quorum := state.Quorum()
			if quorum != tt.quorum {
				t.Errorf("Quorum got %v, want %v", quorum, tt.quorum)
			}
		})
	}
}

func TestIsLogUpToDate(t *testing.T) {
	res := []struct {
		// state
		logs []domain.LogEntry
		// input
		lastLogIndex uint32
		lastLogTerm  uint32
		// output
		upToDate bool
	}{
		{[]domain.LogEntry{}, 0, 0, true},
		{[]domain.LogEntry{}, 1, 0, true},
		{[]domain.LogEntry{}, 0, 1, true},
		{[]domain.LogEntry{}, 1, 1, true},

		{[]domain.LogEntry{{Term: 1}}, 0, 0, false},
		{[]domain.LogEntry{{Term: 1}}, 1, 0, false},
		{[]domain.LogEntry{{Term: 1}}, 1, 1, true},
		{[]domain.LogEntry{{Term: 1}}, 1, 2, true},
		{[]domain.LogEntry{{Term: 1}}, 2, 0, false},
		{[]domain.LogEntry{{Term: 1}}, 2, 1, true},
		{[]domain.LogEntry{{Term: 1}}, 2, 2, true},

		{[]domain.LogEntry{{Term: 1}, {Term: 1}}, 1, 1, false},
		{[]domain.LogEntry{{Term: 1}, {Term: 1}}, 2, 1, true},
		{[]domain.LogEntry{{Term: 1}, {Term: 1}}, 2, 2, true},
		{[]domain.LogEntry{{Term: 1}, {Term: 2}}, 2, 2, true},
		{[]domain.LogEntry{{Term: 1}, {Term: 3}}, 2, 2, false},

		{[]domain.LogEntry{{Term: 1}, {Term: 2}}, 10, 2, true},
		{[]domain.LogEntry{{Term: 1}, {Term: 2}}, 10, 1, false},
	}
	for _, tt := range res {
		t.Run(fmt.Sprintf("%v", tt.logs), func(t *testing.T) {
			state := nodeState{
				fsmState: &fsmState{
					PersistentState: &PersistentState{
						machineLogs: tt.logs,
					},
				},
			}
			upToDate := state.IsUpToDate(tt.lastLogIndex, tt.lastLogTerm)
			if upToDate != tt.upToDate {
				t.Errorf("IsUpToDate got %v, want %v", upToDate, tt.upToDate)
			}
		})
	}
}

func TestCanGrantVote(t *testing.T) {
	res := []struct {
		votedFor string
		// input
		peerId string
		// ouput
		canGrantVote bool
	}{
		{"", "peerId", true},
		{"peerId", "peerId", true},
		{"peerId", "", false},
		{"unknownId", "peerId", false},
	}
	for _, tt := range res {
		t.Run(fmt.Sprintf("%v", tt.votedFor), func(t *testing.T) {
			state := nodeState{
				fsmState: &fsmState{
					PersistentState: &PersistentState{
						votedFor: tt.votedFor,
					},
				},
				peers: domain.Peers{"peerId": domain.Peer{}},
			}
			canGrantVote := state.CanGrantVote(tt.peerId)
			if canGrantVote != tt.canGrantVote {
				t.Errorf("CanGrantVote got %v, want %v", canGrantVote, tt.canGrantVote)
			}
		})
	}
}

func TestAppendEntriesInput(t *testing.T) {
	var currenTerm uint32 = 9
	logs := []domain.LogEntry{
		{Term: 1, Value: "val1"},
		{Term: 2, Value: "val2"},
		{Term: 3, Value: "val3"},
		{Term: 4, Value: "val4"},
		{Term: 5, Value: "val5"},
	}
	nextIndex := peerIndex{
		"peer1": 5,
		"peer2": 0,
		"peer3": 1,
		"peer4": 6,
		"peer5": 3,
	}
	matchIndex := peerIndex{
		"peer1": 0,
		"peer2": 3,
		"peer3": 5,
		"peer4": 2,
		"peer5": 3,
	}
	res := []struct {
		// input
		peerId string
		// output
		out domain.AppendEntriesInput
	}{
		{"peer1", domain.AppendEntriesInput{
			Term:         currenTerm,
			PrevLogIndex: 5,
			PrevLogTerm:  5,
			Entries:      []domain.LogEntry{},
		}},
		{"peer2", domain.AppendEntriesInput{
			Term:         currenTerm,
			PrevLogIndex: 0,
			PrevLogTerm:  0,
			Entries:      logs,
		}},
		{"peer3", domain.AppendEntriesInput{
			Term:         currenTerm,
			PrevLogIndex: 5,
			PrevLogTerm:  5,
			Entries:      []domain.LogEntry{},
		}},
		{"peer4", domain.AppendEntriesInput{
			Term:         currenTerm,
			PrevLogIndex: 6,
			PrevLogTerm:  currenTerm,
			Entries:      []domain.LogEntry{},
		}},
		{"peer5", domain.AppendEntriesInput{
			Term:         currenTerm,
			PrevLogIndex: 3,
			PrevLogTerm:  3,
			Entries:      logs[3:],
		}},
	}
	for _, tt := range res {
		t.Run(tt.peerId, func(t *testing.T) {
			state := nodeState{
				peers: domain.Peers{"peerId": domain.Peer{}},
				fsmState: &fsmState{
					matchIndex: matchIndex,
					nextIndex:  nextIndex,
					PersistentState: &PersistentState{
						currentTerm: currenTerm,
						machineLogs: logs,
					},
				},
			}
			out := state.AppendEntriesInput(tt.peerId)
			if !reflect.DeepEqual(tt.out, out) {
				t.Errorf("AppendEntriesInput got %v, want %v", out, tt.out)
			}
		})
	}
}

func TestComputeNewCommitIndex(t *testing.T) {
	logs := []domain.LogEntry{
		{Term: 1, Value: "val1"},
		{Term: 1, Value: "val2"},
		{Term: 3, Value: "val3"},
		{Term: 3, Value: "val4"},
		{Term: 5, Value: "val5"},
		{Term: 10, Value: "val6"},
	}
	// 4 + 1 nodes cluster
	p := domain.Peer{}
	peers := domain.Peers{
		"peer1": p,
		"peer2": p,
		"peer3": p,
		"peer4": p,
	}
	res := []struct {
		// state
		commitIndex uint32
		currentTerm uint32
		matchIndex  peerIndex
		// output
		newCommitIndex uint32
	}{
		{2, 5, peerIndex{"peer1": 6, "peer2": 10, "peer3": 1, "peer4": 1}, 5},
		{6, 12, peerIndex{"peer1": 1, "peer2": 1, "peer3": 1, "peer4": 1}, 6},
		{5, 10, peerIndex{"peer1": 6, "peer2": 6, "peer3": 1, "peer4": 1}, 6},
		{3, 3, peerIndex{"peer1": 4, "peer2": 4, "peer3": 1, "peer4": 1}, 4},
		{0, 10, peerIndex{"peer1": 6, "peer2": 6, "peer3": 6, "peer4": 6}, 6},
		{0, 10, peerIndex{"peer1": 5, "peer2": 5, "peer3": 5, "peer4": 5}, 0},
	}
	for _, tt := range res {
		t.Run(fmt.Sprintf("%v", tt.commitIndex), func(t *testing.T) {
			state := nodeState{
				peers: peers,
				fsmState: &fsmState{
					commitIndex: tt.commitIndex,
					matchIndex:  tt.matchIndex,
					PersistentState: &PersistentState{
						currentTerm: tt.currentTerm,
						machineLogs: logs,
					},
				},
			}
			newCommitIndex := state.ComputeNewCommitIndex()
			if tt.newCommitIndex != newCommitIndex {
				t.Errorf("ComputeNewCommitIndex got %v, want %v", newCommitIndex, tt.newCommitIndex)
			}
		})
	}
}
