package statenew

import (
	"errors"
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
			state := state{leaderId: tt.leaderId, peers: tt.peers}
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
			state := state{peers: tt.peers}
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
			state := state{
				machineLogs: tt.logs,
			}
			upToDate := state.IsLogUpToDate(tt.lastLogIndex, tt.lastLogTerm)
			if upToDate != tt.upToDate {
				t.Errorf("IsLogUpToDate got %v, want %v", upToDate, tt.upToDate)
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
			state := state{
				votedFor: tt.votedFor,
				peers:    domain.Peers{"peerId": domain.Peer{}},
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
			state := state{
				peers:       domain.Peers{"peerId": domain.Peer{}},
				matchIndex:  matchIndex,
				nextIndex:   nextIndex,
				currentTerm: currenTerm,
				machineLogs: logs,
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
			state := state{
				peers:       peers,
				commitIndex: tt.commitIndex,
				matchIndex:  tt.matchIndex,
				currentTerm: tt.currentTerm,
				machineLogs: logs,
			}
			newCommitIndex := state.ComputeNewCommitIndex()
			if tt.newCommitIndex != newCommitIndex {
				t.Errorf("ComputeNewCommitIndex got %v, want %v", newCommitIndex, tt.newCommitIndex)
			}
		})
	}
}

func TestNextMatchIndexForPeer(t *testing.T) {
	state := state{
		nextIndex:  peerIndex{"peerId": 10},
		matchIndex: peerIndex{"peerId": 10},
	}

	res := []struct {
		in  string
		idx uint32
	}{
		{"peerId", 10},
		{"noId", 0},
	}

	for _, tt := range res {
		t.Run(tt.in, func(t *testing.T) {
			idx := state.NextIndexForPeer(tt.in)
			if idx != tt.idx {
				t.Errorf("got %v, want %v", idx, tt.idx)
			}
			idx = state.MatchIndexForPeer(tt.in)
			if idx != tt.idx {
				t.Errorf("got %v, want %v", idx, tt.idx)
			}
		})
	}
}

func TestWithNextMatchIndex(t *testing.T) {
	state := state{
		nextIndex:  peerIndex{"peerId": 10},
		matchIndex: peerIndex{"peerId": 10},
	}

	res := []struct {
		peerId string
		idx    uint32
		out    peerIndex
	}{
		{"peerId", 12, peerIndex{"peerId": 12}},
		{"peerId", 10, peerIndex{"peerId": 10}},
		{"newPeerId", 0, peerIndex{"peerId": 10}},
	}
	for _, tt := range res {
		t.Run(tt.peerId, func(t *testing.T) {
			out := state.WithNextIndex(tt.peerId, tt.idx)
			if !reflect.DeepEqual(out.nextIndex, tt.out) {
				t.Errorf("got %v, want %v", out.nextIndex, tt.out)
			}
			out = state.WithMatchIndex(tt.peerId, tt.idx)
			if !reflect.DeepEqual(out.matchIndex, tt.out) {
				t.Errorf("got %v, want %v", out.matchIndex, tt.out)
			}
		})
	}

	// Mutate
	newState := state.WithNextIndex("peerId", 2).WithMatchIndex("peerId", 2)
	state.nextIndex["peerId"] += 1
	state.matchIndex["peerId"] += 1

	if newState.nextIndex["peerId"] != 2 {
		t.Error("mutate copy")
	}
	if newState.matchIndex["peerId"] != 2 {
		t.Error("mutate copy")
	}
}

func TestWithDecrementNextIndex(t *testing.T) {
	st := state{
		nextIndex: peerIndex{"id1": 10, "id2": 0, "id3": 1},
	}

	res := []struct {
		peerId string
		out    uint32
	}{
		{"id1", 9},
		{"id2", 0},
		{"id3", 0},
		{"id4", 0},
	}
	for _, tt := range res {
		t.Run(tt.peerId, func(t *testing.T) {
			out := st.WithDecrementNextIndex(tt.peerId)
			if !reflect.DeepEqual(out.nextIndex[tt.peerId], tt.out) {
				t.Errorf("got %v, want %v", out.nextIndex[tt.peerId], tt.out)
			}
		})
	}
}

func TestLastLogIndex(t *testing.T) {
	// Last log index should be len(p.machineLogs) as logs
	// index start at 1 and not 0
	state := state{
		machineLogs: []domain.LogEntry{
			{Term: 1, Value: "val1"},
			{Term: 2},
			{Term: 3, Value: "val3", Type: "ADMIN"},
		},
	}

	if state.LastLogIndex() != 3 {
		t.Fatalf("state.LastLogIndex() = %v, want %v", state.LastLogIndex(), 3)
	}
}

