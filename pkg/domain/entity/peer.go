package entity

import (
	"fmt"
	"net/netip"

	log "github.com/sirupsen/logrus"
)

type Peer struct {
	Id    string
	Host  string
	Ports struct {
		P2p string
		Api string
	}
}

func NewPeer(id string, host string, p2p string, api string) *Peer {
	p2pAddr := fmt.Sprintf("%s:%s", host, p2p)
	apiAddr := fmt.Sprintf("%s:%s", host, api)
	if _, err := netip.ParseAddrPort(p2pAddr); err != nil {
		fmt.Println(err)
		log.Fatalf("Invalid addr %s", p2pAddr)
	}
	if _, err := netip.ParseAddrPort(apiAddr); err != nil {
		log.Fatalf("Invalid addr %s", apiAddr)
	}

	return &Peer{
		Id:   id,
		Host: host,
		Ports: struct {
			P2p string
			Api string
		}{P2p: p2p, Api: api},
	}
}

type Peers map[string]Peer

func (p *Peers) AddPeer(newPeer Peer) {
	(*p)[newPeer.Id] = newPeer
}

func (p *Peers) RemovePeer(peerId string) {
	delete(*p, peerId)
}

func (p Peer) Target() string {
	return fmt.Sprintf("%s:%s", p.Host, p.Ports.P2p)
}
