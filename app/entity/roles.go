package entity

import "sync"

type role struct {
	string
}

var (
	Undefined = role{""}
	Follower  = role{"Follower"}
	Candidate = role{"Candidate"}
	Leader    = role{"Leader"}
)

type Role struct {
	value  role
	signal chan struct{}
	mu     sync.RWMutex
}

func NewRole() *Role {
	r := &Role{signal: make(chan struct{}, 1)} // buffer 1 required for async reading
	r.Shift(Follower)
	return r
}

func (r *Role) Is(rl role) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.value == rl
}

func (r *Role) Shift(rl role) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.value = rl
	r.signal <- struct{}{}
}

func (r *Role) Signal() chan struct{} {
	return r.signal
}
