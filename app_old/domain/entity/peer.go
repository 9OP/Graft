package entity

type Peer struct {
	Id   string
	Host string
	Port string
}

type Peers = map[string]Peer
