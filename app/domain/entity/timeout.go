package entity

import (
	"math/rand"
	"sync"
	"time"
)

type Timeout struct {
	ElectionTimer    *time.Ticker // Should be a timer
	LeaderTicker     *time.Ticker
	electionDuration time.Duration
	heartbeatFreq    time.Duration
	mu               sync.Mutex
}

func getRandomTime(t time.Duration) time.Duration {
	rand.Seed(time.Now().UnixNano())
	return time.Duration(rand.Intn(int(t)/2) + int(t)/2)
}

func NewTimeout(ed int, hf int) *Timeout {
	return &Timeout{
		ElectionTimer:    time.NewTicker(1 * time.Second),
		LeaderTicker:     time.NewTicker(1 * time.Second),
		electionDuration: time.Duration(ed),
		heartbeatFreq:    time.Duration(hf),
	}
}

func (t *Timeout) ResetElectionTimer() {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.LeaderTicker.Stop()
	t.ElectionTimer.Stop()
	rnd := getRandomTime(t.electionDuration)
	t.ElectionTimer.Reset(rnd * time.Millisecond)
}

func (t *Timeout) ResetLeaderTicker() {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.ElectionTimer.Stop()
	t.LeaderTicker.Stop()
	t.LeaderTicker.Reset(t.heartbeatFreq * time.Millisecond)
}
