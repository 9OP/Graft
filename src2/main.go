package main

import (
	"graft/src2/usecase/state"
)

func main() {
	// Parse args

	// Create stateService
	persister := Persister{}
	stateService := state.NewService("state.json", &persister)
	stateService.LoadState()

	// Create rpcService

	// Create raftService

}