func TestLastLog(t *testing.T) {
	state := state{
		machineLogs: []domain.LogEntry{
			{Term: 1, Value: "val1"},
			{Term: 2},
			{Term: 3, Value: "val3", Type: "ADMIN"},
		},
	}

	lastLog := state.LastLog()
	if lastLog != state.machineLogs[2] {
		t.Fatalf("state.LastLog() = %v, want %v", state.machineLogs[2], lastLog)
	}
	if &lastLog == &state.machineLogs[2] {
		t.Fatalf("&state.LastLog() = %p, dont want %p", &state.machineLogs[2], &lastLog)
	}

	// Mutate last log in state
	state.machineLogs[2].Term = 1
	state.machineLogs[2].Value = "val1"

	updatedLastLog := state.LastLog()

	if lastLog == updatedLastLog {
		t.Fatalf("state.LastLog() = %v, dont want %v", updatedLastLog, lastLog)
	}
	if &lastLog == &updatedLastLog {
		t.Fatalf("&updatedLastLog = %p, dont want &lastLog = %p", &updatedLastLog, &lastLog)
	}
}

func TestMachineLog(t *testing.T) {
	state := state{
		machineLogs: []domain.LogEntry{
			{Term: 1, Value: "val1"},
			{Term: 2},
			{Term: 3, Value: "val3", Type: "ADMIN"},
		},
	}

	res := []struct {
		in  uint32
		out domain.LogEntry
		err error
	}{
		{0, domain.LogEntry{}, nil},
		{1, domain.LogEntry{Term: 1, Value: "val1"}, nil},
		{3, domain.LogEntry{Term: 3, Value: "val3", Type: "ADMIN"}, nil},
		{4, domain.LogEntry{}, errIndexOutOfRange},
	}

	for _, tt := range res {
		t.Run(fmt.Sprintf("%d", tt.in), func(t *testing.T) {
			out, err := state.Log(tt.in)
			// Check output and error
			if out != tt.out || !errors.Is(err, tt.err) {
				t.Errorf("got %v %v, want %v %v", out, err, tt.out, tt.err)
			}
		})
	}

	// Mutate
	log, _ := state.Log(1)
	state.machineLogs[0].Term += 1
	if reflect.DeepEqual(log, state.machineLogs[0]) {
		t.Error("mutate copy")
	}
}

func TestMachineLogs(t *testing.T) {
	state := state{
		machineLogs: []domain.LogEntry{
			{Term: 1, Value: "val1"},
			{Term: 2},
			{Term: 3, Value: "val3", Type: "ADMIN"},
		},
	}

	res := state.Logs()
	if !reflect.DeepEqual(res, state.machineLogs) {
		t.Errorf("got %v, want %v", res, state.machineLogs)
	}

	// Mutate state
	state.machineLogs[0].Term += 1
	if reflect.DeepEqual(res, state.machineLogs) {
		t.Error("mutate copy")
	}
}

func TestMachineLogsFrom(t *testing.T) {
	state := state{
		machineLogs: []domain.LogEntry{
			{Term: 1, Value: "val1"},
			{Term: 2, Value: "val2"},
			{Term: 3, Value: "val3"},
			{Term: 4, Value: "val4"},
			{Term: 5, Value: "val5"},
		},
	}

	res := []struct {
		in  uint32
		out []domain.LogEntry
		len int
	}{
		{0, state.machineLogs, 5},
		{1, state.machineLogs, 5},
		{3, state.machineLogs[2:], 3},
		{4, state.machineLogs[3:], 2},
		{5, state.machineLogs[4:], 1},
		{6, []domain.LogEntry{}, 0},
	}

	for _, tt := range res {
		t.Run(fmt.Sprintf("%d", tt.in), func(t *testing.T) {
			out := state.LogsFrom(tt.in)

			if !reflect.DeepEqual(out, tt.out) {
				t.Errorf("got %v, want %v", out, tt.out)
			}

			if len(out) != tt.len {
				t.Errorf("len got %v, want %v", len(out), tt.len)
			}
		})
	}

	// Mutate state
	log := state.LogsFrom(0)[0]
	state.machineLogs[0].Term += 1
	if reflect.DeepEqual(log, state.machineLogs[0]) {
		t.Error("mutate copy")
	}
}

