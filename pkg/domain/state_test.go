package domain

import (
	"errors"
	"fmt"
	"reflect"
	"testing"
)

func TestHasLeader(t *testing.T) {
	peer := Peer{}
	res := []struct {
		// state
		leaderId string
		peers    Peers
		// output
		hasLeader bool
	}{
		{"leaderId", Peers{}, false},
		{"leaderId", Peers{"leaderId": peer}, true},
		{"leaderId", Peers{"id": peer}, false},
		{"", Peers{"leaderId": peer}, false},
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

func TestActivePeers(t *testing.T) {
	peers := Peers{
		"1": Peer{Active: true},
		"2": Peer{Active: false, Id: "2"},
		"3": Peer{Active: false},
		"4": Peer{},
	}
	expected := Peers{
		"1": Peer{Active: true},
	}

	state := state{peers: peers}
	activePeers := state.activePeers()

	if !reflect.DeepEqual(activePeers, expected) {
		t.Errorf("active peers got %v, want %v", activePeers, expected)
	}

	// Mutate
	p := peers["2"]
	p.Active = true
	peers["2"] = p

	if activePeers["2"].Active == true {
		t.Errorf("mutate copy")
	}
}

func TestQuorum(t *testing.T) {
	p1 := Peer{Active: true}
	p2 := Peer{Active: false}
	res := []struct {
		// state
		peers Peers
		// output
		quorum int
	}{
		{Peers{"peer1": p1, "peer2": p1}, 2},                           // tolerate 1 failure
		{Peers{"peer1": p1, "peer2": p1, "peer3": p1}, 3},              // tolerate 1 failure
		{Peers{"peer1": p1, "peer2": p1, "peer3": p1, "peer4": p1}, 3}, // tolerate 2 failures
		{Peers{"peer1": p1}, 2},                                        // tolerate 0 failure
		{Peers{"peer1": p1, "peer2": p2, "peer3": p2}, 2},
		{Peers{}, 1}, // tolerate 0 failure
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
		logs []LogEntry
		// input
		lastLogIndex uint32
		lastLogTerm  uint32
		// output
		upToDate bool
	}{
		{[]LogEntry{}, 0, 0, true},
		{[]LogEntry{}, 1, 0, true},
		{[]LogEntry{}, 0, 1, true},
		{[]LogEntry{}, 1, 1, true},

		{[]LogEntry{{Term: 1}}, 0, 0, false},
		{[]LogEntry{{Term: 1}}, 1, 0, false},
		{[]LogEntry{{Term: 1}}, 1, 1, true},
		{[]LogEntry{{Term: 1}}, 1, 2, true},
		{[]LogEntry{{Term: 1}}, 2, 0, false},
		{[]LogEntry{{Term: 1}}, 2, 1, true},
		{[]LogEntry{{Term: 1}}, 2, 2, true},

		{[]LogEntry{{Term: 1}, {Term: 1}}, 1, 1, false},
		{[]LogEntry{{Term: 1}, {Term: 1}}, 2, 1, true},
		{[]LogEntry{{Term: 1}, {Term: 1}}, 2, 2, true},
		{[]LogEntry{{Term: 1}, {Term: 2}}, 2, 2, true},
		{[]LogEntry{{Term: 1}, {Term: 3}}, 2, 2, false},

		{[]LogEntry{{Term: 1}, {Term: 2}}, 10, 2, true},
		{[]LogEntry{{Term: 1}, {Term: 2}}, 10, 1, false},
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
				peers:    Peers{"peerId": Peer{}},
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
	logs := []LogEntry{
		{Term: 1, Data: []byte("val1")},
		{Term: 2, Data: []byte("val2")},
		{Term: 3, Data: []byte("val3")},
		{Term: 4, Data: []byte("val4")},
		{Term: 5, Data: []byte("val5")},
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
		out AppendEntriesInput
	}{
		{"peer1", AppendEntriesInput{
			Term:         currenTerm,
			PrevLogIndex: 5,
			PrevLogTerm:  5,
			Entries:      []LogEntry{},
		}},
		{"peer2", AppendEntriesInput{
			Term:         currenTerm,
			PrevLogIndex: 0,
			PrevLogTerm:  0,
			Entries:      logs,
		}},
		{"peer3", AppendEntriesInput{
			Term:         currenTerm,
			PrevLogIndex: 5,
			PrevLogTerm:  5,
			Entries:      []LogEntry{},
		}},
		{"peer4", AppendEntriesInput{
			Term:         currenTerm,
			PrevLogIndex: 6,
			PrevLogTerm:  currenTerm,
			Entries:      []LogEntry{},
		}},
		{"peer5", AppendEntriesInput{
			Term:         currenTerm,
			PrevLogIndex: 3,
			PrevLogTerm:  3,
			Entries:      logs[3:],
		}},
	}
	for _, tt := range res {
		t.Run(tt.peerId, func(t *testing.T) {
			state := state{
				peers:       Peers{"peerId": Peer{}},
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
	logs := []LogEntry{
		{Term: 1, Data: []byte("val1")},
		{Term: 1, Data: []byte("val2")},
		{Term: 3, Data: []byte("val3")},
		{Term: 3, Data: []byte("val4")},
		{Term: 5, Data: []byte("val5")},
		{Term: 10, Data: []byte("val6")},
	}
	// 4 + 1 nodes cluster
	p := Peer{Active: true}
	peers := Peers{
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
			newCommitIndex := state.computeNewCommitIndex()
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
			idx := state.nextIndexForPeer(tt.in)
			if idx != tt.idx {
				t.Errorf("got %v, want %v", idx, tt.idx)
			}
			idx = state.matchIndexForPeer(tt.in)
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
		{"newPeerId", 0, peerIndex{"peerId": 10, "newPeerId": 0}},
	}
	for _, tt := range res {
		t.Run(tt.peerId, func(t *testing.T) {
			out := state.withNextIndex(tt.peerId, tt.idx)
			if !reflect.DeepEqual(out.nextIndex, tt.out) {
				t.Errorf("got %v, want %v", out.nextIndex, tt.out)
			}
			out = state.withMatchIndex(tt.peerId, tt.idx)
			if !reflect.DeepEqual(out.matchIndex, tt.out) {
				t.Errorf("got %v, want %v", out.matchIndex, tt.out)
			}
		})
	}

	// Mutate
	newState := state.withNextIndex("peerId", 2).withMatchIndex("peerId", 2)
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
			out := st.withDecrementNextIndex(tt.peerId)
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
		machineLogs: []LogEntry{
			{Term: 1, Data: []byte("val1")},
			{Term: 2},
			{Term: 3, Data: []byte("val3"), Type: LogNoop},
		},
	}

	if state.lastLogIndex() != 3 {
		t.Fatalf("state.lastLogIndex() = %v, want %v", state.lastLogIndex(), 3)
	}
}

func TestLastLog(t *testing.T) {
	state := state{
		machineLogs: []LogEntry{
			{Term: 1, Data: []byte("val1")},
			{Term: 2},
			{Term: 3, Data: []byte("val3"), Type: LogNoop},
		},
	}

	lastLog := state.lastLog()
	if !reflect.DeepEqual(lastLog, state.machineLogs[2]) {
		t.Fatalf("state.LastLog() = %v, want %v", state.machineLogs[2], lastLog)
	}
	if &lastLog == &state.machineLogs[2] {
		t.Fatalf("&state.LastLog() = %p, dont want %p", &state.machineLogs[2], &lastLog)
	}

	// Mutate last log in state
	state.machineLogs[2].Term = 1
	state.machineLogs[2].Data = []byte("val1")

	updatedLastLog := state.lastLog()

	if reflect.DeepEqual(lastLog, updatedLastLog) {
		t.Fatalf("state.LastLog() = %v, dont want %v", updatedLastLog, lastLog)
	}
	if &lastLog == &updatedLastLog {
		t.Fatalf("&updatedLastLog = %p, dont want &lastLog = %p", &updatedLastLog, &lastLog)
	}
}

func TestMachineLog(t *testing.T) {
	state := state{
		machineLogs: []LogEntry{
			{Term: 1, Data: []byte("val1")},
			{Term: 2},
			{Term: 3, Data: []byte("val3"), Type: LogNoop},
		},
	}

	res := []struct {
		in  uint32
		out LogEntry
		err error
	}{
		{0, LogEntry{}, nil},
		{1, LogEntry{Term: 1, Data: []byte("val1")}, nil},
		{3, LogEntry{Term: 3, Data: []byte("val3"), Type: LogNoop}, nil},
		{4, LogEntry{}, errIndexOutOfRange},
	}

	for _, tt := range res {
		t.Run(fmt.Sprintf("%d", tt.in), func(t *testing.T) {
			out, err := state.Log(tt.in)
			// Check output and error
			if !reflect.DeepEqual(out, tt.out) || !errors.Is(err, tt.err) {
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
		machineLogs: []LogEntry{
			{Term: 1, Data: []byte("val1")},
			{Term: 2},
			{Term: 3, Data: []byte("val3"), Type: LogNoop},
		},
	}

	res := state.logs()
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
		machineLogs: []LogEntry{
			{Term: 1, Data: []byte("val1")},
			{Term: 2, Data: []byte("val2")},
			{Term: 3, Data: []byte("val3")},
			{Term: 4, Data: []byte("val4")},
			{Term: 5, Data: []byte("val5")},
		},
	}

	res := []struct {
		in  uint32
		out []LogEntry
		len int
	}{
		{0, state.machineLogs, 5},
		{1, state.machineLogs, 5},
		{3, state.machineLogs[2:], 3},
		{4, state.machineLogs[3:], 2},
		{5, state.machineLogs[4:], 1},
		{6, []LogEntry{}, 0},
	}

	for _, tt := range res {
		t.Run(fmt.Sprintf("%d", tt.in), func(t *testing.T) {
			out := state.logsFrom(tt.in)

			if !reflect.DeepEqual(out, tt.out) {
				t.Errorf("got %v, want %v", out, tt.out)
			}

			if len(out) != tt.len {
				t.Errorf("len got %v, want %v", len(out), tt.len)
			}
		})
	}

	// Mutate state
	log := state.logsFrom(0)[0]
	state.machineLogs[0].Term += 1
	if reflect.DeepEqual(log, state.machineLogs[0]) {
		t.Error("mutate copy")
	}
}

func TestWithers(t *testing.T) {
	st := state{votedFor: "", currentTerm: 0}
	stateWithTerm := st.withCurrentTerm(10)
	stateWithVotedFor := st.withVotedFor("id")

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
	newState := st.withCurrentTerm(10).withVotedFor("id")
	expect := state{votedFor: "id", currentTerm: 10}
	if !reflect.DeepEqual(newState, expect) {
		t.Errorf("newState got %v, want %v", newState, expect)
	}
}

func TestWithDeleteLogsFrom(t *testing.T) {
	state := state{
		machineLogs: []LogEntry{
			{Term: 1, Data: []byte("val1")},
			{Term: 2, Data: []byte("val2")},
			{Term: 3, Data: []byte("val3")},
			{Term: 4, Data: []byte("val4")},
			{Term: 5, Data: []byte("val5")},
		},
	}

	res := []struct {
		in      uint32
		out     []LogEntry
		len     int
		changed bool
	}{
		{0, []LogEntry{}, 0, true},
		{1, []LogEntry{}, 0, true},
		{3, state.machineLogs[:2], 2, true},
		{4, state.machineLogs[:3], 3, true},
		{5, state.machineLogs[:4], 4, true},
		{6, state.machineLogs, 5, false},
		{10, state.machineLogs, 5, false},
	}

	for _, tt := range res {
		t.Run(fmt.Sprintf("%d", tt.in), func(t *testing.T) {
			out, changed := state.withDeleteLogsFrom(tt.in)

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
	newState, _ := state.withDeleteLogsFrom(10)
	state.machineLogs[0].Term += 1
	if reflect.DeepEqual(newState.machineLogs, state.machineLogs) {
		t.Error("mutate copy")
	}
}

func TestWithAppendLogsFrom(t *testing.T) {
	state := state{
		machineLogs: []LogEntry{
			{Term: 1},
			{Term: 2},
			{Term: 3},
		},
	}
	entries := []LogEntry{
		{Term: 4},
		{Term: 5},
	}
	res := []struct {
		prevLogIndex uint32
		entries      []LogEntry
		logs         []LogEntry
		len          int
		changed      bool
	}{
		{3, entries, []LogEntry{{Term: 1}, {Term: 2}, {Term: 3}, {Term: 4}, {Term: 5}}, 5, true},
		{2, entries, []LogEntry{{Term: 1}, {Term: 2}, {Term: 3}, {Term: 5}}, 4, true},
		{1, entries, []LogEntry{{Term: 1}, {Term: 2}, {Term: 3}}, 3, false},
		{0, entries, []LogEntry{{Term: 1}, {Term: 2}, {Term: 3}}, 3, false},
		{4, entries, []LogEntry{{Term: 1}, {Term: 2}, {Term: 3}}, 3, false},
		{5, entries, []LogEntry{{Term: 1}, {Term: 2}, {Term: 3}}, 3, false},
		{2, []LogEntry{}, []LogEntry{{Term: 1}, {Term: 2}, {Term: 3}}, 3, false},
	}

	for _, tt := range res {
		t.Run(fmt.Sprintf("%d", tt.prevLogIndex), func(t *testing.T) {
			out, changed := state.withAppendLogs(tt.prevLogIndex, tt.entries...)

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
