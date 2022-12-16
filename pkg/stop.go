package pkg

import (
	secondaryAdapter "graft/pkg/infrastructure/adapter/secondary"
	secondaryPort "graft/pkg/infrastructure/port/secondary"
)

func Shutdown(peer string) error {
	client := secondaryPort.NewRpcClientPort(secondaryAdapter.NewClusterClient())
	return client.Shutdown(peer)
}
