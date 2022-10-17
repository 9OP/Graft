package main

import (
	"encoding/json"
	"graft/src2/entity"
	"os"
)

type Persister struct{}

func (p *Persister) Load(location string) (*entity.PersistentState, error) {
	data, err := os.ReadFile(location)
	if err != nil {
		return nil, err
	}
	state := &entity.PersistentState{}
	err = json.Unmarshal(data, state)
	return state, err
}

func (p *Persister) Save(location string, state *entity.PersistentState) error {
	data, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(location, data, 0644)
}
