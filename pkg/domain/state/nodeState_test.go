package state

import (
	"fmt"
	"testing"

	"graft/pkg/domain"
)

func TestHasLeader(t *testing.T) {
	peer := domain.Peer{}
	res := []struct {
		leaderId  string
		peers     domain.Peers
		hasLeader bool
	}{
		{"leaderId", domain.Peers{}, false},
		{"leaderId", domain.Peers{"leaderId": peer}, true},
		{"leaderId", domain.Peers{"id": peer}, false},
		{"", domain.Peers{"leaderId": peer}, false},
	}
	for _, tt := range res {
		t.Run(tt.leaderId, func(t *testing.T) {
			ok := nodeState{leaderId: tt.leaderId, peers: tt.peers}.HasLeader()
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
		peers  domain.Peers
		quorum int
	}{
		{domain.Peers{"peer1": p, "peer2": p}, 2},
		{domain.Peers{"peer1": p, "peer2": p, "peer3": p}, 2},
		{domain.Peers{"peer1": p, "peer2": p, "peer3": p, "peer4": p}, 3},
		{domain.Peers{"peer1": p}, 1},
		{domain.Peers{}, 1},
	}
	for _, tt := range res {
		t.Run(fmt.Sprintf("%v", tt.peers), func(t *testing.T) {
			quorum := nodeState{peers: tt.peers}.Quorum()
			if quorum != tt.quorum {
				t.Errorf("Quorum got %v, want %v", quorum, tt.quorum)
			}
		})
	}
}

func TestDoesLogExists(t *testing.T) {
	res := []struct {
		logs         []domain.LogEntry
		lastLogIndex uint32
		lastLogTerm  uint32
		upToDate     bool
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
			upToDate := nodeState{
				fsmState: &fsmState{
					PersistentState: &PersistentState{
						machineLogs: tt.logs,
					},
				},
			}.IsUpToDate(tt.lastLogIndex, tt.lastLogTerm)

			if upToDate != tt.upToDate {
				t.Errorf("IsUpToDate got %v, want %v", upToDate, tt.upToDate)
			}
		})
	}
}

func TestCanGrantVote(t *testing.T) {
}

func TestAppendEntriesInput(t *testing.T) {
}

func TestComputeNewCommitIndex(t *testing.T) {
}
