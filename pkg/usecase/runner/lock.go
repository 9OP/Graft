package runner

import "sync/atomic"

type binaryLock struct {
	c uint32
}

func (l *binaryLock) Lock() bool {
	return atomic.CompareAndSwapUint32(&l.c, 0, 1)
}

func (l *binaryLock) Unlock() {
	atomic.StoreUint32(&l.c, 0)
}
