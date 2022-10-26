package adapter

import (
	"encoding/json"
	"graft/app/domain/entity"
	"os"
)

type jsonPersister struct {
	location string
}

func NewJsonPersister(location string) *jsonPersister {
	return &jsonPersister{location}
}

func (p jsonPersister) Load() (*entity.Persistent, error) {
	data, err := os.ReadFile(p.location)
	if err != nil {
		return nil, err
	}
	state := &entity.Persistent{}
	err = json.Unmarshal(data, state)
	return state, err
}

func (p jsonPersister) Save(state *entity.Persistent) error {
	data, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(p.location, data, 0644)
}
