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
	value  role
	signal chan struct{}
}

func NewRole() *Role {
	return &Role{
		value:  Follower,
		signal: make(chan struct{}, 1),
	}
}

func (r *Role) Is(rl role) bool {
	return r.value == rl
}

func (r *Role) Shift(rl role) {
	r.value = rl
	r.signal <- struct{}{}
}

func (r *Role) Signal() chan struct{} {
	return r.signal
}
