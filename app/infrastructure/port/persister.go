package port

import "graft/app/domain/entity"

// Move to its own use case when implementating log compaction

type Persister interface {
	Load() (*entity.Persistent, error)
	Save(state *entity.Persistent) error
}