func TestWithers(t *testing.T) {
	st := state{votedFor: "", currentTerm: 0}
	stateWithTerm := st.WithCurrentTerm(10)
	stateWithVotedFor := st.WithVotedFor("id")

	// Does not mutate input
	if st.currentTerm != 0 {
		t.Error("state.currentTerm muted")
	}
	if st.votedFor != "" {
		t.Error("state.votedFor muted")
	}
	if stateWithTerm.currentTerm != 10 {
		t.Errorf("stateWithTerm.currentTerm got %v, want %v", stateWithTerm.currentTerm, 10)
	}
	if stateWithVotedFor.votedFor != "id" {
		t.Errorf("stateWithVotedFor.votedFor got %v, want %v", stateWithVotedFor.votedFor, 10)
	}

	// Can chain mutation
	newState := st.WithCurrentTerm(10).WithVotedFor("id")
	expect := state{votedFor: "id", currentTerm: 10}
	if !reflect.DeepEqual(newState, expect) {
		t.Errorf("newState got %v, want %v", newState, expect)
	}
}

func TestWithDeleteLogsFrom(t *testing.T) {
	state := state{
		machineLogs: []domain.LogEntry{
			{Term: 1, Value: "val1"},
			{Term: 2, Value: "val2"},
			{Term: 3, Value: "val3"},
			{Term: 4, Value: "val4"},
			{Term: 5, Value: "val5"},
		},
	}

	res := []struct {
		in      uint32
		out     []domain.LogEntry
		len     int
		changed bool
	}{
		{0, []domain.LogEntry{}, 0, true},
		{1, []domain.LogEntry{}, 0, true},
		{3, state.machineLogs[:2], 2, true},
		{4, state.machineLogs[:3], 3, true},
		{5, state.machineLogs[:4], 4, true},
		{6, state.machineLogs, 5, false},
		{10, state.machineLogs, 5, false},
	}

	for _, tt := range res {
		t.Run(fmt.Sprintf("%d", tt.in), func(t *testing.T) {
			out, changed := state.WithDeleteLogsFrom(tt.in)

			if !reflect.DeepEqual(out.machineLogs, tt.out) {
				t.Errorf("got %v, want %v", out.machineLogs, tt.out)
			}

			if len(out.machineLogs) != tt.len {
				t.Errorf("len got %v, want %v", len(out.machineLogs), tt.len)
			}

			if changed != tt.changed {
				t.Errorf("changed got %v, want %v", changed, tt.changed)
			}
		})
	}

	// Mutate
	newState, _ := state.WithDeleteLogsFrom(10)
	state.machineLogs[0].Term += 1
	if reflect.DeepEqual(newState.machineLogs, state.machineLogs) {
		t.Error("mutate copy")
	}
}

func TestWithAppendLogsFrom(t *testing.T) {
	state := state{
		machineLogs: []domain.LogEntry{
			{Term: 1},
			{Term: 2},
			{Term: 3},
		},
	}
	entries := []domain.LogEntry{
		{Term: 4},
		{Term: 5},
	}
	res := []struct {
		prevLogIndex uint32
		entries      []domain.LogEntry
		logs         []domain.LogEntry
		len          int
		changed      bool
	}{
		{3, entries, []domain.LogEntry{{Term: 1}, {Term: 2}, {Term: 3}, {Term: 4}, {Term: 5}}, 5, true},
		{2, entries, []domain.LogEntry{{Term: 1}, {Term: 2}, {Term: 3}, {Term: 5}}, 4, true},
		{1, entries, []domain.LogEntry{{Term: 1}, {Term: 2}, {Term: 3}}, 3, false},
		{0, entries, []domain.LogEntry{{Term: 1}, {Term: 2}, {Term: 3}}, 3, false},
		{4, entries, []domain.LogEntry{{Term: 1}, {Term: 2}, {Term: 3}}, 3, false},
		{5, entries, []domain.LogEntry{{Term: 1}, {Term: 2}, {Term: 3}}, 3, false},
		{2, []domain.LogEntry{}, []domain.LogEntry{{Term: 1}, {Term: 2}, {Term: 3}}, 3, false},
	}

	for _, tt := range res {
		t.Run(fmt.Sprintf("%d", tt.prevLogIndex), func(t *testing.T) {
			out, changed := state.WithAppendLogs(tt.prevLogIndex, tt.entries...)

			if !reflect.DeepEqual(out.machineLogs, tt.logs) {
				t.Errorf("got %v, want %v", out.machineLogs, tt.logs)
			}

			if changed != tt.changed {
				t.Errorf("got %v, want %v", changed, tt.changed)
			}

			if len(out.machineLogs) != len(tt.logs) {
				t.Errorf("got %v, want %v", len(out.machineLogs), len(tt.logs))
			}
		})
	}
}