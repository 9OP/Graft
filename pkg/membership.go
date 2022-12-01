package pkg

import (
	"encoding/json"

	"graft/pkg/domain"
	secondaryAdapter "graft/pkg/infrastructure/adapter/secondary"
	secondaryPort "graft/pkg/infrastructure/port/secondary"
)

func AddSelf(peer domain.Peer) error {
	// 1. Add newPeer to cluster configuration
	err := executeConfigurationUpdate(domain.ConfAddPeer, peer, peer.Target())
	if err != nil {
		return err
	}

	// 2. Set newPeer to active
	err = executeConfigurationUpdate(domain.ConfActivatePeer, peer, peer.Target())
	if err != nil {
		return err
	}

	return nil
}

func AddClusterPeer(newPeer domain.Peer, cluster string) (chan struct{}, error) {
	client := secondaryPort.NewRpcClientPort(secondaryAdapter.NewClusterClient())

	// 1. Get cluster leader
	leaderNotFound := false
	leader, err := getClusterLeader(cluster)
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
		config, err := client.Configuration(cluster)
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
			config.Fsm,
			config.ElectionTimeout,
			config.LeaderHeartbeat,
		)

		return quit, nil

	} else {
		// Cluster configuration is up to date, because it is from the leader
		config, _ := client.Configuration(*leader)

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
			config.Fsm,
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

func RemoveClusterPeer(oldPeer domain.Peer, cluster string) error {
	// 1. Get cluster leader
	leader, err := getClusterLeader(cluster)
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

	return nil
}

func executeConfigurationUpdate(tp domain.ConfigurationUpdateType, peer domain.Peer, leader string) error {
	client := secondaryPort.NewRpcClientPort(secondaryAdapter.NewClusterClient())

	data, _ := json.Marshal(&domain.ConfigurationUpdate{
		Type: tp,
		Peer: peer,
	})
	input := domain.ExecuteInput{
		Type: domain.LogConfiguration,
		Data: data,
	}
	if _, err := client.Execute(leader, &input); err != nil {
		return err
	}
	return nil
}
