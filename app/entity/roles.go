package entity

type role struct {
	string
}

var (
	Undefined = role{""}
	Follower  = role{"Follower"}
	Candidate = role{"Candidate"}
	Leader    = role{"Leader"}
)

type Role struct {
	Value  role
	Signal chan struct{}
}

func NewRole() *Role {
	return &Role{
		Value:  Follower,
		Signal: make(chan struct{}, 1),
	}
}
