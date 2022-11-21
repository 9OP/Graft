package domain

import (
	"reflect"
	"testing"
)

func TestNewPeer(t *testing.T) {
	res := []struct {
		// input
		id   string
		host string
		p2p  string
		api  string
		// output
		err error
	}{
		{"id", "localhost", "1234", "1234", errInvalidAddr("localhost:1234")},
		{"id", "127.0.0.1", "1234", "blublu", errInvalidAddr("127.0.0.1:blublu")},
		{"id", "127.0.0.1", "blublu", "1234", errInvalidAddr("127.0.0.1:blublu")},
		{"id", "255.255.255.255", "1234", "1234", errInvalidAddr("255.255.255.255:1234")},
		{"id", "blublu", "1234", "1234", errInvalidAddr("blublu:1234")},
		{"id", "127.0.0.1", "1234", "1234", nil},
	}
	for _, tt := range res {
		t.Run(tt.id, func(t *testing.T) {
			peer, err := NewPeer(tt.id, true, tt.host, tt.p2p, tt.api)
			if err != nil && tt.err != nil {
				if err.Error() != tt.err.Error() {
					t.Errorf("NewPeer err got %v, want %v", err.Error(), tt.err.Error())
				}
				return
			}
			expected := Peer{
				Id:     tt.id,
				Host:   tt.host,
				Active: true,
				Ports: struct {
					P2p string
					Api string
				}{
					P2p: tt.p2p,
					Api: tt.api,
				},
			}
			if !reflect.DeepEqual(expected, *peer) {
				t.Errorf("NewPeer got %v, want %v", peer, expected)
			}
		})
	}
}

func TestAddPeer(t *testing.T) {
	newPeer := Peer{Id: "newPeerId"}
	peers := Peers{
		"peerId": Peer{Id: "peerId"},
	}

	// Add peer
	res := peers.AddPeer(newPeer)
	expected := Peers{
		"peerId":    Peer{Id: "peerId"},
		"newPeerId": Peer{Id: "newPeerId"},
	}
	if !reflect.DeepEqual(res, expected) {
		t.Errorf("AddPeer got %v, want %v", res, expected)
	}

	// Mutate
	if _, ok := peers["newPeerId"]; ok {
		t.Error("mutate copy")
	}
}

func TestRemovePeer(t *testing.T) {
	peers := Peers{
		"peerId":    Peer{Id: "peerId"},
		"oldPeerId": Peer{Id: "oldPeerId"},
	}

	// Add peer
	res := peers.RemovePeer("oldPeerId")
	expected := Peers{
		"peerId": Peer{Id: "peerId"},
	}
	if !reflect.DeepEqual(res, expected) {
		t.Errorf("AddPeer got %v, want %v", res, expected)
	}

	// Mutate
	if _, ok := peers["oldPeerId"]; !ok {
		t.Error("mutate copy")
	}
}
