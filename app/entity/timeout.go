package entity

import (
	"math/rand"
	"sync"
	"time"
)

type base struct {
	duration int
	mu       sync.Mutex
}

type Timeout struct {
	base
	*time.Timer
}

type Ticker struct {
	base
	*time.Ticker
}

func NewTimeout(t int) *Timeout {
	return &Timeout{
		base{
			t,
			sync.Mutex{},
		},
		time.NewTimer(time.Duration(t) * time.Millisecond),
	}
}

func (t *Timeout) RReset() {
	t.mu.Lock()
	defer t.mu.Unlock()

	rand.Seed(time.Now().UnixNano())
	timeout := (rand.Intn(t.duration/2) + t.duration/2)
	t.Reset(time.Duration(timeout) * time.Millisecond)
}

func NewTicker(t int) *Ticker {
	return &Ticker{
		base{
			t,
			sync.Mutex{},
		},
		time.NewTicker(time.Duration(t) * time.Millisecond),
	}
}

func (t *Ticker) Start() {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.Reset(time.Duration(t.duration) * time.Millisecond)
}
