package cluster

type repository interface {
	Execute(entry string) chan interface{}
	IsLeader() bool
}

type UseCase interface {
	// Writes with ExecuteCommand (leader only)
	// Reads with ExecuteQuery (leader + followers)
	ExecuteCommand(command string) (interface{}, error)
	//ExecuteQuery(query string) (interface{}, error)
}
