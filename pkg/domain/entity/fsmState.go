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

func (f FsmState) NextIndexForPeer(peerId string) uint32 {
	// Not safe when peerId not in nextIndex
	return f.nextIndex[peerId]
}

func (f FsmState) MatchIndexForPeer(peerId string) uint32 {
	// Not safe when peerId not in nextIndex
	return f.matchIndex[peerId]
}

func (f FsmState) MatchIndex() peerIndex {
	return utils.CopyMap(f.matchIndex)
}

func (f FsmState) WithCommitIndex(index uint32) FsmState {
	f.commitIndex = index
	return f
}

func (f FsmState) WithLastApplied(index uint32) FsmState {
	f.lastApplied = index
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

func (f FsmState) WithDecrementNextIndex(peerId string) FsmState {
	nextIndex := f.NextIndex()
	nextIndex[peerId] -= 1
	f.nextIndex = nextIndex
	return f
}
