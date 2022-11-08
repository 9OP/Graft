package entity

type domainError struct {
	message string
}

func (e *domainError) Error() string {
	return e.message
}

type NotLeaderError struct {
	Leader Peer
	*domainError
}
type TimeoutError struct {
	*domainError
}

func NewNotLeaderError(leader Peer) *NotLeaderError {
	return &NotLeaderError{leader, &domainError{"not leader"}}
}

func NewTimeoutError() *TimeoutError {
	return &TimeoutError{&domainError{"timeout"}}
}
