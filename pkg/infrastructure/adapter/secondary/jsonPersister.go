package secondaryAdapter

import (
	"encoding/json"
	"graft/pkg/domain/entity"
	"os"
)

type UseCaseJsonPersisterAdapter interface {
	Load(location string) (*entity.PersistentState, error)
	Save(state *entity.PersistentState, location string) error
}

type jsonPersister struct{}

func NewJsonPersister() *jsonPersister {
	return &jsonPersister{}
}

func (p jsonPersister) Load(location string) (*entity.PersistentState, error) {
	data, err := os.ReadFile(location)
	if err != nil {
		return nil, err
	}
	state := &entity.PersistentState{}
	err = json.Unmarshal(data, state)
	return state, err
}

func (p jsonPersister) Save(state *entity.PersistentState, location string) error {
	data, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(location, data, 0644)
}
