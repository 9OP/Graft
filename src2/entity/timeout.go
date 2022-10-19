package entity

import (
	"math/rand"
	"sync"
	"time"
)

type Timeout struct {
	duration int
	*time.Timer
	mu sync.Mutex
}

func NewTimeout(t int) *Timeout {
	return &Timeout{
		t,
		time.NewTimer(time.Duration(t) * time.Millisecond),
		sync.Mutex{},
	}
}

func (t *Timeout) RReset() {
	t.mu.Lock()
	defer t.mu.Unlock()

	rand.Seed(time.Now().UnixNano())
	timeout := (rand.Intn(t.duration/2) + t.duration/2)
	t.Reset(time.Duration(timeout) * time.Millisecond)
}
