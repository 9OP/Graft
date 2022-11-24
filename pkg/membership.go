package pkg

import (
	"encoding/json"
	"errors"

	"graft/pkg/domain"
	secondaryAdapter "graft/pkg/infrastructure/adapter/secondary"
	secondaryPort "graft/pkg/infrastructure/port/secondary"
)

func AddClusterPeer(newPeer domain.Peer, clusterPeer domain.Peer, quit chan struct{}) error {
	// 1. Get cluster leader
	leaderNotFound := false
	leader, err := getClusterLeader(clusterPeer)
	switch err {
	case nil:
		break
	case errLeaderNotFound:
		leaderNotFound = true
	default:
		return err
	}

	if leaderNotFound {
		// Leader is not found, we cannot get an up-to-date configuration
		// We should exit now and not attempt to start a node

		// Fall back to a cluster peer to fetch the configuration
		// the cluster configuration is potentially stale
		config, err := client.ClusterConfiguration(clusterPeer)
		if err != nil {
			return err
		}

		// newPeer is unknown to the cluster configuration
		// or is not an active cluster peer.
		// Because leader is not found, we cannot perform
		// a cluster config update
		if p, ok := config.Peers[newPeer.Id]; !ok || !p.Active {
			return errLeaderNotFound
		}

		// 2. Start newPeer
		go Start(
			newPeer.Id,
			newPeer.Host,
			config.Peers,
			config.ElectionTimeout,
			config.LeaderHeartbeat,
			quit,
		)

	} else {
		// Cluster configuration is up to date, because it is from the leader
		config, _ := client.ClusterConfiguration(*leader)

		// 2. Add newPeer to cluster configuration
		// only if unknown
		newPeerFromLeader, newPeerIsKnown := config.Peers[newPeer.Id]
		if !newPeerIsKnown {
			err = addPeerConfiguration(newPeer, *leader)
			if err != nil {
				return err
			}
		}

		// 3. Start newPeer
		go Start(
			newPeer.Id,
			newPeer.Host,
			config.Peers,
			config.ElectionTimeout,
			config.LeaderHeartbeat,
			quit,
		)

		// 4. Set newPeer to active
		// only if peer is unknown or inactive
		if !newPeerIsKnown || !newPeerFromLeader.Active {
			err = activatePeerConfiguration(newPeer, *leader)
			if err != nil {
				return err
			}
		}
	}

	// Block forever: crado
	<-make(chan int)

	return nil
}

func RemoveClusterPeer() {}

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
			// No worries, peer might be out
			// ask the next one for leader info
			continue
		}

		if leader, ok := config.Peers[config.LeaderId]; ok {
			return &leader, nil
		}
	}

	// The leader is still unknown.
	// either:
	// - network partition
	// - cluster has failed quorum and no leader can be elected
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
	_, err := client.Execute(leader, &input)
	if err != nil {
		return err
	}
	// if res.Err != nil {
	// 	return res.Err
	// }
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
	_, err := client.Execute(leader, &input)
	if err != nil {
		return err
	}
	// if res.Err != nil {
	// 	return res.Err
	// }
	return nil
}
