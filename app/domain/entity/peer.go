package entity

type Peer struct {
	Id   string
	Host string
	Port string
}

type Peers map[string]Peer

func (p *Peers) AddPeer(newPeer Peer) {
	(*p)[newPeer.Id] = newPeer
}
func (p *Peers) RemovePeer(peerId string) {
	delete(*p, peerId)
}
