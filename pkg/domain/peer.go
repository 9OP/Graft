package domain

import (
	"net/netip"

	"graft/pkg/utils"
)

type Peer struct {
	Id     string
	Addr   netip.AddrPort
	Active bool
}

type Peers map[string]Peer

func (p Peers) addPeer(newPeer Peer) Peers {
	peersCopy := utils.CopyMap(p)
	peersCopy[newPeer.Id] = newPeer
	return peersCopy
}

func (p Peers) removePeer(peerId string) Peers {
	peersCopy := utils.CopyMap(p)
	delete(peersCopy, peerId)
	return peersCopy
}

func (p Peers) setPeerStatus(peerId string, activate bool) Peers {
	peersCopy := utils.CopyMap(p)
	if peer, ok := peersCopy[peerId]; ok {
		peer.Active = activate
		peersCopy[peerId] = peer
		return peersCopy
	}
	return p
}

func (p Peers) activatePeer(peerId string) Peers {
	return p.setPeerStatus(peerId, true)
}

func (p Peers) deactivatePeer(peerId string) Peers {
	return p.setPeerStatus(peerId, false)
}

func (p Peer) Target() string {
	return p.Addr.String()
}
