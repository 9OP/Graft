package entity

import (
	"math/rand"
	"time"
)

const ELECTION_TIMEOUT = 350

type Timeout struct {
	time.Timer // Use ticker ?
}

func (t *Timeout) RReset() {
	rand.Seed(time.Now().UnixNano())
	timeout := (rand.Intn(ELECTION_TIMEOUT/2) + ELECTION_TIMEOUT/2)
	t.Reset(time.Duration(timeout) * time.Millisecond)
}
