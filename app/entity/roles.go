package entity

type Role struct {
	string
}

var (
	Undefined = Role{""}
	Follower  = Role{"Follower"}
	Candidate = Role{"Candidate"}
	Leader    = Role{"Leader"}
)
