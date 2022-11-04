package secondaryAdapter

import (
	"encoding/json"
	"graft/pkg/domain/entity"
	"os"
)

type UseCaseJsonPersisterAdapter interface {
	Load(location string) (*entity.Persistent, error)
	Save(state *entity.Persistent, location string) error
}

type jsonPersister struct{}

func NewJsonPersister() *jsonPersister {
	return &jsonPersister{}
}

func (p jsonPersister) Load(location string) (*entity.Persistent, error) {
	data, err := os.ReadFile(location)
	if err != nil {
		return nil, err
	}
	state := &entity.Persistent{}
	err = json.Unmarshal(data, state)
	return state, err
}

func (p jsonPersister) Save(state *entity.Persistent, location string) error {
	data, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(location, data, 0644)
}
