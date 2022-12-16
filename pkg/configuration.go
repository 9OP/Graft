package pkg

import (
	"errors"
	"time"

	"graft/pkg/domain"
	secondaryAdapter "graft/pkg/infrastructure/adapter/secondary"
	secondaryPort "graft/pkg/infrastructure/port/secondary"
)

func ClusterConfiguration(cluster string) (*domain.ClusterConfiguration, error) {
	client := secondaryPort.NewRpcClientPort(secondaryAdapter.NewClusterClient())

	leader, err := getClusterLeader(cluster)
	if err != nil {
		return nil, err
	}
	return client.Configuration(*leader)
}

var errLeaderNotFound = errors.New("leader not found")

func getClusterLeader(cluster string) (*string, error) {
	client := secondaryPort.NewRpcClientPort(secondaryAdapter.NewClusterClient())

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

	// The leader is still unknown. Either:
	// - network partition
	// - cluster has failed quorum and no leader can be elected
	return nil, errLeaderNotFound
}

func getClusterLeaderWithTimeout(cluster string) (*string, error) {
	timeout := time.NewTimer(5 * time.Second)
	polling := time.NewTicker(200 * time.Millisecond)

	for {
		select {
		case <-polling.C:
			if l, err := getClusterLeader(cluster); err == nil {
				return l, nil
			}
		case <-timeout.C:
			return nil, errLeaderNotFound

		}
	}
}
