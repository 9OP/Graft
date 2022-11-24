package domain

import "errors"

var ErrNotLeader = errors.New("node is not leader")
