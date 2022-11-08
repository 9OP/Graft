package secondaryAdapter

import (
	"encoding/json"
	"os"

	"graft/pkg/domain/entity"
)

type UseCaseJsonPersisterAdapter interface {
	Load(location string) (*persistent, error)
	Save(location string, currentTerm uint32, votedFor string, machineLogs []entity.LogEntry) error
}

type jsonPersister struct{}

func NewJsonPersister() *jsonPersister {
	return &jsonPersister{}
}

type persistent struct {
	CurrentTerm uint32            `json:"current_term"`
	VotedFor    string            `json:"voted_for"`
	MachineLogs []entity.LogEntry `json:"machine_logs"`
}

func (p jsonPersister) Load(location string) (*persistent, error) {
	data, err := os.ReadFile(location)
	if err != nil {
		return nil, err
	}
	state := &persistent{}
	err = json.Unmarshal(data, state)
	return state, err
}

func (p jsonPersister) Save(location string, currentTerm uint32, votedFor string, machineLogs []entity.LogEntry) error {
	state := persistent{
		CurrentTerm: currentTerm,
		VotedFor:    votedFor,
		MachineLogs: machineLogs,
	}
	data, err := json.MarshalIndent(&state, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(location, data, 0o644)
}
