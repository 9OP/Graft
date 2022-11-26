package domain

import "errors"

var (
	ErrNotLeader    = errors.New("node is not leader")
	ErrShuttingDown = errors.New("node shutting down")
	ErrNotActive    = errors.New("node is not active yet")
)
