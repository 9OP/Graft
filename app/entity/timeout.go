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

func randomizeTime(t int) time.Duration {
	rand.Seed(time.Now().UnixNano())
	timeout := (rand.Intn(t/2) + t/2)
	return time.Duration(timeout)
}

func NewTimeout(t int) *Timeout {
	rnd := randomizeTime(t)

	return &Timeout{
		base{
			t,
			sync.Mutex{},
		},
		time.NewTimer(rnd * time.Millisecond),
	}
}

func (t *Timeout) RReset() {
	t.mu.Lock()
	defer t.mu.Unlock()

	rnd := randomizeTime(t.duration)
	t.Stop()
	t.Reset(rnd * time.Millisecond)
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
