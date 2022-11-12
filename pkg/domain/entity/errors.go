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
type UnknownLeaderError struct {
	*domainError
}

func NewNotLeaderError(leader Peer) *NotLeaderError {
	return &NotLeaderError{leader, &domainError{"not leader"}}
}

func NewUnknownLeaderError() *UnknownLeaderError {
	return &UnknownLeaderError{&domainError{"unknown leader"}}
}
