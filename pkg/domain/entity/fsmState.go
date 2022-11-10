package entity

import utils "graft/pkg/domain"

type FsmState struct {
	commitIndex uint32
	lastApplied uint32
	nextIndex   map[string]uint32
	matchIndex  map[string]uint32
	*PersistentState
}

// Maps peerId to log index
type peerIndex map[string]uint32

func NewFsmState(persistent *PersistentState) FsmState {
	return FsmState{
		commitIndex:     0,
		lastApplied:     0,
		nextIndex:       peerIndex{},
		matchIndex:      peerIndex{},
		PersistentState: persistent,
	}
}

func (f FsmState) CommitIndex() uint32 {
	return f.commitIndex
}

func (f FsmState) LastApplied() uint32 {
	return f.lastApplied
}

func (f FsmState) NextIndex() peerIndex {
	return utils.CopyMap(f.nextIndex)
}

func (f FsmState) MatchIndex() peerIndex {
	return utils.CopyMap(f.matchIndex)
}

func (f FsmState) NextIndexForPeer(peerId string) uint32 {
	if idx, ok := f.nextIndex[peerId]; ok {
		return idx
	}
	return 0
}

func (f FsmState) MatchIndexForPeer(peerId string) uint32 {
	if idx, ok := f.matchIndex[peerId]; ok {
		return idx
	}
	return 0
}

func (f FsmState) WithCommitIndex(index uint32) (FsmState, bool) {
	changed := f.commitIndex != index
	f.commitIndex = index
	return f, changed
}

func (f FsmState) WithLastApplied(lastApplied uint32) FsmState {
	f.lastApplied = lastApplied
	return f
}

func (f FsmState) WithIncrementLastApplied() FsmState {
	f.lastApplied += 1
	return f
}

func (f FsmState) WithNextIndex(peerId string, index uint32) FsmState {
	nextIndex := f.NextIndex()
	nextIndex[peerId] = index
	f.nextIndex = nextIndex
	return f
}

func (f FsmState) WithMatchIndex(peerId string, index uint32) FsmState {
	matchIndex := f.MatchIndex()
	matchIndex[peerId] = index
	f.matchIndex = matchIndex
	return f
}

func (f FsmState) WithDecrementNextIndex(peerId string) (FsmState, bool) {
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
