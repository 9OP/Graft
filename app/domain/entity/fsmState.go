package entity

type fsmState struct {
	persistent
	commitIndex uint32
	lastApplied uint32
	// leader only
	nextIndex  map[string]uint32
	matchIndex map[string]uint32
}

func NewFsmState() *fsmState {
	return &fsmState{}
}

// Move into domain/service ?
// Move to node ?
func (s *fsmState) InitializeLeader(peers Peers) {
	lastLogIndex := s.GetLastLogIndex()
	nextIndex := make(map[string]uint32)
	matchIndex := make(map[string]uint32)

	for _, peer := range peers {
		nextIndex[peer.Id] = lastLogIndex + 1
		matchIndex[peer.Id] = 0
	}

	s.nextIndex = nextIndex
	s.matchIndex = matchIndex
}

func (s *fsmState) IncrementLastApplied() {
	s.lastApplied += 1
}

func (s *fsmState) SetCommitIndex(index uint32) {
	s.commitIndex = index
}

type FsmState struct {
	Persistent
	CommitIndex uint32
	LastApplied uint32
	NextIndex   map[string]uint32
	MatchIndex  map[string]uint32
}

func (s fsmState) getStateCopy() FsmState {
	persistent := s.getPersistentCopy()

	nextIndex := map[string]uint32{}
	for k, v := range s.nextIndex {
		nextIndex[k] = v
	}

	matchIndex := map[string]uint32{}
	for k, v := range s.matchIndex {
		matchIndex[k] = v
	}

	return FsmState{
		Persistent:  *persistent,
		CommitIndex: s.commitIndex,
		LastApplied: s.lastApplied,
		NextIndex:   nextIndex,
		MatchIndex:  matchIndex,
	}
}
