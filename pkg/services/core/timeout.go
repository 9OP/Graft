package core

import (
	"math/rand"
	"os"
	"time"
)

type timeout struct {
	ElectionTimer    *time.Ticker
	LeaderTicker     *time.Ticker
	electionDuration time.Duration
	heartbeatFreq    time.Duration
}

func getRandomTime(t time.Duration) time.Duration {
	rand.Seed(time.Now().UnixNano() + int64(os.Getpid()))
	return time.Duration(rand.Intn(int(t)/2) + int(t)/2)
}

func newTimeout(ed int, hf int) *timeout {
	return &timeout{
		ElectionTimer:    time.NewTicker(1 * time.Second),
		LeaderTicker:     time.NewTicker(1 * time.Second),
		electionDuration: time.Duration(ed),
		heartbeatFreq:    time.Duration(hf),
	}
}

func (t *timeout) resetElectionTimer() {
	t.LeaderTicker.Stop()
	t.ElectionTimer.Stop()
	t.ElectionTimer.Reset(getRandomTime(t.electionDuration) * time.Millisecond)
}

func (t *timeout) resetLeaderTicker() {
	t.ElectionTimer.Stop()
	t.LeaderTicker.Stop()
	t.LeaderTicker.Reset(t.heartbeatFreq * time.Millisecond)
}
