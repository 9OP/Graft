package domain

import (
	"fmt"
	"net/netip"

	"graft/pkg/utils"
)

type Peer struct {
	Id     string
	Host   string
	Active bool
	Ports  struct {
		P2p string
		Api string
	}
}

func errInvalidAddr(addr string) error {
	return fmt.Errorf("invalid addr format %s", addr)
}

func NewPeer(id string, active bool, host string, p2p string, api string) (*Peer, error) {
	p2pAddr := fmt.Sprintf("%s:%s", host, p2p)
	apiAddr := fmt.Sprintf("%s:%s", host, api)
	if _, err := netip.ParseAddrPort(p2pAddr); err != nil {
		return nil, errInvalidAddr(p2pAddr)
	}
	if _, err := netip.ParseAddrPort(apiAddr); err != nil {
		return nil, errInvalidAddr(apiAddr)
	}

	return &Peer{
		Id:     id,
		Host:   host,
		Active: active,
		Ports: struct {
			P2p string
			Api string
		}{P2p: p2p, Api: api},
	}, nil
}

type Peers map[string]Peer

func (p Peers) AddPeer(newPeer Peer) Peers {
	peersCopy := utils.CopyMap(p)
	peersCopy[newPeer.Id] = newPeer
	return peersCopy
}

func (p Peers) RemovePeer(peerId string) Peers {
	peersCopy := utils.CopyMap(p)
	delete(peersCopy, peerId)
	return peersCopy
}

func (p Peer) TargetP2p() string {
	return fmt.Sprintf("%s:%s", p.Host, p.Ports.P2p)
}

func (p Peer) TargetApi() string {
	return fmt.Sprintf("%s:%s", p.Host, p.Ports.Api)
}
