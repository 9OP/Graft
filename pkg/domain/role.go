package domain

type Role uint8

const (
	_ Role = iota
	Follower
	Candidate
	Leader
)

func (r Role) String() string {
	switch r {
	case Follower:
		return "follower"
	case Candidate:
		return "candidate"
	case Leader:
		return "leader"
	}
	return "undefined"
}
