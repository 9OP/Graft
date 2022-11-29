package secondaryAdapter

import (
	"encoding/json"
	"os"
	"sync"

	"graft/pkg/domain"
)

type jsonPersister struct {
	mu sync.Mutex
}

func NewJsonPersister() *jsonPersister {
	return &jsonPersister{}
}

func (p *jsonPersister) Load(location string) (*domain.PersistentState, error) {
	data, err := os.ReadFile(location)
	if err != nil {
		return nil, err
	}
	state := &domain.PersistentState{}
	err = json.Unmarshal(data, state)
	return state, err
}

func (p *jsonPersister) Save(location string, currentTerm uint32, votedFor string, machineLogs []domain.LogEntry) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	state := domain.PersistentState{
		CurrentTerm: currentTerm,
		VotedFor:    votedFor,
		MachineLogs: machineLogs,
	}
	data, err := json.MarshalIndent(&state, "", "  ")
	if err != nil {
		return err
	}
	// return nil
	return os.WriteFile(location, data, 0o644)
}
