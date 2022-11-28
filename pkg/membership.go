package pkg

import (
	"encoding/json"
	"errors"

	"graft/pkg/domain"
	secondaryAdapter "graft/pkg/infrastructure/adapter/secondary"
	secondaryPort "graft/pkg/infrastructure/port/secondary"
)

func AddClusterPeer(newPeer domain.Peer, cluster string) (chan struct{}, error) {
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

func Execute(entry string, logType domain.LogType, cluster string) (*domain.ExecuteOutput, error) {
	var peer string

	if logType == domain.LogCommand || logType == domain.LogConfiguration {
		p, err := getClusterLeader(cluster)
		if err != nil {
			return nil, err
		}
		peer = *p
	} else {
		peer = cluster
	}

	input := domain.ExecuteInput{
		Type: logType,
		Data: []byte(entry),
	}

	return client.Execute(peer, &input)
}

func LeadeshipTransfer(peer string) error {
	return client.LeadershipTransfer(peer)
}

func Shutdown(peer string) error {
	return client.Shutdown(peer)
}

func ClusterConfiguration(cluster string) (*domain.ClusterConfiguration, error) {
	leader, err := getClusterLeader(cluster)
	if err != nil {
		return nil, err
	}
	return client.Configuration(*leader)
}

var (
	client            = secondaryPort.NewRpcClientPort(secondaryAdapter.NewClusterClient())
	errLeaderNotFound = errors.New("leader not found")
)

func getClusterLeader(cluster string) (*string, error) {
	config, err := client.Configuration(cluster)
	if err != nil {
		return nil, err
	}

	// If leader available, return it
	if leader, ok := config.Peers[config.LeaderId]; ok {
		target := leader.Target()
		return &target, nil
	}

	// Otherwise ask known peers about the leader
	for _, peer := range config.Peers {
		config, err := client.Configuration(peer.Target())
		if err != nil {
			// No worries, peer might be out
			// ask the next one for leader info
			continue
		}

		if leader, ok := config.Peers[config.LeaderId]; ok {
			target := leader.Target()
			return &target, nil
		}
	}

	// The leader is still unknown.
	// either:
	// - network partition
	// - cluster has failed quorum and no leader can be elected
	return nil, errLeaderNotFound
}

func executeConfigurationUpdate(tp domain.ConfigurationUpdateType, peer domain.Peer, leader string) error {
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
