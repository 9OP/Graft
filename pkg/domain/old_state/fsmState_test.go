package state

import (
	"reflect"
	"testing"
)

func TestNextMatchIndex(t *testing.T) {
	state := fsmState{
		nextIndex:  peerIndex{"peerId": 10},
		matchIndex: peerIndex{"peerId": 12},
	}
	nextIndex := state.NextIndex()
	matchIndex := state.MatchIndex()

	if !reflect.DeepEqual(nextIndex, state.nextIndex) {
		t.Errorf("got %v, want %v", nextIndex, state.nextIndex)
	}
	if !reflect.DeepEqual(matchIndex, state.matchIndex) {
		t.Errorf("got %v, want %v", matchIndex, state.matchIndex)
	}

	// Mutate
	state.nextIndex["peerId"] += 1
	state.matchIndex["peerId"] += 1
	if reflect.DeepEqual(nextIndex, state.nextIndex) {
		t.Error("mutate copy")
	}
	if reflect.DeepEqual(matchIndex, state.matchIndex) {
		t.Error("mutate copy")
	}
}

func TestNextMatchIndexForPeer(t *testing.T) {
	state := fsmState{
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
	state := fsmState{
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
		{"newPeerId", 2, peerIndex{"newPeerId": 2, "peerId": 10}},
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
	state := fsmState{
		nextIndex: peerIndex{"id1": 10, "id2": 0, "id3": 1},
	}

	res := []struct {
		peerId string
		out    uint32
		ok     bool
	}{
		{"id1", 9, true},
		{"id2", 0, false},
		{"id3", 0, true},
		{"id4", 0, false},
	}
	for _, tt := range res {
		t.Run(tt.peerId, func(t *testing.T) {
			out, ok := state.WithDecrementNextIndex(tt.peerId)
			if !reflect.DeepEqual(out.nextIndex[tt.peerId], tt.out) {
				t.Errorf("got %v, want %v", out.nextIndex[tt.peerId], tt.out)
			}
			if ok != tt.ok {
				t.Errorf("got %v, want %v", ok, tt.ok)
			}
		})
	}
}
