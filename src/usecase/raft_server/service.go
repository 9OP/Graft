package raft_server

type Service struct {
	repository *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repository: repo}
}

func (s *Service) RunFollower(follower Follower)    {}
func (s *Service) RunCandidate(candidate Candidate) {}
func (s *Service) RunLeader(leader Leader)          {}
