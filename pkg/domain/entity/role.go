package entity

type Role struct {
	string
}

func (r Role) String() string {
	return r.string
}

var (
	Undefined = Role{""}
	Follower  = Role{"Follower"}
	Candidate = Role{"Candidate"}
	Leader    = Role{"Leader"}
)
