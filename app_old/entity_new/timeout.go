package entitynew

import (
	"math/rand"
	"sync"
	"time"
)

type Timeout struct {
	electionDuration int
	leaderDuration   int
	ElectionTimer    *time.Timer
	LeaderTicker     *time.Ticker
	mu               sync.Mutex
}

func randomizeTime(t int) time.Duration {
	rand.Seed(time.Now().UnixNano())
	timeout := (rand.Intn(t/2) + t/2)
	return time.Duration(timeout)
}

func NewTimeout(ed int, ld int) *Timeout {
	return &Timeout{electionDuration: ed, leaderDuration: ld}
}

func (t *Timeout) ResetElectionTimer() {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.LeaderTicker.Stop()
	t.ElectionTimer.Stop()
	t.ElectionTimer = time.NewTimer(randomizeTime(t.electionDuration) * time.Millisecond)
}

func (t *Timeout) ResetLeaderTicker() {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.ElectionTimer.Stop()
	t.LeaderTicker.Stop()
	t.LeaderTicker = time.NewTicker(time.Duration(t.leaderDuration) * time.Millisecond)
}
