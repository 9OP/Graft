package repository

import (
	"encoding/json"
	"graft/app/domain/entity"
	"os"
)

// Implements state.Repository
type JsonPersister struct{}

func (p *JsonPersister) Load(location string) (*entity.Persistent, error) {
	data, err := os.ReadFile(location)
	if err != nil {
		return nil, err
	}
	state := &entity.Persistent{}
	err = json.Unmarshal(data, state)
	return state, err
}

func (p *JsonPersister) Save(location string, state *entity.Persistent) error {
	data, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(location, data, 0644)
}
