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
