package pkg

import (
	"graft/pkg/domain"
	secondaryAdapter "graft/pkg/infrastructure/adapter/secondary"
	secondaryPort "graft/pkg/infrastructure/port/secondary"
)

func Execute(entry string, logType domain.LogType, cluster string) (*domain.ExecuteOutput, error) {
	client := secondaryPort.NewRpcClientPort(secondaryAdapter.NewClusterClient())

	var peer string

	if logType == domain.LogCommand || logType == domain.LogConfiguration {
		leader, err := getClusterLeader(cluster)
		if err != nil {
			return nil, err
		}
		peer = *leader
	} else {
		peer = cluster
	}

	input := domain.ExecuteInput{
		Type: logType,
		Data: []byte(entry),
	}

	return client.Execute(peer, &input)
}
