package pkg

import (
	"encoding/json"
	"errors"

	"graft/pkg/domain"
	secondaryAdapter "graft/pkg/infrastructure/adapter/secondary"
	secondaryPort "graft/pkg/infrastructure/port/secondary"
)

func AddClusterPeer(newPeer domain.Peer, clusterPeer domain.Peer) error {
	// 1. Get cluster leader
	leader, err := getClusterLeader(clusterPeer)
	if err != nil {
		return err
	}

	// 2. Add newPeer to cluster configuration
	err = addPeerConfiguration(newPeer, *leader)
	if err != nil {
		return err
	}

	// 3. Start newPeer

	// 4. Set newPeer to active
	err = activatePeerConfiguration(newPeer, *leader)
	if err != nil {
		return err
	}

	return nil
}

var (
	client            = secondaryPort.NewRpcClientPort(secondaryAdapter.NewGrpcClient())
	errLeaderNotFound = errors.New("leader not found")
)

func getClusterLeader(clusterPeer domain.Peer) (*domain.Peer, error) {
	config, err := client.ClusterConfiguration(clusterPeer)
	if err != nil {
		return nil, err
	}

	// If leader available, return it
	if leader, ok := config.Peers[config.LeaderId]; ok {
		return &leader, nil
	}

	// Otherwise ask known peers about the leader
	for _, peer := range config.Peers {
		config, err := client.ClusterConfiguration(peer)
		if err != nil {
			return nil, err
		}

		if leader, ok := config.Peers[config.LeaderId]; ok {
			return &leader, nil
		}
	}

	return nil, errLeaderNotFound
}

func addPeerConfiguration(newPeer domain.Peer, leader domain.Peer) error {
	data, _ := json.Marshal(&domain.ConfigurationUpdate{
		Type: domain.ConfAddPeer,
		Peer: domain.Peer{Id: newPeer.Id, Host: newPeer.Host, Active: false},
	})
	input := domain.ExecuteInput{
		Type: domain.LogConfiguration,
		Data: data,
	}
	res, err := client.Execute(leader, &input)
	if err != nil {
		return err
	}
	if res.Err != nil {
		return res.Err
	}
	return nil
}

func activatePeerConfiguration(newPeer domain.Peer, leader domain.Peer) error {
	data, _ := json.Marshal(&domain.ConfigurationUpdate{
		Type: domain.ConfActivatePeer,
		Peer: domain.Peer{Id: newPeer.Id, Host: newPeer.Host, Active: true},
	})
	input := domain.ExecuteInput{
		Type: domain.LogConfiguration,
		Data: data,
	}
	res, err := client.Execute(leader, &input)
	if err != nil {
		return err
	}
	if res.Err != nil {
		return res.Err
	}
	return nil
}
