package entity

// Create /role/entity
type Role struct {
	string
}

var (
	Undefined = Role{""}
	Follower  = Role{"Follower"}
	Candidate = Role{"Candidate"}
	Leader    = Role{"Leader"}
)

// UseCase

type IFollower interface {
	upCandidate()
	grantVote()
}

type ICandidate interface {
	downFollower()
	upLeader()
	grantVote()
}

type ILeader interface {
	downFollower()
	// broadcast ?
}
