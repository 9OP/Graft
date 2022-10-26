package entity

import (
	"fmt"
)

type Peer struct {
	Id   string
	host string
	Port string
}

func NewPeer(id string, host string, port string) *Peer {
	// validate IP
	// net.ParseIP("").String()
	return &Peer{Id: id, host: host, Port: port}
}

type Peers map[string]Peer

func (p *Peers) AddPeer(newPeer Peer) {
	(*p)[newPeer.Id] = newPeer
}
func (p *Peers) RemovePeer(peerId string) {
	delete(*p, peerId)
}
func (p Peer) Target() string {
	return fmt.Sprintf("%s:%s", p.host, p.Port)
}
