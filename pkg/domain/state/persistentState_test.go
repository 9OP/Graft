package state

import (
	"errors"
	"fmt"
	"reflect"
	"testing"

	"graft/pkg/domain"
)

func TestNewPersistentState(t *testing.T) {
	var currentTerm uint32 = 12
	votedFor := "id"
	machineLogs := []domain.LogEntry{
		{Term: 1, Value: "val1"},
		{Term: 2},
		{Term: 3, Value: "val3", Type: "ADMIN"},
	}

	state := NewPersistentState(currentTerm, votedFor, machineLogs)

	if state.currentTerm != currentTerm {
		t.Fatalf("state.currentTerm = %v, want %v", state.currentTerm, currentTerm)
	}
	if state.votedFor != votedFor {
		t.Fatalf("state.currentTerm = %v, want %v", state.votedFor, votedFor)
	}

	for i := range state.machineLogs {
		log := state.machineLogs[i]
		expected := machineLogs[i]

		if log != expected {
			t.Fatalf("state.machineLogs[%d] = %v, want %v", i, log, expected)
		}
	}
}

func TestNewDefaultPersistentState(t *testing.T) {
	state := NewDefaultPersistentState()

	if state.currentTerm != 0 {
		t.Fatalf("state.currentTerm = %v, want %v", state.currentTerm, 0)
	}
	if state.votedFor != "" {
		t.Fatalf("state.votedFor = %v, want %v", state.votedFor, "")
	}
	if len(state.machineLogs) != 0 {
		t.Fatalf("len(state.machineLogs) = %v, want %v", len(state.machineLogs), 0)
	}
}

// func TestCurrentTerm(t *testing.T) {
// 	// TODO
// }
// func TestVotedFor(t *testing.T) {
// 	// TODO
// }

func TestLastLogIndex(t *testing.T) {
	// Last log index should be len(p.machineLogs) as logs
	// index start at 1 and not 0
	state := PersistentState{
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
	state := PersistentState{
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
	state := PersistentState{
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
			out, err := state.MachineLog(tt.in)
			// Check output and error
			if out != tt.out || !errors.Is(err, tt.err) {
				t.Errorf("got %v %v, want %v %v", out, err, tt.out, tt.err)
			}
		})
	}
}

func TestMachineLogs(t *testing.T) {
	state := PersistentState{
		machineLogs: []domain.LogEntry{
			{Term: 1, Value: "val1"},
			{Term: 2},
			{Term: 3, Value: "val3", Type: "ADMIN"},
		},
	}

	res := state.MachineLogs()
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
	state := PersistentState{
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
			out := state.MachineLogsFrom(tt.in)

			if !reflect.DeepEqual(out, tt.out) {
				t.Errorf("got %v, want %v", out, tt.out)
			}

			if len(out) != tt.len {
				t.Errorf("len got %v, want %v", len(out), tt.len)
			}
		})
	}

	// Mutate state
	log := state.MachineLogsFrom(0)[0]
	state.machineLogs[0].Term += 1
	if reflect.DeepEqual(log, state.machineLogs[0]) {
		t.Error("mutate copy")
	}
}

func TestWithers(t *testing.T) {
	state := PersistentState{votedFor: "", currentTerm: 0}
	stateWithTerm := state.WithCurrentTerm(10)
	stateWithVotedFor := state.WithVotedFor("id")

	// Does not mutate input
	if state.currentTerm != 0 {
		t.Error("state.currentTerm muted")
	}
	if state.votedFor != "" {
		t.Error("state.votedFor muted")
	}
	if stateWithTerm.currentTerm != 10 {
		t.Errorf("stateWithTerm.currentTerm got %v, want %v", stateWithTerm.currentTerm, 10)
	}
	if stateWithVotedFor.votedFor != "id" {
		t.Errorf("stateWithVotedFor.votedFor got %v, want %v", stateWithVotedFor.votedFor, 10)
	}

	// Can chain mutation
	newState := state.WithCurrentTerm(10).WithVotedFor("id")
	expect := PersistentState{votedFor: "id", currentTerm: 10}
	if !reflect.DeepEqual(newState, expect) {
		t.Errorf("newState got %v, want %v", newState, expect)
	}
}

func TestDeleteLogsFrom(t *testing.T) {
	//
}

func TestAppendLogsFrom(t *testing.T) {
	//
}
