package api

type UseCase interface {
	ExecuteCommand(command string) ([]byte, error)
	ExecuteQuery(query string, weakConsistency bool) ([]byte, error)
}
