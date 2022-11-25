package pkg

import (
	"encoding/json"
	"errors"

	"graft/pkg/domain"
	secondaryAdapter "graft/pkg/infrastructure/adapter/secondary"
	secondaryPort "graft/pkg/infrastructure/port/secondary"
)

func AddClusterPeer(newPeer domain.Peer, clusterPeer domain.Peer) (chan struct{}, error) {
	// 1. Get cluster leader
	leaderNotFound := false
	leader, err := getClusterLeader(clusterPeer)
	switch err {
	case nil:
		break
	case errLeaderNotFound:
		leaderNotFound = true
	default:
		return nil, err
	}

	if leaderNotFound {
		// Leader is not found, we cannot get an up-to-date configuration
		// We should exit now and not attempt to start a node

		// Fall back to a cluster peer to fetch the configuration
		// the cluster configuration is potentially stale
		config, err := client.ClusterConfiguration(clusterPeer)
		if err != nil {
			return nil, err
		}

		// newPeer is unknown to the cluster configuration
		// or is not an active cluster peer.
		// Because leader is not found, we cannot perform
		// a cluster config update
		if p, ok := config.Peers[newPeer.Id]; !ok || !p.Active {
			return nil, errLeaderNotFound
		}

		// 2. Start newPeer
		quit := Start(
			newPeer.Id,
			newPeer.Host,
			config.Peers,
			config.ElectionTimeout,
			config.LeaderHeartbeat,
		)

		return quit, nil

	} else {
		// Cluster configuration is up to date, because it is from the leader
		config, _ := client.ClusterConfiguration(*leader)

		// 2. Add newPeer to cluster configuration
		// only if unknown
		newPeerFromLeader, newPeerIsKnown := config.Peers[newPeer.Id]
		if !newPeerIsKnown {
			err = executeConfigurationUpdate(domain.ConfAddPeer, newPeer, *leader)
			if err != nil {
				return nil, err
			}
		}

		// 3. Start newPeer
		quit := Start(
			newPeer.Id,
			newPeer.Host,
			config.Peers,
			config.ElectionTimeout,
			config.LeaderHeartbeat,
		)

		// 4. Set newPeer to active
		// only if peer is unknown or inactive
		if !newPeerIsKnown || !newPeerFromLeader.Active {
			err = executeConfigurationUpdate(domain.ConfActivatePeer, newPeer, *leader)
			if err != nil {
				return nil, err
			}
		}

		return quit, nil
	}
}

func RemoveClusterPeer(oldPeer domain.Peer) error {
	// 1. Get cluster leader
	leader, err := getClusterLeader(oldPeer)
	if err != nil {
		return err
	}

	// 2. Set peer as inactive
	err = executeConfigurationUpdate(domain.ConfDeactivatePeer, oldPeer, *leader)
	if err != nil {
		return err
	}

	// 3. Remove peer
	err = executeConfigurationUpdate(domain.ConfRemovePeer, oldPeer, *leader)
	if err != nil {
		return err
	}

	// 4. Shutdown peer
	err = Shutdown(oldPeer)
	if err != nil {
		return err
	}

	return nil
}

var (
	client                      = secondaryPort.NewRpcClientPort(secondaryAdapter.NewGrpcClient())
	errLeaderNotFound           = errors.New("leader not found")
	errConfigurationNotFound    = errors.New("cannot fetch cluster configuration")
	errApplyConfigurationUpdate = errors.New("cannot apply cluster configuration update")
)

func Shutdown(peer domain.Peer) error {
	return client.Shutdown(peer)
}

func getClusterLeader(clusterPeer domain.Peer) (*domain.Peer, error) {
	config, err := client.ClusterConfiguration(clusterPeer)
	if err != nil {
		return nil, errConfigurationNotFound
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

func executeConfigurationUpdate(tp domain.ConfigurationUpdateType, peer domain.Peer, leader domain.Peer) error {
	data, _ := json.Marshal(&domain.ConfigurationUpdate{
		Type: tp,
		Peer: peer,
	})
	input := domain.ExecuteInput{
		Type: domain.LogConfiguration,
		Data: data,
	}
	if _, err := client.Execute(leader, &input); err != nil {
		return errApplyConfigurationUpdate
	}
	return nil
}
