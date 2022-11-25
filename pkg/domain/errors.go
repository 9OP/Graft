package domain

import "errors"

var (
	ErrNotLeader    = errors.New("node is not leader")
	ErrShuttingDown = errors.New("node shutting down")
)
