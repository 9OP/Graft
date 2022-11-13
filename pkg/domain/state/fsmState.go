package state

import "graft/pkg/utils"

type fsmState struct {
	commitIndex uint32
	lastApplied uint32
	nextIndex   map[string]uint32
	matchIndex  map[string]uint32
	*PersistentState
}

// Maps peerId to log index
type peerIndex map[string]uint32

func NewFsmState(persistent *PersistentState) fsmState {
	return fsmState{
		commitIndex:     0,
		lastApplied:     0,
		nextIndex:       peerIndex{},
		matchIndex:      peerIndex{},
		PersistentState: persistent,
	}
}

func (f fsmState) CommitIndex() uint32 {
	return f.commitIndex
}

func (f fsmState) LastApplied() uint32 {
	return f.lastApplied
}

func (f fsmState) NextIndex() peerIndex {
	return utils.CopyMap(f.nextIndex)
}

func (f fsmState) MatchIndex() peerIndex {
	return utils.CopyMap(f.matchIndex)
}

func (f fsmState) NextIndexForPeer(peerId string) uint32 {
	if idx, ok := f.nextIndex[peerId]; ok {
		return idx
	}
	return 0
}

func (f fsmState) MatchIndexForPeer(peerId string) uint32 {
	if idx, ok := f.matchIndex[peerId]; ok {
		return idx
	}
	return 0
}

func (f fsmState) WithCommitIndex(index uint32) (fsmState, bool) {
	changed := f.commitIndex != index
	f.commitIndex = index
	return f, changed
}

func (f fsmState) WithLastApplied(lastApplied uint32) fsmState {
	f.lastApplied = lastApplied
	return f
}

func (f fsmState) WithIncrementLastApplied() fsmState {
	f.lastApplied += 1
	return f
}

func (f fsmState) WithNextIndex(peerId string, index uint32) fsmState {
	nextIndex := f.NextIndex()
	nextIndex[peerId] = index
	f.nextIndex = nextIndex
	return f
}

func (f fsmState) WithMatchIndex(peerId string, index uint32) fsmState {
	matchIndex := f.MatchIndex()
	matchIndex[peerId] = index
	f.matchIndex = matchIndex
	return f
}

func (f fsmState) WithDecrementNextIndex(peerId string) (fsmState, bool) {
	// Copy of f.nextIndex is expensive
	// avoid it when next index is already 0
	if f.nextIndex[peerId] == 0 {
		return f, false
	}
	nextIndex := f.NextIndex()
	nextIndex[peerId] -= 1
	f.nextIndex = nextIndex
	return f, true
}
