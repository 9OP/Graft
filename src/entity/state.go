package entity

type Peer struct {
	id   string
	host string
	port string
}

type State struct {
	commitIndex uint32
	lastApplied uint32
	nextIndex   []string // leader only
	matchIndex  []string // leader only
}
