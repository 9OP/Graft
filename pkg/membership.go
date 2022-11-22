package pkg

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"

	"graft/pkg/domain"
)

func AddMember(peer domain.Peer, addr net.IP, port uint16) {
	// Take first peer and seed dummy request
	resp, err := http.Get(fmt.Sprintf("http://%v:%v/query", addr, port))
	if err != nil {
		log.Fatalln(err)
	}

	// Get Leader addr
	location, err := resp.Location()
	if err != nil {
		log.Fatalln(err)
	}

	// Create body
	data, _ := json.Marshal(&domain.ConfigurationUpdate{
		Type: domain.ConfAddPeer,
		Peer: peer,
	})
	body, _ := json.Marshal(&domain.ApiCommand{
		Type: domain.LogConfiguration,
		Data: data,
	})

	// Send request
	resp, err = http.Post(location.String(), "application/json", bytes.NewBuffer(body))
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatalln(err)
		}
		log.Println(string(body))
	} else {
		log.Fatalln(resp.Status)
	}
}

func RemoveMember(peerId string, addr net.IP, port uint16) {
	//
}
